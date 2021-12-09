package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
	"strings"
)

type (
	// SwIDInterface provides the functions to interact with SwID tags in fabric.
	SwIDInterface interface {
		// ReportSwID is used by Accounts to report to Blossom when a software user has installed a piece of software associated
		// with an asset that account has checked out. This function will invoke NGAC chaincode to add the SwID to the NGAC graph.
		ReportSwID(stub shim.ChaincodeStubInterface, account, primaryTag, asset, license, xml string) error

		// GetSwID returns the SwID object including the XML that matches the provided primaryTag parameter.
		GetSwID(stub shim.ChaincodeStubInterface, account, primaryTag string) (*model.SwID, error)

		// GetSwIDsAssociatedWithAsset returns the SwIDs that are associated with the given asset.
		GetSwIDsAssociatedWithAsset(stub shim.ChaincodeStubInterface, account, assetID string) ([]*model.SwID, error)
	}
)

func NewSwIDContract() SwIDInterface {
	return &BlossomSmartContract{}
}

func (b *BlossomSmartContract) swidExists(stub shim.ChaincodeStubInterface, account, primaryTag string) (bool, error) {
	data, err := stub.GetPrivateData(AccountCollectionName(account), model.SwIDKey(primaryTag))
	if err != nil {
		return false, errors.Wrapf(err, "error checking if SwID with primary tag %q already exists on the ledger", primaryTag)
	}

	return data != nil, nil
}

func (b *BlossomSmartContract) ReportSwID(stub shim.ChaincodeStubInterface, account, primaryTag, asset, license, xml string) error {
	if ok, err := b.swidExists(stub, account, primaryTag); err != nil {
		return errors.Wrapf(err, "error checking if SwID with primary tag %s already exists", primaryTag)
	} else if ok {
		return errors.Errorf("a SwID tag with the primary tag %s has already been reported", primaryTag)
	}

	collection := AccountCollectionName(account)

	// ngac check
	if err := pdp.CanReportSwID(stub, collection, account); err != nil {
		return errors.Wrapf(err, "ngac check failed")
	}

	swid := &model.SwID{
		PrimaryTag: primaryTag,
		XML:        xml,
		Asset:      asset,
		License:    license,
	}

	swidBytes, err := json.Marshal(swid)
	if err != nil {
		return errors.Wrapf(err, "error serializing swid tag")
	}

	if err = stub.PutPrivateData(collection, model.SwIDKey(swid.PrimaryTag), swidBytes); err != nil {
		return errors.Wrapf(err, "error updating SwID %s", swid.PrimaryTag)
	}

	return nil
}

func (b *BlossomSmartContract) GetSwID(stub shim.ChaincodeStubInterface, account, primaryTag string) (*model.SwID, error) {
	if ok, err := b.swidExists(stub, account, primaryTag); err != nil {
		return nil, errors.Wrapf(err, "error checking if SwID with primary tag %s already exists", primaryTag)
	} else if !ok {
		return nil, errors.Errorf("a SwID tag with the primary tag %s has not been reported", primaryTag)
	}

	var (
		swidBytes []byte
		err       error
	)

	if swidBytes, err = stub.GetPrivateData(AccountCollectionName(account), model.SwIDKey(primaryTag)); err != nil {
		return nil, errors.Wrapf(err, "error getting SwID %s", primaryTag)
	}

	swid := &model.SwID{}
	if err = json.Unmarshal(swidBytes, swid); err != nil {
		return nil, errors.Wrapf(err, "error deserializing SwID tag %s", primaryTag)
	}

	return swid, nil
}

func (b *BlossomSmartContract) GetSwIDsAssociatedWithAsset(stub shim.ChaincodeStubInterface, account, asset string) ([]*model.SwID, error) {
	resultsIterator, err := stub.GetPrivateDataByRange(AccountCollectionName(account), "", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	swids := make([]*model.SwID, 0)
	for resultsIterator.HasNext() {
		var queryResponse *queryresult.KV
		if queryResponse, err = resultsIterator.Next(); err != nil {
			return nil, err
		}

		// assets on the ledger begin with the asset prefix -- ignore other results
		if !strings.HasPrefix(queryResponse.Key, model.SwIDPrefix) {
			continue
		}

		swid := &model.SwID{}
		if err = json.Unmarshal(queryResponse.Value, swid); err != nil {
			return nil, err
		}

		// continue if the asset associated with this swid tag does not match the given asset ID
		if swid.Asset != asset {
			continue
		}

		swids = append(swids, swid)
	}

	return swids, nil
}
