package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
	"strings"
	"time"
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
		// and from NGAC. An error will be returned if there are any agencies that have checked out the asset
		// and the licenses are not returned
		OffboardAsset(stub shim.ChaincodeStubInterface, id string) error

		// Assets returns all software assets in Blossom. This information includes which agencies have licenses for each
		// asset.
		Assets(stub shim.ChaincodeStubInterface) ([]*model.Asset, error)

		// AssetInfo returns the info for the asset with the given asset ID.
		AssetInfo(stub shim.ChaincodeStubInterface, id string) (*model.Asset, error)

		// Checkout requests software licenses for an agency.  The requesting user must have permission to request
		// (i.e. System Administrator). The amount parameter is the amount of software licenses the agency is requesting.
		// This number is subtracted from the total available for the asset. Returns the set of licenses that are now assigned to
		// the agency.
		Checkout(stub shim.ChaincodeStubInterface, assetID string, agency string, amount int) (map[string]time.Time, error)

		// Checkin returns the licenses to Blossom.  The return of these licenses is reflected in the amount available for
		// the asset, and the licenses assigned to the agency on the ledger.
		Checkin(stub shim.ChaincodeStubInterface, assetID string, licenses []string, agencyName string) error
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
	asset.OnboardingDate = time.Now()
	asset.CheckedOut = make(map[string]map[string]time.Time)

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
	agencyName string,
	amount int) (map[string]time.Time, error) {

	var (
		asset  = &model.Asset{}
		agency = &model.Agency{}
		err    error
	)

	// get the agency that will be leasing the licenses
	if agency, err = b.Agency(stub, agencyName); err != nil {
		return nil, errors.Wrapf(err, "error getting agency %q", agencyName)
	}

	// get asset being requested
	if asset, err = b.AssetInfo(stub, assetID); err != nil {
		return nil, errors.Wrapf(err, "error getting info for asset %q", assetID)
	}

	// checkout the asset
	var checkedOutLicenses map[string]time.Time
	if checkedOutLicenses, err = checkout(agency, asset, amount); err != nil {
		return nil, errors.Wrapf(err, "error checking out %q", asset.ID)
	}

	// update agency's record of checked out licenses
	var bytes []byte
	if bytes, err = json.Marshal(agency); err != nil {
		return nil, errors.Wrapf(err, "error marshaling agency %q", agency.Name)
	}

	if err = stub.PutState(model.AgencyKey(agency.Name), bytes); err != nil {
		return nil, errors.Wrapf(err, "error updating agency state")
	}

	// update asset to reflect the licenses being leased to the agency
	if bytes, err = json.Marshal(asset); err != nil {
		return nil, errors.Wrapf(err, "error marshaling asset %q", asset.ID)
	}

	if err = stub.PutState(model.AssetKey(asset.ID), bytes); err != nil {
		return nil, errors.Wrapf(err, "error updating asset state")
	}

	// begin NGAC
	// record the checkout in NGAC
	// provide NGAC with the licenses that were checked out in order to reflect the change in the graph
	// this change will provide the users of the requesting agency access to the licenses, nobody else
	// will be able to access them
	if err := pdp.NewAssetDecider().Checkout(stub, agencyName, assetID, checkedOutLicenses); err != nil {
		return nil, errors.Wrapf(err, "error checking out asset in NGAC")
	}
	// end NGAC

	return checkedOutLicenses, nil
}

func checkout(agency *model.Agency, asset *model.Asset, amount int) (map[string]time.Time, error) {
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
	retCheckedOutLicenses := make(map[string]time.Time, 0)
	expiration := time.Now().AddDate(1, 0, 0)
	for _, license := range fromAvailable {
		// set the expiration of the license to one year from now
		retCheckedOutLicenses[license] = expiration
	}

	// update the agency assets
	// add to existing asset if they are checking out more of a software asset
	allCheckedOutAssets, ok := agency.Assets[asset.ID]
	if ok {
		allCheckedOutAssets = retCheckedOutLicenses
	} else {
		allCheckedOutAssets = make(map[string]time.Time)
		for license, exp := range retCheckedOutLicenses {
			allCheckedOutAssets[license] = exp
		}
	}

	// update asset in the agency
	agency.Assets[asset.ID] = allCheckedOutAssets

	// update the asset's agency tracker
	agencyCheckedOut := make(map[string]time.Time)
	for k, t := range allCheckedOutAssets {
		agencyCheckedOut[k] = t
	}
	asset.CheckedOut[agency.Name] = agencyCheckedOut

	return retCheckedOutLicenses, nil
}

func (b *BlossomSmartContract) Checkin(stub shim.ChaincodeStubInterface, assetID string, licenses []string, agencyName string) error {
	var (
		asset  = &model.Asset{}
		agency = &model.Agency{}
		err    error
	)

	// get agency
	if agency, err = b.Agency(stub, agencyName); err != nil {
		return errors.Wrapf(err, "error getting agency %q", agencyName)
	}

	// get asset
	if asset, err = b.AssetInfo(stub, assetID); err != nil {
		return errors.Wrapf(err, "error getting info for asset %q", assetID)
	}

	// check in asset logic
	if err = checkin(agency, asset, licenses); err != nil {
		return err
	}

	// update agency
	var bytes []byte
	if bytes, err = json.Marshal(agency); err != nil {
		return errors.Wrapf(err, "error marshaling agency %q", agency.Name)
	}

	if err = stub.PutState(model.AgencyKey(agency.Name), bytes); err != nil {
		return errors.Wrapf(err, "error updating agency state")
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
	// the agency users will no longer be able to see the licenses
	if err := pdp.NewAssetDecider().Checkin(stub, agencyName, assetID, licenses); err != nil {
		return errors.Wrapf(err, "error checking in licenses in NGAC")
	}
	// end NGAC

	return nil
}

func checkin(agency *model.Agency, asset *model.Asset, licenses []string) error {
	checkedOut := agency.Assets[asset.ID]
	for _, returnedKey := range licenses {
		// check that the returned license is leased to the agency
		if _, ok := agency.Assets[asset.ID][returnedKey]; !ok {
			return errors.Errorf("returned key %s was not checked out by %s", returnedKey, agency.Name)
		}

		delete(checkedOut, returnedKey)
	}

	// if all licenses were returned remove asset from agency's checked out
	if len(checkedOut) == 0 {
		delete(agency.Assets, asset.ID)
	} else {
		// update agency licenses
		agency.Assets[asset.ID] = checkedOut
	}

	agencyCheckedOut, ok := asset.CheckedOut[agency.Name]
	if !ok {
		return errors.Errorf("agency %s has not checked out any licenses for asset %s", agency.Name, asset.ID)
	}

	for _, license := range licenses {
		// check that the agency has the license checked out
		if _, ok = agencyCheckedOut[license]; !ok {
			return errors.Errorf("returned license %s was not checked out by %s", license, agency.Name)
		}

		// remove the returned license from the checked out licenses
		delete(agencyCheckedOut, license)

		// add the returned license to the available licenses
		asset.AvailableLicenses = append(asset.AvailableLicenses, license)
	}

	// if all licenses are returned, remove the agency from the asset
	if len(agencyCheckedOut) == 0 {
		delete(asset.CheckedOut, agency.Name)
	}

	// update number of available licenses
	asset.Available += len(licenses)

	return nil
}
