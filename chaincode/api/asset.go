package api

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/usnistgov/blossom/chaincode/collections"
	events "github.com/usnistgov/blossom/chaincode/ngac/epp"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
	decider "github.com/usnistgov/blossom/chaincode/ngac/pdp"
	"strings"

	"github.com/usnistgov/blossom/chaincode/model"
)

type (
	CheckoutRequest struct {
		Asset  string `json:"asset,omitempty"`
		Amount int    `json:"amount,omitempty"`
	}

	CheckinRequest struct {
		Asset    string   `json:"asset,omitempty"`
		Licenses []string `json:"licenses,omitempty"`
	}
)

func NewLicenseContract() AssetInterface {
	return &BlossomSmartContract{}
}

func (b *BlossomSmartContract) assetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetPrivateData(collections.Catalog(), model.AssetKey(id))
	if err != nil {
		return false, fmt.Errorf("error checking if asset id %q already exists on the ledger: %w", id, err)
	}

	return data != nil, nil
}

func (b *BlossomSmartContract) OnboardAsset(ctx contractapi.TransactionContextInterface, id string, name string, onboardDate string, expiration string) error {
	if ok, err := b.assetExists(ctx, id); err != nil {
		return fmt.Errorf("error checking if asset already exists: %w", err)
	} else if ok {
		return fmt.Errorf("an asset with the ID %q already exists", id)
	}

	assetInput, err := getOnboardAssetTransientInput(ctx)
	if err != nil {
		return fmt.Errorf("error getting transient input: %w", err)
	}

	if len(assetInput.Licenses) == 0 {
		return fmt.Errorf("licenses cannot be nil")
	}

	// ngac check
	if err = pdp.CanOnboardAsset(ctx); err != nil {
		return fmt.Errorf("ngac check failed: %w", err)
	}

	// public info - id, name, available (=total), expiration
	assetPub := &model.AssetPublic{
		ID:             id,
		Name:           name,
		Available:      len(assetInput.Licenses),
		OnboardingDate: onboardDate,
		Expiration:     expiration,
	}

	bytes, err := json.Marshal(assetPub)
	if err != nil {
		return fmt.Errorf("error marshaling asset %q: %w", name, err)
	}

	// put in catalog pdc
	if err = ctx.GetStub().PutPrivateData(collections.Catalog(), model.AssetKey(id), bytes); err != nil {
		return fmt.Errorf("error adding asset to catalog private data collection: %w", err)
	}

	licenses := make([]string, 0)
	licenseMap := make(map[string]string)
	for _, license := range assetInput.Licenses {
		licenses = append(licenses, license.LicenseID)
		licenseMap[license.LicenseID] = license.Expiration
	}

	assetPvt := model.AssetPrivate{
		TotalAmount:       len(assetInput.Licenses),
		Licenses:          licenseMap,
		AvailableLicenses: licenses,
		CheckedOut:        make(map[string]map[string]string),
	}

	if bytes, err = json.Marshal(assetPvt); err != nil {
		return fmt.Errorf("error marshaling asset %q: %w", name, err)
	}

	// add license to licenses private data
	if err = ctx.GetStub().PutPrivateData(collections.Licenses(), model.AssetKey(id), bytes); err != nil {
		return fmt.Errorf("error adding asset to ledger: %w", err)
	}

	// ngac event
	return events.ProcessOnboardAsset(ctx, collections.Catalog(), id)
}

func (b *BlossomSmartContract) OffboardAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	if ok, err := b.assetExists(ctx, assetID); err != nil {
		return fmt.Errorf("error checking if asset exists: %w", err)
	} else if !ok {
		return nil
	}

	var (
		asset *model.Asset
		err   error
	)

	// ngac check
	if err := pdp.CanOffboardAsset(ctx); err != nil {
		return fmt.Errorf("ngac check failed: %w", err)
	}

	if asset, err = b.GetAsset(ctx, assetID); err != nil {
		return fmt.Errorf("error getting asset info: %w", err)
	}

	// check that all licenses have been returned
	if asset.Available != len(asset.Licenses) {
		return fmt.Errorf("asset %s still has licenses checked out: %w", assetID, err)
	}

	// remove asset from catalog
	if err = ctx.GetStub().DelPrivateData(collections.Catalog(), model.AssetKey(assetID)); err != nil {
		return fmt.Errorf("error offboarding asset from catalog pdc: %w", err)
	}

	// remove license licenses pdc
	if err = ctx.GetStub().DelPrivateData(collections.Licenses(), model.AssetKey(assetID)); err != nil {
		return fmt.Errorf("error offboarding asset from licenses pdc: %w", err)
	}

	// ngac event
	return events.ProcessOffboardAsset(ctx, collections.Catalog(), assetID)
}

func (b *BlossomSmartContract) GetAssets(ctx contractapi.TransactionContextInterface) ([]*model.AssetPublic, error) {
	// ngac check
	if err := pdp.CanViewAssets(ctx); err != nil {
		return nil, fmt.Errorf("ngac check failed: %w", err)
	}

	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collections.Catalog(), "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	assets := make([]*model.AssetPublic, 0)
	for resultsIterator.HasNext() {
		var queryResponse *queryresult.KV
		if queryResponse, err = resultsIterator.Next(); err != nil {
			return nil, err
		}

		// assets on the ledger begin with the asset prefix -- ignore other results
		if !strings.HasPrefix(queryResponse.Key, model.AssetPrefix) {
			continue
		}

		asset := model.NewAssetPublic()
		if err = json.Unmarshal(queryResponse.Value, asset); err != nil {
			return nil, err
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

func (b *BlossomSmartContract) GetAsset(ctx contractapi.TransactionContextInterface, id string) (*model.Asset, error) {
	if ok, err := b.assetExists(ctx, id); err != nil {
		return nil, fmt.Errorf("error checking if asset exists: %w", err)
	} else if !ok {
		return nil, fmt.Errorf("an asset with the ID %q does not exist", id)
	}

	var (
		assetPub = model.NewAssetPublic()
		assetPvt = model.NewAssetPrivate()
		bytes    []byte
		err      error
	)

	if bytes, err = ctx.GetStub().GetPrivateData(collections.Catalog(), model.AssetKey(id)); err != nil {
		return nil, fmt.Errorf("error getting asset from private data: %w", err)
	}

	// ngac check
	if err = pdp.CanViewAssetPublic(ctx); err != nil {
		return nil, fmt.Errorf("ngac check on asset public failed: %w", err)
	}

	if err = json.Unmarshal(bytes, assetPub); err != nil {
		return nil, fmt.Errorf("error unmarshaling asset public info: %w", err)
	}

	if bytes, err = ctx.GetStub().GetPrivateData(collections.Licenses(), model.AssetKey(id)); err != nil {
		mspid, _ := ctx.GetClientIdentity().GetMSPID()
		fmt.Printf("error occurred reading pvtdata for user in org %s: %v\n", mspid, err)
	} else {
		if err = json.Unmarshal(bytes, assetPvt); err != nil {
			return nil, fmt.Errorf("error deserializing account private info: %w", err)
		}
	}

	return &model.Asset{
		ID:                assetPub.ID,
		Name:              assetPub.Name,
		Available:         assetPub.Available,
		OnboardingDate:    assetPub.OnboardingDate,
		Expiration:        assetPub.Expiration,
		TotalAmount:       assetPvt.TotalAmount,
		Licenses:          assetPvt.Licenses,
		AvailableLicenses: assetPvt.AvailableLicenses,
		CheckedOut:        assetPvt.CheckedOut,
	}, nil
}

func (b *BlossomSmartContract) RequestCheckout(ctx contractapi.TransactionContextInterface) error {
	transientInput, err := getRequestCheckoutTransientInput(ctx)
	if err != nil {
		return fmt.Errorf("error getting transient input: %w", err)
	}

	var (
		account string
		bytes   []byte
	)

	// check requested asset exists
	if bytes, err = ctx.GetStub().GetPrivateData(collections.Catalog(), model.AssetKey(transientInput.AssetID)); err != nil {
		return err
	} else if bytes == nil {
		return fmt.Errorf("asset with id %s does not exist", transientInput.AssetID)
	}

	if account, err = accountName(ctx); err != nil {
		return fmt.Errorf("error getting MSPID from stub: %w", err)
	}

	collection := collections.Account(account)

	// ngac check
	if err = decider.CanRequestCheckout(ctx, account); err != nil {
		return fmt.Errorf("ngac check failed: %w", err)
	}

	key := checkoutRequestKey(account, transientInput.AssetID)

	// check if request has already been made and not approved
	if bytes, err = ctx.GetStub().GetPrivateData(collection, key); err != nil {
		return fmt.Errorf("error reading private data to check if request has been made but not approved: %w", err)
	} else if bytes != nil {
		return fmt.Errorf("request for asset %s alreadys exists for account %s and has not been approved yet", transientInput.AssetID, account)
	}

	req := &CheckoutRequest{transientInput.AssetID, transientInput.Amount}

	if bytes, err = json.Marshal(req); err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	return ctx.GetStub().PutPrivateData(collection, key, bytes)
}

func checkoutRequestKey(account, assetID string) string {
	return fmt.Sprintf("checkout=%s:%s", account, assetID)
}

func (b *BlossomSmartContract) GetCheckoutRequests(ctx contractapi.TransactionContextInterface, account string) ([]CheckoutRequest, error) {
	collection := collections.Account(account)

	key := checkoutRequestKey(account, "")

	iter, err := ctx.GetStub().GetPrivateDataByRange(collection, "", "")
	if err != nil {
		return nil, err
	}

	reqs := make([]CheckoutRequest, 0)
	for iter.HasNext() {
		next := &queryresult.KV{}
		if next, err = iter.Next(); err != nil {
			return nil, fmt.Errorf("error getting next KV: %w", err)
		}

		if !strings.HasPrefix(next.Key, key) {
			continue
		}

		req := CheckoutRequest{}
		if err = json.Unmarshal(next.Value, &req); err != nil {
			return nil, fmt.Errorf("error unmarshaling request: %w", err)
		}

		reqs = append(reqs, req)
	}

	return reqs, nil
}

func (b *BlossomSmartContract) ApproveCheckout(ctx contractapi.TransactionContextInterface) error {
	transientInput, err := getApproveCheckoutTransientInput(ctx)
	if err != nil {
		return fmt.Errorf("error getting transient input: %w", err)
	}

	var (
		acctColl = collections.Account(transientInput.Account)
		key      = checkoutRequestKey(transientInput.Account, transientInput.AssetID)
		bytes    []byte
	)

	// ngac check
	if err = decider.CanApproveCheckout(ctx, transientInput.Account); err != nil {
		return fmt.Errorf("ngac check failed: %w", err)
	}

	// check that request exists
	if bytes, err = ctx.GetStub().GetPrivateData(acctColl, key); err != nil {
		return fmt.Errorf("error checking if request exists: %w", err)
	} else if bytes == nil {
		return fmt.Errorf("request for asset %s does not exist for account %s", transientInput.AssetID, transientInput.Account)
	}

	// delete request key
	if err = ctx.GetStub().DelPrivateData(acctColl, key); err != nil {
		return fmt.Errorf("error deleting request: %w", err)
	}

	req := &CheckoutRequest{}
	if err = json.Unmarshal(bytes, req); err != nil {
		return fmt.Errorf("error unmarshaling request: %w", err)
	}

	acctPub, acctPvt, assetPub, assetPvt, err := getAcctAndAsset(ctx, transientInput.Account, transientInput.AssetID)
	if err != nil {
		return fmt.Errorf("error getting account and asset to process checkout: %w", err)
	}

	if err = checkout(assetPub, assetPvt, acctPub, acctPvt, req.Amount); err != nil {
		return fmt.Errorf("error checking out %s for account %s: %w", transientInput.AssetID, transientInput.Account, err)
	}

	return putAcctAndAsset(ctx, acctPub, acctPvt, assetPub, assetPvt)
}

func checkout(assetPub *model.AssetPublic, assetPvt *model.AssetPrivate, acctPub *model.AccountPublic, acctPvt *model.AccountPrivate, amount int) error {
	// check that the amount requested is less than the amount available
	if amount > assetPub.Available {
		return fmt.Errorf("requested amount %v cannot be greater than the available amount %v",
			amount, assetPub.Available)
	}

	// update available amount
	assetPub.Available -= amount

	// get the available licenses
	fromAvailable := assetPvt.AvailableLicenses[0:amount]
	// update available licenses
	assetPvt.AvailableLicenses = assetPvt.AvailableLicenses[amount:]

	// create the set of licenses that are checked out including expiration dates
	retCheckedOutLicenses := make(map[string]string)
	for _, license := range fromAvailable {
		retCheckedOutLicenses[license] = assetPvt.Licenses[license]
	}

	// update the account assets
	// add to existing asset if they are checking out more of a software asset
	allCheckedOutAssets, ok := acctPvt.Assets[assetPub.ID]
	if !ok {
		allCheckedOutAssets = make(map[string]string)
	}

	for license, exp := range retCheckedOutLicenses {
		allCheckedOutAssets[license] = exp
	}

	// update asset in the account
	acctPvt.Assets[assetPub.ID] = allCheckedOutAssets

	// update the asset's account tracker
	accountCheckedOut := make(map[string]string)
	for license, exp := range allCheckedOutAssets {
		accountCheckedOut[license] = exp
	}
	assetPvt.CheckedOut[acctPub.Name] = accountCheckedOut

	return nil
}

func (b *BlossomSmartContract) GetLicenses(ctx contractapi.TransactionContextInterface, account, assetID string) (map[string]string, error) {
	bytes, err := ctx.GetStub().GetPrivateData(collections.Account(account), model.AccountKey(account))
	if err != nil {
		return nil, fmt.Errorf("error reading account private data: %w", err)
	}

	acctPvt := model.NewAccountPrivate()
	if err = json.Unmarshal(bytes, acctPvt); err != nil {
		return nil, fmt.Errorf("error unmarshaling account private data: %w", err)
	}

	return acctPvt.Assets[assetID], nil
}

func (b *BlossomSmartContract) InitiateCheckin(ctx contractapi.TransactionContextInterface) error {
	transientInput, err := getInitiateCheckinTransientInput(ctx)
	if err != nil {
		return fmt.Errorf("error getting transient input: %w", err)
	}

	account, err := accountName(ctx)
	if err != nil {
		return fmt.Errorf("error getting MSPID from stub: %w", err)
	}

	collection := collections.Account(account)

	// ngac check
	if err = decider.CanInitiateCheckIn(ctx, account); err != nil {
		return fmt.Errorf("ngac check failed: %w", err)
	}

	var (
		key   = checkinRequestKey(account, transientInput.AssetID)
		bytes []byte
	)

	// check if the licenses in the request are really checked out by the account
	if bytes, err = ctx.GetStub().GetPrivateData(collections.Account(account), model.AccountKey(account)); err != nil {
		return fmt.Errorf("error getting account private info from private data: %w", err)
	}

	acctPvt := model.NewAccountPrivate()
	if err = json.Unmarshal(bytes, &acctPvt); err != nil {
		return fmt.Errorf("error unmarshaling account private info: %w", err)
	}

	checkedOut := acctPvt.Assets[transientInput.AssetID]
	for _, returnedKey := range transientInput.Licenses {
		// check that the returned license is leased to the account
		if _, ok := checkedOut[returnedKey]; !ok {
			return fmt.Errorf("returned key %s was not checked out by %s: %w", returnedKey, account, err)
		}
	}

	// check if request has already been made and not approved
	if bytes, err = ctx.GetStub().GetPrivateData(collection, key); err != nil {
		return err
	} else if bytes != nil {
		return fmt.Errorf("request to checkin %s has already been initiated for account %s and has not been processed yet: %w", transientInput.AssetID, account, err)
	}

	req := CheckinRequest{
		Asset:    transientInput.AssetID,
		Licenses: transientInput.Licenses,
	}

	if bytes, err = json.Marshal(req); err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	return ctx.GetStub().PutPrivateData(collection, key, bytes)
}

func (b *BlossomSmartContract) GetInitiatedCheckins(ctx contractapi.TransactionContextInterface, account string) ([]CheckinRequest, error) {
	collection := collections.Account(account)

	key := checkinRequestKey(account, "")

	iter, err := ctx.GetStub().GetPrivateDataByRange(collection, "", "")
	if err != nil {
		return nil, err
	}

	reqs := make([]CheckinRequest, 0)
	for iter.HasNext() {
		next := &queryresult.KV{}
		if next, err = iter.Next(); err != nil {
			return nil, fmt.Errorf("error getting next KV: %w", err)
		}

		if !strings.HasPrefix(next.Key, key) {
			continue
		}

		req := CheckinRequest{}
		if err = json.Unmarshal(next.Value, &req); err != nil {
			return nil, fmt.Errorf("error unmarshaling request: %w", err)
		}

		reqs = append(reqs, req)
	}

	return reqs, nil
}

func checkinRequestKey(account, assetID string) string {
	return fmt.Sprintf("checkin=%s:%s", account, assetID)
}

func checkin(assetPub *model.AssetPublic, assetPvt *model.AssetPrivate, acctPub *model.AccountPublic, acctPvt *model.AccountPrivate, licenses []string) error {
	checkedOut := acctPvt.Assets[assetPub.ID]
	for _, license := range licenses {
		delete(checkedOut, license)
	}

	// if all licenses were returned remove asset from account's checked out
	if len(checkedOut) == 0 {
		delete(acctPvt.Assets, assetPub.ID)
	} else {
		acctPvt.Assets[assetPub.ID] = checkedOut
	}

	accountCheckedOut, ok := assetPvt.CheckedOut[acctPub.Name]
	if !ok {
		return fmt.Errorf("account %s has not checked out any licenses for asset %s", acctPub.Name, assetPub.ID)
	}

	for _, license := range licenses {
		// check that the account has the license checked out
		if _, ok = accountCheckedOut[license]; !ok {
			return fmt.Errorf("returned license %s was not checked out by %s", license, acctPub.Name)
		}

		// remove the returned license from the checked out licenses
		delete(accountCheckedOut, license)

		// add the returned license to the available licenses
		assetPvt.AvailableLicenses = append(assetPvt.AvailableLicenses, license)
	}

	// if all licenses are returned, remove the account from the asset
	if len(accountCheckedOut) == 0 {
		delete(assetPvt.CheckedOut, acctPub.Name)
	}

	// update number of available licenses
	assetPub.Available += len(licenses)

	return nil
}

func (b *BlossomSmartContract) ProcessCheckin(ctx contractapi.TransactionContextInterface) error {
	transientInput, err := getProcessCheckinTransientInput(ctx)
	if err != nil {
		return fmt.Errorf("error getting transient input: %w", err)
	}

	var (
		acctColl = collections.Account(transientInput.Account)
		key      = checkinRequestKey(transientInput.Account, transientInput.AssetID)
		bytes    []byte
	)

	// ngac check
	if err = decider.CanProcessCheckIn(ctx, transientInput.Account); err != nil {
		return fmt.Errorf("ngac check failed: %w", err)
	}

	// check that request exists
	if bytes, err = ctx.GetStub().GetPrivateData(acctColl, key); err != nil {
		return fmt.Errorf("error checking if checkin request exists: %w", err)
	} else if bytes == nil {
		return fmt.Errorf("request to checkin asset %s does not exist for account %s", transientInput.AssetID, transientInput.Account)
	}

	// delete request key
	if err = ctx.GetStub().DelPrivateData(acctColl, key); err != nil {
		return fmt.Errorf("error deleting request: %w", err)
	}

	req := &CheckinRequest{}
	if err = json.Unmarshal(bytes, req); err != nil {
		return fmt.Errorf("error unmarshaling request: %w", err)
	}

	acctPub, acctPvt, assetPub, assetPvt, err := getAcctAndAsset(ctx, transientInput.Account, transientInput.AssetID)
	if err != nil {
		return err
	}

	if err = checkin(assetPub, assetPvt, acctPub, acctPvt, req.Licenses); err != nil {
		return fmt.Errorf("error checking out %s for account %s: %w", transientInput.AssetID, transientInput.Account, err)
	}

	return putAcctAndAsset(ctx, acctPub, acctPvt, assetPub, assetPvt)
}

func putAcctAndAsset(ctx contractapi.TransactionContextInterface, acctPub *model.AccountPublic, acctPvt *model.AccountPrivate,
	assetPub *model.AssetPublic, assetPvt *model.AssetPrivate) (err error) {
	var (
		bytes    []byte
		acctKey  = model.AccountKey(acctPub.Name)
		acctColl = collections.Account(acctPub.Name)
	)

	// put account public
	if bytes, err = json.Marshal(acctPub); err != nil {
		return
	}

	if err = ctx.GetStub().PutState(acctKey, bytes); err != nil {
		return
	}

	// put account private
	if bytes, err = json.Marshal(acctPvt); err != nil {
		return
	}

	if err = ctx.GetStub().PutPrivateData(acctColl, acctKey, bytes); err != nil {
		return
	}

	// put asset public (still pdc)
	if bytes, err = json.Marshal(assetPub); err != nil {
		return
	}

	if err = ctx.GetStub().PutPrivateData(collections.Catalog(), model.AssetKey(assetPub.ID), bytes); err != nil {
		return
	}

	// put asset private
	if bytes, err = json.Marshal(assetPvt); err != nil {
		return
	}

	return ctx.GetStub().PutPrivateData(collections.Licenses(), model.AssetKey(assetPub.ID), bytes)
}

func getAcctAndAsset(ctx contractapi.TransactionContextInterface, account, assetID string) (*model.AccountPublic, *model.AccountPrivate, *model.AssetPublic, *model.AssetPrivate, error) {
	var (
		bytes []byte
		err   error
	)

	// get licenses from license collection
	if bytes, err = ctx.GetStub().GetPrivateData(collections.Licenses(), model.AssetKey(assetID)); err != nil {
		return nil, nil, nil, nil, err
	}

	assetPvt := model.NewAssetPrivate()
	if err = json.Unmarshal(bytes, &assetPvt); err != nil {
		return nil, nil, nil, nil, err
	}

	// get asset public info from catalog collection to update available
	if bytes, err = ctx.GetStub().GetPrivateData(collections.Catalog(), model.AssetKey(assetID)); err != nil {
		return nil, nil, nil, nil, err
	}

	assetPub := model.NewAssetPublic()
	if err = json.Unmarshal(bytes, &assetPub); err != nil {
		return nil, nil, nil, nil, err
	}

	// get account private info from account collection to update available
	if bytes, err = ctx.GetStub().GetPrivateData(collections.Account(account), model.AccountKey(account)); err != nil {
		return nil, nil, nil, nil, err
	}

	acctPvt := model.NewAccountPrivate()
	if err = json.Unmarshal(bytes, &acctPvt); err != nil {
		return nil, nil, nil, nil, err
	}

	// get account private info from account collection to update available
	if bytes, err = ctx.GetStub().GetState(model.AccountKey(account)); err != nil {
		return nil, nil, nil, nil, err
	}

	acctPub := model.NewAccountPublic()
	if err = json.Unmarshal(bytes, &acctPub); err != nil {
		return nil, nil, nil, nil, err
	}

	return acctPub, acctPvt, assetPub, assetPvt, nil
}
