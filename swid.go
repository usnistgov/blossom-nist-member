package main

import (
	"encoding/json"
	"fmt"
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
		ReportSwID(stub shim.ChaincodeStubInterface) error

		// GetSwID returns the SwID object including the XML that matches the provided primaryTag parameter.
		GetSwID(stub shim.ChaincodeStubInterface) (*model.SwID, error)

		// GetSwIDsAssociatedWithAsset returns the SwIDs that are associated with the given asset.
		GetSwIDsAssociatedWithAsset(stub shim.ChaincodeStubInterface) ([]*model.SwID, error)
	}
)

func NewSwIDContract() SwIDInterface {
	return &BlossomSmartContract{}
}

func (b *BlossomSmartContract) swidExists(stub shim.ChaincodeStubInterface, account, primaryTag string) (bool, error) {
	data, err := stub.GetPrivateData(AccountCollection(account), model.SwIDKey(primaryTag))
	if err != nil {
		return false, errors.Wrapf(err, "error checking if SwID with primary tag %q already exists on the ledger", primaryTag)
	}

	return data != nil, nil
}

func (b *BlossomSmartContract) ReportSwID(stub shim.ChaincodeStubInterface) error {
	transientInput, err := getReportSwIDTransientInput(stub)
	if err != nil {
		return fmt.Errorf("error getting transient input: %v", err)
	}

	if ok, err := b.swidExists(stub, transientInput.Account, transientInput.PrimaryTag); err != nil {
		return errors.Wrapf(err, "error checking if SwID with primary tag %s already exists", transientInput.PrimaryTag)
	} else if ok {
		return errors.Errorf("a SwID tag with the primary tag %s has already been reported", transientInput.PrimaryTag)
	}

	collection := AccountCollection(transientInput.Account)

	// ngac check
	if err = pdp.CanReportSwID(stub, collection, transientInput.Account); err != nil {
		return errors.Wrapf(err, "ngac check failed")
	}

	swid := &model.SwID{
		PrimaryTag: transientInput.PrimaryTag,
		XML:        transientInput.Xml,
		Asset:      transientInput.Asset,
		License:    transientInput.License,
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

func (b *BlossomSmartContract) GetSwID(stub shim.ChaincodeStubInterface) (*model.SwID, error) {
	transientInput, err := getGetSwIDTransientInput(stub)
	if err != nil {
		return nil, fmt.Errorf("error getting transient input: %v", err)
	}

	if ok, err := b.swidExists(stub, transientInput.Account, transientInput.PrimaryTag); err != nil {
		return nil, errors.Wrapf(err, "error checking if SwID with primary tag %s already exists", transientInput.PrimaryTag)
	} else if !ok {
		return nil, errors.Errorf("a SwID tag with the primary tag %s has not been reported", transientInput.PrimaryTag)
	}
	var swidBytes []byte

	if swidBytes, err = stub.GetPrivateData(AccountCollection(transientInput.Account), model.SwIDKey(transientInput.PrimaryTag)); err != nil {
		return nil, errors.Wrapf(err, "error getting SwID %s", transientInput.PrimaryTag)
	}

	swid := &model.SwID{}
	if err = json.Unmarshal(swidBytes, swid); err != nil {
		return nil, errors.Wrapf(err, "error deserializing SwID tag %s", transientInput.PrimaryTag)
	}

	return swid, nil
}

func (b *BlossomSmartContract) GetSwIDsAssociatedWithAsset(stub shim.ChaincodeStubInterface) ([]*model.SwID, error) {
	transientInput, err := getGetSwIDsAssociatedWithAssetTransientInput(stub)
	if err != nil {
		return nil, fmt.Errorf("error getting transient input: %v", err)
	}

	resultsIterator, err := stub.GetPrivateDataByRange(AccountCollection(transientInput.Account), "", "")

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
		if swid.Asset != transientInput.AssetID {
			continue
		}

		swids = append(swids, swid)
	}

	return swids, nil
}
