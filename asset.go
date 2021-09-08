package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
)

type (
	// AssetsInterface provides the functions to interact with Assets in fabric.
	AssetsInterface interface {
		// OnboardAsset adds a new software asset to Blossom.  This will create a new asset object on the ledger and in the
		// NGAC graph. Assets are identified by the ID field. The user performing the request will need to
		// have permission to add an asset to the ledger. The asset will be an object attribute in NGAC and the
		// asset licenses will be objects that are assigned to the asset.
		OnboardAsset(stub shim.ChaincodeStubInterface, asset *model.Asset) error

		// OffboardAsset removes an existing asset in Blossom.  This will remove the license from the ledger
		// and from NGAC. An error will be returned if there are any accounts that have checked out the asset
		// and the licenses are not returned
		OffboardAsset(stub shim.ChaincodeStubInterface, id string) error

		// Assets returns all software assets in Blossom. This information includes which accounts have licenses for each
		// asset.
		Assets(stub shim.ChaincodeStubInterface) ([]*model.Asset, error)

		// AssetInfo returns the info for the asset with the given asset ID.
		AssetInfo(stub shim.ChaincodeStubInterface, id string) (*model.Asset, error)

		// Checkout requests software licenses for an account.  The requesting user must have permission to request
		// (i.e. System Administrator). The amount parameter is the amount of software licenses the account is requesting.
		// This number is subtracted from the total available for the asset. Returns the set of licenses that are now assigned to
		// the account.
		Checkout(stub shim.ChaincodeStubInterface, assetID string, account string, amount int) (map[string]model.DateTime, error)

		// Checkin returns the licenses to Blossom.  The return of these licenses is reflected in the amount available for
		// the asset, and the licenses assigned to the account on the ledger.
		Checkin(stub shim.ChaincodeStubInterface, assetID string, licenses []string, accountName string) error
	}
)

func NewLicenseContract() AssetsInterface {
	return &BlossomSmartContract{}
}

func (b *BlossomSmartContract) assetExists(stub shim.ChaincodeStubInterface, id string) (bool, error) {
	data, err := stub.GetState(model.AssetKey(id))
	if err != nil {
		return false, errors.Wrapf(err, "error checking if asset id %q already exists on the ledger", id)
	}

	return data != nil, nil
}

func (b *BlossomSmartContract) OnboardAsset(stub shim.ChaincodeStubInterface, asset *model.Asset) error {
	if ok, err := b.assetExists(stub, asset.ID); err != nil {
		return errors.Wrapf(err, "error checking if asset already exists")
	} else if ok {
		return errors.Errorf("an asset with the ID %q already exists", asset.ID)
	}

	// begin NGAC
	if err := pdp.NewAssetDecider().OnboardAsset(stub, asset); err != nil {
		return errors.Wrapf(err, "error onboarding asset %q", asset.Name)
	}
	// end NGAC

	// at the time of onboarding all licenses are available
	asset.AvailableLicenses = asset.Licenses
	asset.OnboardingDate = model.DateTime(time.Now().String())
	asset.CheckedOut = make(map[string]map[string]model.DateTime)

	// convert license to bytes
	bytes, err := json.Marshal(asset)
	if err != nil {
		return errors.Wrapf(err, "error marshaling asset %q", asset.Name)
	}

	// add license to world state
	if err = stub.PutState(model.AssetKey(asset.ID), bytes); err != nil {
		return errors.Wrapf(err, "error adding asset to ledger")
	}

	return nil
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

	if asset, err = b.AssetInfo(stub, assetID); err != nil {
		return errors.Wrapf(err, "error getting asset info")
	}

	// check that all licenses have been returned
	if len(asset.CheckedOut) != 0 {
		return errors.Errorf("asset %s still has licenses checked out", assetID)
	}

	// begin NGAC
	if err = pdp.NewAssetDecider().OffboardAsset(stub, assetID); err != nil {
		return errors.Wrapf(err, "error offboarding asset %q in NGAC", assetID)
	}
	// end NGAC

	// remove license from world state
	if err = stub.DelState(model.AssetKey(assetID)); err != nil {
		return errors.Wrapf(err, "error offboarding asset from ledger")
	}

	return nil
}

func (b *BlossomSmartContract) Assets(stub shim.ChaincodeStubInterface) ([]*model.Asset, error) {
	// retrieve the assets from the ledger
	assets, err := assets(stub)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting licenses")
	}

	// begin NGAC
	// filter any asset information the requesting user may not have permission to see
	if assets, err = pdp.NewAssetDecider().FilterAssets(stub, assets); err != nil {
		return nil, errors.Wrapf(err, "error filtering assets")
	}
	// end NGAC

	return assets, nil
}

func assets(stub shim.ChaincodeStubInterface) ([]*model.Asset, error) {
	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	assets := make([]*model.Asset, 0)
	for resultsIterator.HasNext() {
		var queryResponse *queryresult.KV
		if queryResponse, err = resultsIterator.Next(); err != nil {
			return nil, err
		}

		// assets on the ledger begin with the asset prefix -- ignore other results
		if !strings.HasPrefix(queryResponse.Key, model.AssetPrefix) {
			continue
		}

		asset := &model.Asset{}
		if err = json.Unmarshal(queryResponse.Value, asset); err != nil {
			return nil, err
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

func (b *BlossomSmartContract) AssetInfo(stub shim.ChaincodeStubInterface, id string) (*model.Asset, error) {
	asset, err := b.assetInfo(stub, id)
	if err != nil {
		return nil, errors.Wrapf(err, "error retrieving asset info")
	}

	// begin NGAC
	// filter any asset information the requesting user may not have permission to see
	if err := pdp.NewAssetDecider().FilterAsset(stub, asset); err != nil {
		return nil, errors.Wrapf(err, "error filtering asset")
	}
	// end NGAC

	return asset, nil
}

func (b *BlossomSmartContract) Checkout(
	stub shim.ChaincodeStubInterface,
	assetID string,
	accountName string,
	amount int) (map[string]model.DateTime, error) {

	var (
		asset   = &model.Asset{}
		account = &model.Account{}
		err     error
	)

	// get the account that will be leasing the licenses
	if account, err = b.Account(stub, accountName); err != nil {
		return nil, errors.Wrapf(err, "error getting account %q", accountName)
	}

	// get asset being requested
	if asset, err = b.assetInfo(stub, assetID); err != nil {
		return nil, errors.Wrapf(err, "error getting info for asset %q", assetID)
	}

	// checkout the asset
	var checkedOutLicenses map[string]model.DateTime
	if checkedOutLicenses, err = checkout(account, asset, amount); err != nil {
		return nil, errors.Wrapf(err, "error checking out %q", asset.ID)
	}

	// update account's record of checked out licenses
	var bytes []byte
	if bytes, err = json.Marshal(account); err != nil {
		return nil, errors.Wrapf(err, "error marshaling account %q", account.Name)
	}

	if err = stub.PutState(model.AccountKey(account.Name), bytes); err != nil {
		return nil, errors.Wrapf(err, "error updating account state")
	}

	// update asset to reflect the licenses being leased to the account
	if bytes, err = json.Marshal(asset); err != nil {
		return nil, errors.Wrapf(err, "error marshaling asset %q", asset.ID)
	}

	if err = stub.PutState(model.AssetKey(asset.ID), bytes); err != nil {
		return nil, errors.Wrapf(err, "error updating asset state")
	}

	// begin NGAC
	// record the checkout in NGAC
	// provide NGAC with the licenses that were checked out in order to reflect the change in the graph
	// this change will provide the users of the requesting account access to the licenses, nobody else
	// will be able to access them
	if err := pdp.NewAssetDecider().Checkout(stub, accountName, assetID, checkedOutLicenses); err != nil {
		return nil, errors.Wrapf(err, "error checking out asset in NGAC")
	}
	// end NGAC

	return checkedOutLicenses, nil
}

func (b *BlossomSmartContract) assetInfo(stub shim.ChaincodeStubInterface, id string) (*model.Asset, error) {
	if ok, err := b.assetExists(stub, id); err != nil {
		return nil, errors.Wrapf(err, "error checking if asset exists")
	} else if !ok {
		return nil, errors.Errorf("an asset with the ID %q does not exist", id)
	}

	var (
		asset = &model.Asset{}
		bytes []byte
		err   error
	)

	if bytes, err = stub.GetState(model.AssetKey(id)); err != nil {
		return nil, errors.Wrapf(err, "error getting asset from ledger")
	}

	if err = json.Unmarshal(bytes, asset); err != nil {
		return nil, errors.Wrapf(err, "error deserializing license")
	}

	return asset, nil
}

func checkout(account *model.Account, asset *model.Asset, amount int) (map[string]model.DateTime, error) {
	// check that the amount requested is less than the amount available
	if amount > asset.Available {
		return nil, errors.Errorf("requested amount %v cannot be greater than the available amount %v",
			amount, asset.Available)
	}

	// update available amount
	asset.Available -= amount

	// get the available licenses
	fromAvailable := asset.AvailableLicenses[0:amount]
	// update available licenses
	asset.AvailableLicenses = asset.AvailableLicenses[amount:]

	// create the array of licenses that are checked out including expiration dates
	retCheckedOutLicenses := make(map[string]model.DateTime)
	expiration := model.DateTime(time.Now().AddDate(1, 0, 0).String())
	for _, license := range fromAvailable {
		// set the expiration of the license to one year from now
		retCheckedOutLicenses[license] = expiration
	}

	// update the account assets
	// add to existing asset if they are checking out more of a software asset
	allCheckedOutAssets, ok := account.Assets[asset.ID]
	if ok {
		allCheckedOutAssets = retCheckedOutLicenses
	} else {
		allCheckedOutAssets = make(map[string]model.DateTime)
		for license, exp := range retCheckedOutLicenses {
			allCheckedOutAssets[license] = exp
		}
	}

	// update asset in the account
	account.Assets[asset.ID] = allCheckedOutAssets

	// update the asset's account tracker
	accountCheckedOut := make(map[string]model.DateTime)
	for k, t := range allCheckedOutAssets {
		accountCheckedOut[k] = t
	}
	asset.CheckedOut[account.Name] = accountCheckedOut

	return retCheckedOutLicenses, nil
}

func (b *BlossomSmartContract) Checkin(stub shim.ChaincodeStubInterface, assetID string, licenses []string, accountName string) error {
	var (
		asset   = &model.Asset{}
		account = &model.Account{}
		err     error
	)

	// get account
	if account, err = b.Account(stub, accountName); err != nil {
		return errors.Wrapf(err, "error getting account %q", accountName)
	}

	// get asset
	if asset, err = b.assetInfo(stub, assetID); err != nil {
		return errors.Wrapf(err, "error getting info for asset %q", assetID)
	}

	// check in asset logic
	if err = checkin(account, asset, licenses); err != nil {
		return err
	}

	// update account
	var bytes []byte
	if bytes, err = json.Marshal(account); err != nil {
		return errors.Wrapf(err, "error marshaling account %q", account.Name)
	}

	if err = stub.PutState(model.AccountKey(account.Name), bytes); err != nil {
		return errors.Wrapf(err, "error updating account state")
	}

	// update asset
	if bytes, err = json.Marshal(asset); err != nil {
		return errors.Wrapf(err, "error marshaling asset %q", asset.ID)
	}

	if err = stub.PutState(model.AssetKey(asset.Name), bytes); err != nil {
		return errors.Wrapf(err, "error updating asset state")
	}

	// begin NGAC
	// record the checkin in NGAC
	// provide NGAC with the licenses that were checked in in order to reflect the change in the graph
	// this will move the licenses back into the pool of available licenses
	// the account users will no longer be able to see the licenses
	if err := pdp.NewAssetDecider().Checkin(stub, accountName, assetID, licenses); err != nil {
		return errors.Wrapf(err, "error checking in licenses in NGAC")
	}
	// end NGAC

	return nil
}

func checkin(account *model.Account, asset *model.Asset, licenses []string) error {
	checkedOut := account.Assets[asset.ID]
	for _, returnedKey := range licenses {
		// check that the returned license is leased to the account
		if _, ok := account.Assets[asset.ID][returnedKey]; !ok {
			return errors.Errorf("returned key %s was not checked out by %s", returnedKey, account.Name)
		}

		delete(checkedOut, returnedKey)
	}

	// if all licenses were returned remove asset from account's checked out
	if len(checkedOut) == 0 {
		delete(account.Assets, asset.ID)
	} else {
		// update account licenses
		account.Assets[asset.ID] = checkedOut
	}

	accountCheckedOut, ok := asset.CheckedOut[account.Name]
	if !ok {
		return errors.Errorf("account %s has not checked out any licenses for asset %s", account.Name, asset.ID)
	}

	for _, license := range licenses {
		// check that the account has the license checked out
		if _, ok = accountCheckedOut[license]; !ok {
			return errors.Errorf("returned license %s was not checked out by %s", license, account.Name)
		}

		// remove the returned license from the checked out licenses
		delete(accountCheckedOut, license)

		// add the returned license to the available licenses
		asset.AvailableLicenses = append(asset.AvailableLicenses, license)
	}

	// if all licenses are returned, remove the account from the asset
	if len(accountCheckedOut) == 0 {
		delete(asset.CheckedOut, account.Name)
	}

	// update number of available licenses
	asset.Available += len(licenses)

	return nil
}
