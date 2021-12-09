package main

import (
	"encoding/json"
	"fmt"
	events "github.com/usnistgov/blossom/chaincode/ngac/epp"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
	decider "github.com/usnistgov/blossom/chaincode/ngac/pdp"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
)

type (
	// AssetsInterface provides the functions to interact with Assets in fabric.
	AssetsInterface interface {
		// OnboardAsset adds a new software asset to Blossom.  This will create a new asset object on the ledger and in the
		// NGAC graph. Assets are identified by the ID field. The user performing the request will need to
		// have permission to add an asset to the ledger. The asset will be an object attribute in NGAC and the
		// asset licenses will be objects that are assigned to the asset.
		// TRANSIENT MAP: export ATO=$(echo -n "{\"licenses\":\"\"}" | base64 | tr -d \\n)
		OnboardAsset(stub shim.ChaincodeStubInterface, id string, name string, expiration model.DateTime) error

		// OffboardAsset removes an existing asset in Blossom.  This will remove the license from the ledger
		// and from NGAC. An error will be returned if there are any accounts that have checked out the asset
		// and the licenses are not returned
		OffboardAsset(stub shim.ChaincodeStubInterface, id string) error

		// Assets returns all software assets in Blossom. This information includes which accounts have licenses for each
		// asset.
		Assets(stub shim.ChaincodeStubInterface) ([]*model.AssetPublic, error)

		// AssetInfo returns the info for the asset with the given asset ID.
		AssetInfo(stub shim.ChaincodeStubInterface, id string) (*model.Asset, error)

		// RequestCheckout requests software licenses for an account.  The requesting user must have permission to request
		// (i.e. System Administrator). The amount parameter is the amount of software licenses the account is requesting.
		// This number is subtracted from the total available for the asset. Returns the set of licenses that are now assigned to
		// the account.
		RequestCheckout(stub shim.ChaincodeStubInterface) error

		// ApproveCheckout approves a checkout request made by an account.  The requested licenses for the asset will be
		// added to the account's private data collection. A user on the account can then call Licenses to get the approved
		// license keys.
		ApproveCheckout(stub shim.ChaincodeStubInterface) error

		// Licenses get the license keys for an asset that an account has access to in their private data collection
		Licenses(stub shim.ChaincodeStubInterface, account, assetID string) (map[string]model.DateTime, error)

		// InitiateCheckin starts the process of returning licenses to Blossom. This is serves as a request to the blossom
		// admin to process the return of the licenses. This is because only the blossom admin can write to the licenses
		// private data collection to return the licenses to the available pool.
		InitiateCheckin(stub shim.ChaincodeStubInterface) error

		// ProcessCheckin processes an account's checkin request (from InitiateCheckin) and returns the licenses to the
		// available pool in the licenses private data collection.
		ProcessCheckin(stub shim.ChaincodeStubInterface) error
	}

	checkoutRequest struct {
		Asset  string `json:"asset,omitempty"`
		Amount int    `json:"amount,omitempty"`
	}

	checkinRequest struct {
		Asset    string   `json:"asset,omitempty"`
		Licenses []string `json:"licenses,omitempty"`
	}
)

func NewLicenseContract() AssetsInterface {
	return &BlossomSmartContract{}
}

func (b *BlossomSmartContract) assetExists(stub shim.ChaincodeStubInterface, id string) (bool, error) {
	data, err := stub.GetPrivateData(CatalogCollectionName(), model.AssetKey(id))
	if err != nil {
		return false, errors.Wrapf(err, "error checking if asset id %q already exists on the ledger", id)
	}

	return data != nil, nil
}

func (b *BlossomSmartContract) OnboardAsset(stub shim.ChaincodeStubInterface, id string, name string, expiration model.DateTime) error {
	if ok, err := b.assetExists(stub, id); err != nil {
		return errors.Wrapf(err, "error checking if asset already exists")
	} else if ok {
		return errors.Errorf("an asset with the ID %q already exists", id)
	}

	assetInput, err := getOnboardAssetTransientInput(stub)
	if err != nil {
		return err
	}

	if len(assetInput.Licenses) == 0 {
		return fmt.Errorf("licenses cannot be nil")
	}

	// ngac check
	if err = pdp.CanOnboardAsset(stub, CatalogCollectionName()); err != nil {
		return errors.Wrapf(err, "ngac check failed")
	}

	// public info - id, name, available (=total), expiration
	assetPub := &model.AssetPublic{
		ID:             id,
		Name:           name,
		Available:      len(assetInput.Licenses),
		OnboardingDate: model.DateTime(time.Now().String()),
		Expiration:     expiration,
	}

	bytes, err := json.Marshal(assetPub)
	if err != nil {
		return errors.Wrapf(err, "error marshaling asset %q", name)
	}

	// put in catalog pdc
	if err = stub.PutPrivateData(CatalogCollectionName(), model.AssetKey(id), bytes); err != nil {
		return errors.Wrap(err, "error adding asset to catalog private data collection")
	}

	assetPvt := model.AssetPrivate{
		TotalAmount:       len(assetInput.Licenses),
		Licenses:          assetInput.Licenses,
		AvailableLicenses: assetInput.Licenses,
		CheckedOut:        make(map[string]map[string]model.DateTime),
	}

	if bytes, err = json.Marshal(assetPvt); err != nil {
		return errors.Wrapf(err, "error marshaling asset %q", name)
	}

	// add license to licenses private data
	if err = stub.PutPrivateData(LicensesCollectionName(), model.AssetKey(id), bytes); err != nil {
		return errors.Wrapf(err, "error adding asset to ledger")
	}

	// ngac event
	return events.ProcessOnboardAsset(stub, LicensesCollectionName(), id)
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
	if err := pdp.CanOffboardAsset(stub, CatalogCollectionName()); err != nil {
		return errors.Wrapf(err, "ngac check failed")
	}

	if asset, err = b.AssetInfo(stub, assetID); err != nil {
		return errors.Wrapf(err, "error getting asset info")
	}

	// check that all licenses have been returned
	if len(asset.CheckedOut) != 0 {
		return errors.Errorf("asset %s still has licenses checked out", assetID)
	}

	// remove asset from catalog
	if err = stub.DelPrivateData(CatalogCollectionName(), model.AssetKey(assetID)); err != nil {
		return errors.Wrapf(err, "error offboarding asset from catalog pdc")
	}

	// remove license licenses pdc
	if err = stub.DelPrivateData(LicensesCollectionName(), model.AssetKey(assetID)); err != nil {
		return errors.Wrapf(err, "error offboarding asset from licenses pdc")
	}

	// ngac event
	return events.ProcessOffboardAsset(stub, CatalogCollectionName(), assetID)
}

func (b *BlossomSmartContract) Assets(stub shim.ChaincodeStubInterface) ([]*model.AssetPublic, error) {
	resultsIterator, err := stub.GetPrivateDataByRange(CatalogCollectionName(), "", "")
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

func (b *BlossomSmartContract) AssetInfo(stub shim.ChaincodeStubInterface, id string) (*model.Asset, error) {
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

	if bytes, err = stub.GetPrivateData(CatalogCollectionName(), model.AssetKey(id)); err != nil {
		return nil, errors.Wrapf(err, "error getting asset from private data")
	}

	if err = json.Unmarshal(bytes, assetPub); err != nil {
		return nil, errors.Wrapf(err, "error deserializing license")
	}

	if bytes, err = stub.GetPrivateData(LicensesCollectionName(), model.AssetKey(id)); err != nil {
		// ignore error if a user does not have access to the private data collection of the asset
		// they can still have access to the public info
		fmt.Printf("error occurred reading pvtdata: %v\n", err)
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
		return fmt.Errorf("error getting transient input")
	}

	var (
		account string
		bytes   []byte
	)

	if account, err = accountName(stub); err != nil {
		return errors.Wrap(err, "error getting MSPID from stub")
	}

	collection := AccountCollectionName(account)

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

	req := &checkoutRequest{transientInput.AssetID, transientInput.Amount}

	if bytes, err = json.Marshal(req); err != nil {
		return err
	}

	return stub.PutPrivateData(collection, key, bytes)
}

func checkoutRequestKey(account, assetID string) string {
	return fmt.Sprintf("checkout=%s:%s", account, assetID)
}

func (b *BlossomSmartContract) ApproveCheckout(stub shim.ChaincodeStubInterface) error {
	transientInput, err := getApproveCheckoutTransientInput(stub)
	if err != nil {
		return fmt.Errorf("error getting transient input")
	}

	var (
		acctColl = AccountCollectionName(transientInput.Account)
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

	req := &checkoutRequest{}
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

	// create the array of licenses that are checked out including expiration dates
	retCheckedOutLicenses := make(map[string]model.DateTime)
	expiration := model.DateTime(time.Now().AddDate(1, 0, 0).String())
	for _, license := range fromAvailable {
		// set the expiration of the license to one year from now
		retCheckedOutLicenses[license] = expiration
	}

	// update the account assets
	// add to existing asset if they are checking out more of a software asset
	allCheckedOutAssets, ok := acctPvt.Assets[assetPub.ID]
	if ok {
		allCheckedOutAssets = retCheckedOutLicenses
	} else {
		allCheckedOutAssets = make(map[string]model.DateTime)
		for license, exp := range retCheckedOutLicenses {
			allCheckedOutAssets[license] = exp
		}
	}

	// update asset in the account
	acctPvt.Assets[assetPub.ID] = allCheckedOutAssets

	// update the asset's account tracker
	accountCheckedOut := make(map[string]model.DateTime)
	for k, t := range allCheckedOutAssets {
		accountCheckedOut[k] = t
	}
	assetPvt.CheckedOut[acctPub.Name] = accountCheckedOut

	return nil
}

func (b *BlossomSmartContract) Licenses(stub shim.ChaincodeStubInterface, account, assetID string) (map[string]model.DateTime, error) {
	bytes, err := stub.GetPrivateData(AccountCollectionName(account), model.AccountKey(account))
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
		return fmt.Errorf("error getting transient input")
	}

	account, err := accountName(stub)
	if err != nil {
		return errors.Wrapf(err, "error getting MSPID from stub")
	}

	collection := AccountCollectionName(account)

	// ngac check
	if err = decider.CanInitiateCheckIn(stub, collection, account); err != nil {
		return errors.Wrapf(err, "ngac check failed")
	}

	var (
		key   = checkinRequestKey(account, transientInput.AssetID)
		bytes []byte
	)

	// check if request has already been made and not approved
	if bytes, err = stub.GetPrivateData(collection, key); err != nil {
		return err
	} else if bytes != nil {
		return fmt.Errorf("request to checkin %s has already been initiated for account %s and has not been processed yet", transientInput.AssetID, account)
	}

	req := checkinRequest{
		Asset:    transientInput.AssetID,
		Licenses: transientInput.Licenses,
	}

	if bytes, err = json.Marshal(req); err != nil {
		return err
	}

	return stub.PutPrivateData(collection, key, bytes)
}

func checkinRequestKey(account, assetID string) string {
	return fmt.Sprintf("checkin=%s:%s", account, assetID)
}

func checkin(assetPub *model.AssetPublic, assetPvt *model.AssetPrivate, acctPub *model.AccountPublic, acctPvt *model.AccountPrivate, licenses []string) error {
	checkedOut := acctPvt.Assets[assetPub.ID]
	for _, returnedKey := range licenses {
		// check that the returned license is leased to the account
		if _, ok := checkedOut[returnedKey]; !ok {
			return errors.Errorf("returned key %s was not checked out by %s", returnedKey, acctPub.Name)
		}

		delete(checkedOut, returnedKey)
	}

	// if all licenses were returned remove asset from account's checked out
	if len(checkedOut) == 0 {
		delete(acctPvt.Assets, assetPub.ID)
	} else {
		// update account licenses
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
		return fmt.Errorf("error getting transient input")
	}

	var (
		acctColl = AccountCollectionName(transientInput.Account)
		key      = checkoutRequestKey(transientInput.Account, transientInput.AssetID)
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

	req := &checkinRequest{}
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
		acctColl = AccountCollectionName(acctPub.Name)
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

	if err = stub.PutPrivateData(CatalogCollectionName(), model.AssetKey(assetPub.ID), bytes); err != nil {
		return
	}

	// put asset private
	if bytes, err = json.Marshal(assetPvt); err != nil {
		return
	}

	return stub.PutPrivateData(LicensesCollectionName(), model.AssetKey(assetPub.ID), bytes)
}

func getAcctAndAsset(stub shim.ChaincodeStubInterface, account, assetID string) (*model.AccountPublic, *model.AccountPrivate, *model.AssetPublic, *model.AssetPrivate, error) {
	var (
		bytes []byte
		err   error
	)

	// get licenses from license collection
	if bytes, err = stub.GetPrivateData(LicensesCollectionName(), model.AssetKey(assetID)); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error getting license info from private data")
	}

	assetPvt := &model.AssetPrivate{}
	if err = json.Unmarshal(bytes, &assetPvt); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error unmarshaling asset private info")
	}

	// get asset public info from catalog collection to update available
	if bytes, err = stub.GetPrivateData(CatalogCollectionName(), model.AssetKey(assetID)); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error getting asset public info from private data")
	}

	assetPub := &model.AssetPublic{}
	if err = json.Unmarshal(bytes, &assetPub); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error unmarshaling asset private info")
	}

	// get account private info from account collection to update available
	if bytes, err = stub.GetPrivateData(AccountCollectionName(account), model.AccountKey(account)); err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "error getting account pruvate info from private data")
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
