package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	events "github.com/usnistgov/blossom/chaincode/ngac/epp"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
	decider "github.com/usnistgov/blossom/chaincode/ngac/pdp"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/pkg/errors"
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

func NewLicenseContract() AssetsInterface {
	return &BlossomSmartContract{}
}

func (b *BlossomSmartContract) assetExists(stub shim.ChaincodeStubInterface, id string) (bool, error) {
	data, err := stub.GetPrivateData(CatalogCollection(), model.AssetKey(id))
	if err != nil {
		return false, errors.Wrapf(err, "error checking if asset id %q already exists on the ledger", id)
	}

	return data != nil, nil
}

func (b *BlossomSmartContract) OnboardAsset(stub shim.ChaincodeStubInterface, id string, name string, onboardDate string, expiration string) error {
	if ok, err := b.assetExists(stub, id); err != nil {
		return errors.Wrapf(err, "error checking if asset already exists")
	} else if ok {
		return errors.Errorf("an asset with the ID %q already exists", id)
	}

	assetInput, err := getOnboardAssetTransientInput(stub)
	if err != nil {
		return fmt.Errorf("error getting transient input: %v", err)
	}

	if len(assetInput.Licenses) == 0 {
		return fmt.Errorf("licenses cannot be nil")
	}

	// ngac check
	if err = pdp.CanOnboardAsset(stub, CatalogCollection()); err != nil {
		return errors.Wrapf(err, "ngac check failed")
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
		return errors.Wrapf(err, "error marshaling asset %q", name)
	}

	// put in catalog pdc
	if err = stub.PutPrivateData(CatalogCollection(), model.AssetKey(id), bytes); err != nil {
		return errors.Wrap(err, "error adding asset to catalog private data collection")
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
		return errors.Wrapf(err, "error marshaling asset %q", name)
	}

	// add license to licenses private data
	if err = stub.PutPrivateData(LicensesCollection(), model.AssetKey(id), bytes); err != nil {
		return errors.Wrapf(err, "error adding asset to ledger")
	}

	// ngac event
	return events.ProcessOnboardAsset(stub, CatalogCollection(), id)
}

func (b *BlossomSmartContract) OffboardAsset(stub shim.ChaincodeStubInterface, assetID string) error {
	if ok, err := b.assetExists(stub, assetID); err != nil {
		return errors.Wrapf(err, "error checking if asset exists")
	} else if !ok {
		return nil
	}

	var (
		asset *model.Asset
		err   error
	)

	// ngac check
	if err := pdp.CanOffboardAsset(stub, CatalogCollection()); err != nil {
		return errors.Wrapf(err, "ngac check failed")
	}

	if asset, err = b.GetAsset(stub, assetID); err != nil {
		return errors.Wrapf(err, "error getting asset info")
	}

	// check that all licenses have been returned
	if asset.Available != len(asset.Licenses) {
		return errors.Errorf("asset %s still has licenses checked out", assetID)
	}

	// remove asset from catalog
	if err = stub.DelPrivateData(CatalogCollection(), model.AssetKey(assetID)); err != nil {
		return errors.Wrapf(err, "error offboarding asset from catalog pdc")
	}

	// remove license licenses pdc
	if err = stub.DelPrivateData(LicensesCollection(), model.AssetKey(assetID)); err != nil {
		return errors.Wrapf(err, "error offboarding asset from licenses pdc")
	}

	// ngac event
	return events.ProcessOffboardAsset(stub, CatalogCollection(), assetID)
}

func (b *BlossomSmartContract) GetAssets(stub shim.ChaincodeStubInterface) ([]*model.AssetPublic, error) {
	resultsIterator, err := stub.GetPrivateDataByRange(CatalogCollection(), "", "")
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

		asset := &model.AssetPublic{}
		if err = json.Unmarshal(queryResponse.Value, asset); err != nil {
			return nil, err
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

func (b *BlossomSmartContract) GetAsset(stub shim.ChaincodeStubInterface, id string) (*model.Asset, error) {
	if ok, err := b.assetExists(stub, id); err != nil {
		return nil, errors.Wrapf(err, "error checking if asset exists")
	} else if !ok {
		return nil, errors.Errorf("an asset with the ID %q does not exist", id)
	}

	var (
		assetPub = &model.AssetPublic{}
		assetPvt = &model.AssetPrivate{}
		bytes    []byte
		err      error
	)

	if bytes, err = stub.GetPrivateData(CatalogCollection(), model.AssetKey(id)); err != nil {
		return nil, errors.Wrapf(err, "error getting asset from private data")
	}

	if err = json.Unmarshal(bytes, assetPub); err != nil {
		return nil, errors.Wrapf(err, "error deserializing license")
	}

	if bytes, err = stub.GetPrivateData(LicensesCollection(), model.AssetKey(id)); err != nil {
		// ignore error if a user does not have access to the private data collection of the asset
		// they can still have access to the public info
		mspid, _ := cid.GetMSPID(stub)
		fmt.Printf("error occurred reading pvtdata for user in org %s: %v\n", mspid, err)
	} else {
		if err = json.Unmarshal(bytes, assetPvt); err != nil {
			return nil, errors.Wrapf(err, "error deserializing account private info")
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

func (b *BlossomSmartContract) RequestCheckout(stub shim.ChaincodeStubInterface) error {
	transientInput, err := getRequestCheckoutTransientInput(stub)
	if err != nil {
		return fmt.Errorf("error getting transient input: %v", err)
	}

	var (
		account string
		bytes   []byte
	)

	// check requested asset exists
	if bytes, err = stub.GetPrivateData(CatalogCollection(), model.AssetKey(transientInput.AssetID)); err != nil {
		return err
	} else if bytes == nil {
		return fmt.Errorf("asset with id %s does not exist", transientInput.AssetID)
	}

	if account, err = accountName(stub); err != nil {
		return errors.Wrap(err, "error getting MSPID from stub")
	}

	collection := AccountCollection(account)

	// ngac check
	if err = decider.CanRequestCheckout(stub, collection, account); err != nil {
		return errors.Wrapf(err, "ngac check failed")
	}

	key := checkoutRequestKey(account, transientInput.AssetID)

	// check if request has already been made and not approved
	if bytes, err = stub.GetPrivateData(collection, key); err != nil {
		return err
	} else if bytes != nil {
		return fmt.Errorf("request for asset %s alreadys exists for account %s and has not been approved yet", transientInput.AssetID, account)
	}

	req := &CheckoutRequest{transientInput.AssetID, transientInput.Amount}

	if bytes, err = json.Marshal(req); err != nil {
		return err
	}

	return stub.PutPrivateData(collection, key, bytes)
}

func checkoutRequestKey(account, assetID string) string {
	return fmt.Sprintf("checkout=%s:%s", account, assetID)
}

func (b *BlossomSmartContract) GetCheckoutRequests(stub shim.ChaincodeStubInterface, account string) ([]CheckoutRequest, error) {
	collection := AccountCollection(account)

	key := checkoutRequestKey(account, "")

	iter, err := stub.GetPrivateDataByRange(collection, "", "")
	if err != nil {
		return nil, err
	}

	reqs := make([]CheckoutRequest, 0)
	for iter.HasNext() {
		next, err := iter.Next()
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(next.Key, key) {
			continue
		}

		req := CheckoutRequest{}
		if err = json.Unmarshal(next.Value, &req); err != nil {
			return nil, err
		}

		reqs = append(reqs, req)
	}

	return reqs, nil
}

func (b *BlossomSmartContract) ApproveCheckout(stub shim.ChaincodeStubInterface) error {
	transientInput, err := getApproveCheckoutTransientInput(stub)
	if err != nil {
		return fmt.Errorf("error getting transient input: %v", err)
	}

	var (
		acctColl = AccountCollection(transientInput.Account)
		key      = checkoutRequestKey(transientInput.Account, transientInput.AssetID)
		bytes    []byte
	)

	// ngac check
	if err = decider.CanApproveCheckout(stub, acctColl, transientInput.Account); err != nil {
		return errors.Wrapf(err, "ngac check failed")
	}

	// check that request exists
	if bytes, err = stub.GetPrivateData(acctColl, key); err != nil {
		return errors.Wrapf(err, "error checking if request exists")
	} else if bytes == nil {
		return fmt.Errorf("request for asset %s does not exist for account %s", transientInput.AssetID, transientInput.Account)
	}

	// delete request key
	if err = stub.DelPrivateData(acctColl, key); err != nil {
		return errors.Wrapf(err, "error deleting request")
	}

	req := &CheckoutRequest{}
	if err = json.Unmarshal(bytes, req); err != nil {
		return errors.Wrapf(err, "error unmarshaling request")
	}

	acctPub, acctPvt, assetPub, assetPvt, err := getAcctAndAsset(stub, transientInput.Account, transientInput.AssetID)
	if err != nil {
		return err
	}

	if err = checkout(assetPub, assetPvt, acctPub, acctPvt, req.Amount); err != nil {
		return errors.Wrapf(err, "error checking out %s for account %s", transientInput.AssetID, transientInput.Account)
	}

	return putAcctAndAsset(stub, acctPub, acctPvt, assetPub, assetPvt)
}

func checkout(assetPub *model.AssetPublic, assetPvt *model.AssetPrivate, acctPub *model.AccountPublic, acctPvt *model.AccountPrivate, amount int) error {
	// check that the amount requested is less than the amount available
	if amount > assetPub.Available {
		return errors.Errorf("requested amount %v cannot be greater than the available amount %v",
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

func (b *BlossomSmartContract) GetLicenses(stub shim.ChaincodeStubInterface, account, assetID string) (map[string]string, error) {
	bytes, err := stub.GetPrivateData(AccountCollection(account), model.AccountKey(account))
	if err != nil {
		return nil, errors.Wrapf(err, "error reading account private data")
	}

	acctPvt := &model.AccountPrivate{}
	if err = json.Unmarshal(bytes, acctPvt); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling account private data")
	}

	return acctPvt.Assets[assetID], nil
}

func (b *BlossomSmartContract) InitiateCheckin(stub shim.ChaincodeStubInterface) error {
	transientInput, err := getInitiateCheckinTransientInput(stub)
	if err != nil {
		return fmt.Errorf("error getting transient input: %v", err)
	}

	account, err := accountName(stub)
	if err != nil {
		return errors.Wrapf(err, "error getting MSPID from stub")
	}

	collection := AccountCollection(account)

	// ngac check
	if err = decider.CanInitiateCheckIn(stub, collection, account); err != nil {
		return errors.Wrapf(err, "ngac check failed")
	}

	var (
		key   = checkinRequestKey(account, transientInput.AssetID)
		bytes []byte
	)

	// check if the licenses in the request are really checked out by the account
	if bytes, err = stub.GetPrivateData(AccountCollection(account), model.AccountKey(account)); err != nil {
		return fmt.Errorf("error getting account private info from private data: %v", err)
	}

	acctPvt := &model.AccountPrivate{}
	if err = json.Unmarshal(bytes, &acctPvt); err != nil {
		return fmt.Errorf("error unmarshaling account private info: %v", err)
	}

	checkedOut := acctPvt.Assets[transientInput.AssetID]
	for _, returnedKey := range transientInput.Licenses {
		// check that the returned license is leased to the account
		if _, ok := checkedOut[returnedKey]; !ok {
			return errors.Errorf("returned key %s was not checked out by %s", returnedKey, account)
		}
	}

	// check if request has already been made and not approved
	if bytes, err = stub.GetPrivateData(collection, key); err != nil {
		return err
	} else if bytes != nil {
		return fmt.Errorf("request to checkin %s has already been initiated for account %s and has not been processed yet", transientInput.AssetID, account)
	}

	req := CheckinRequest{
		Asset:    transientInput.AssetID,
		Licenses: transientInput.Licenses,
	}

	if bytes, err = json.Marshal(req); err != nil {
		return err
	}

	return stub.PutPrivateData(collection, key, bytes)
}

func (b *BlossomSmartContract) GetInitiatedCheckins(stub shim.ChaincodeStubInterface, account string) ([]CheckinRequest, error) {
	collection := AccountCollection(account)

	key := checkinRequestKey(account, "")

	iter, err := stub.GetPrivateDataByRange(collection, "", "")
	if err != nil {
		return nil, err
	}

	reqs := make([]CheckinRequest, 0)
	for iter.HasNext() {
		next, err := iter.Next()
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(next.Key, key) {
			continue
		}

		req := CheckinRequest{}
		if err = json.Unmarshal(next.Value, &req); err != nil {
			return nil, err
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
		return errors.Errorf("account %s has not checked out any licenses for asset %s", acctPub.Name, assetPub.ID)
	}

	for _, license := range licenses {
		// check that the account has the license checked out
		if _, ok = accountCheckedOut[license]; !ok {
			return errors.Errorf("returned license %s was not checked out by %s", license, acctPub.Name)
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

func (b *BlossomSmartContract) ProcessCheckin(stub shim.ChaincodeStubInterface) error {
	transientInput, err := getProcessCheckinTransientInput(stub)
	if err != nil {
		return fmt.Errorf("error getting transient input: %v", err)
	}

	var (
		acctColl = AccountCollection(transientInput.Account)
		key      = checkinRequestKey(transientInput.Account, transientInput.AssetID)
		bytes    []byte
	)

	// ngac check
	if err = decider.CanProcessCheckIn(stub, acctColl, transientInput.Account); err != nil {
		return errors.Wrapf(err, "ngac check failed")
	}

	// check that request exists
	if bytes, err = stub.GetPrivateData(acctColl, key); err != nil {
		return errors.Wrapf(err, "error checking if checkin request exists")
	} else if bytes == nil {
		return fmt.Errorf("request to checkin asset %s does not exist for account %s", transientInput.AssetID, transientInput.Account)
	}

	// delete request key
	if err = stub.DelPrivateData(acctColl, key); err != nil {
		return errors.Wrapf(err, "error deleting request")
	}

	req := &CheckinRequest{}
	if err = json.Unmarshal(bytes, req); err != nil {
		return errors.Wrapf(err, "error unmarshaling request")
	}

	acctPub, acctPvt, assetPub, assetPvt, err := getAcctAndAsset(stub, transientInput.Account, transientInput.AssetID)
	if err != nil {
		return err
	}

	if err = checkin(assetPub, assetPvt, acctPub, acctPvt, req.Licenses); err != nil {
		return errors.Wrapf(err, "error checking out %s for account %s", transientInput.AssetID, transientInput.Account)
	}

	return putAcctAndAsset(stub, acctPub, acctPvt, assetPub, assetPvt)
}

func putAcctAndAsset(stub shim.ChaincodeStubInterface, acctPub *model.AccountPublic, acctPvt *model.AccountPrivate,
	assetPub *model.AssetPublic, assetPvt *model.AssetPrivate) (err error) {
	var (
		bytes    []byte
		acctKey  = model.AccountKey(acctPub.Name)
		acctColl = AccountCollection(acctPub.Name)
	)

	// put account public
	if bytes, err = json.Marshal(acctPub); err != nil {
		return
	}

	if err = stub.PutState(acctKey, bytes); err != nil {
		return
	}

	// put account private
	if bytes, err = json.Marshal(acctPvt); err != nil {
		return
	}

	if err = stub.PutPrivateData(acctColl, acctKey, bytes); err != nil {
		return
	}

	// put asset public (still pdc)
	if bytes, err = json.Marshal(assetPub); err != nil {
		return
	}

	if err = stub.PutPrivateData(CatalogCollection(), model.AssetKey(assetPub.ID), bytes); err != nil {
		return
	}

	// put asset private
	if bytes, err = json.Marshal(assetPvt); err != nil {
		return
	}

	return stub.PutPrivateData(LicensesCollection(), model.AssetKey(assetPub.ID), bytes)
}

func getAcctAndAsset(stub shim.ChaincodeStubInterface, account, assetID string) (*model.AccountPublic, *model.AccountPrivate, *model.AssetPublic, *model.AssetPrivate, error) {
	var (
		bytes []byte
		err   error
	)

	// get licenses from license collection
	if bytes, err = stub.GetPrivateData(LicensesCollection(), model.AssetKey(assetID)); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error getting license info from private data")
	}

	assetPvt := &model.AssetPrivate{}
	if err = json.Unmarshal(bytes, &assetPvt); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error unmarshaling asset private info")
	}

	// get asset public info from catalog collection to update available
	if bytes, err = stub.GetPrivateData(CatalogCollection(), model.AssetKey(assetID)); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error getting asset public info from private data")
	}

	assetPub := &model.AssetPublic{}
	if err = json.Unmarshal(bytes, &assetPub); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error unmarshaling asset private info")
	}

	// get account private info from account collection to update available
	if bytes, err = stub.GetPrivateData(AccountCollection(account), model.AccountKey(account)); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error getting account private info from private data")
	}

	acctPvt := &model.AccountPrivate{}
	if err = json.Unmarshal(bytes, &acctPvt); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error unmarshaling account private info")
	}

	// get account private info from account collection to update available
	if bytes, err = stub.GetState(model.AccountKey(account)); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error getting account public info from private data")
	}

	acctPub := &model.AccountPublic{}
	if err = json.Unmarshal(bytes, &acctPub); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error unmarshaling account public info")
	}

	return acctPub, acctPvt, assetPub, assetPvt, nil
}
