package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
	"strings"
)

type (
	// SwIDInterface provides the functions to interact with SwID tags in fabric.
	SwIDInterface interface {
		// ReportSwID is used by Accounts to report to Blossom when a software user has installed a piece of software associated
		// with an asset that account has checked out. The account is extracted from the requesting identity.  The account
		// must have checked out the defined license or this function will fail.
		// TRANSIENT MAP: export ATO=$(echo -n "{\"primary_tag\":\"123\",\"asset\":\"101\",\"license\":\"asset1-license-1\",\"xml\":\"<swid></swid>\"}" | base64 | tr -d \\n)
		ReportSwID(stub shim.ChaincodeStubInterface) error

		// DeleteSwID deletes a swid from the ledger. This would happen in the case of an organziation returning licenses,
		// and the swid no longer being valid.  The requesting user will need to have the correct permissions in NGAC
		// to do so.  The user with pemrission is the system_owner as defined in the account info.
		// TRANSIENT MAP: export swid=$(echo -n "{\"primary_tag\":\"\",\"account\":\"\"}" | base64 | tr -d \\n)
		DeleteSwID(stub shim.ChaincodeStubInterface) error

		// GetSwID returns the SwID object including the XML that matches the provided primaryTag parameter.
		// TRANSIENT MAP: export swid=$(echo -n "{\"primary_tag\":\"\",\"account\":\"\"}" | base64 | tr -d \\n)
		GetSwID(stub shim.ChaincodeStubInterface) (*model.SwID, error)

		// GetSwIDsAssociatedWithAsset returns the SwIDs that are associated with the given asset for an account.
		GetSwIDsAssociatedWithAsset(stub shim.ChaincodeStubInterface, account string, assetID string) ([]*model.SwID, error)
	}
)

func NewSwIDContract() SwIDInterface {
	return &BlossomSmartContract{}
}

func (b *BlossomSmartContract) swidExists(stub shim.ChaincodeStubInterface, account, primaryTag string) (bool, error) {
	data, err := stub.GetPrivateData(AccountCollection(account), model.SwIDKey(primaryTag))
	if err != nil {
		return false, fmt.Errorf("error checking if SwID with primary tag %q already exists on the ledger: %v", primaryTag, err)
	}

	return data != nil, nil
}

func (b *BlossomSmartContract) ReportSwID(stub shim.ChaincodeStubInterface) error {
	transientInput, err := getReportSwIDTransientInput(stub)
	if err != nil {
		return fmt.Errorf("error getting transient input: %v", err)
	}

	account, err := accountName(stub)
	if err != nil {
		return fmt.Errorf("error getting account name from stub: %v", err)
	}

	if ok, err := b.swidExists(stub, account, transientInput.PrimaryTag); err != nil {
		return fmt.Errorf("error checking if SwID with primary tag %s already exists: %v", transientInput.PrimaryTag, err)
	} else if ok {
		return fmt.Errorf("a SwID tag with the primary tag %s has already been reported", transientInput.PrimaryTag)
	}

	collection := AccountCollection(account)

	// check if this account did indeed checkout the license in the request
	licenses, err := b.Licenses(stub, account, transientInput.Asset)
	if err != nil {
		return err
	}

	ok := false
	for license := range licenses {
		if license == transientInput.License {
			ok = true
		}
	}

	if !ok {
		return fmt.Errorf("account %s cannot report a swid using license %s", account, transientInput.License)
	}

	// ngac check
	if err = pdp.CanReportSwID(stub, collection, account); err != nil {
		return fmt.Errorf("ngac check failed: %v", err)
	}

	swid := &model.SwID{
		PrimaryTag: transientInput.PrimaryTag,
		XML:        transientInput.Xml,
		Asset:      transientInput.Asset,
		License:    transientInput.License,
	}

	swidBytes, err := json.Marshal(swid)
	if err != nil {
		return fmt.Errorf("error serializing swid tag: %v", err)
	}

	if err = stub.PutPrivateData(collection, model.SwIDKey(swid.PrimaryTag), swidBytes); err != nil {
		return fmt.Errorf("error updating SwID %s: %v", swid.PrimaryTag, err)
	}

	return nil
}

func (b *BlossomSmartContract) DeleteSwID(stub shim.ChaincodeStubInterface) error {
	transientInput, err := getGetSwIDTransientInput(stub)
	if err != nil {
		return fmt.Errorf("error getting transient input: %v", err)
	}

	if ok, err := b.swidExists(stub, transientInput.Account, transientInput.PrimaryTag); err != nil {
		return fmt.Errorf("error checking if SwID with primary tag %s already exists: %v", transientInput.PrimaryTag, err)
	} else if !ok {
		return fmt.Errorf("a SwID tag with the primary tag %s has not been reported", transientInput.PrimaryTag)
	}

	// ngac check
	if err = pdp.CanReportSwID(stub, AccountCollection(transientInput.Account), transientInput.Account); err != nil {
		return fmt.Errorf("ngac check failed: %v", err)
	}

	if err = stub.DelPrivateData(AccountCollection(transientInput.Account), model.SwIDKey(transientInput.PrimaryTag)); err != nil {
		return fmt.Errorf("error getting SwID %s: %v", transientInput.PrimaryTag, err)
	}

	return nil
}

func (b *BlossomSmartContract) GetSwID(stub shim.ChaincodeStubInterface) (*model.SwID, error) {
	transientInput, err := getGetSwIDTransientInput(stub)
	if err != nil {
		return nil, fmt.Errorf("error getting transient input: %v", err)
	}

	if ok, err := b.swidExists(stub, transientInput.Account, transientInput.PrimaryTag); err != nil {
		return nil, fmt.Errorf("error checking if SwID with primary tag %s already exists: %v", transientInput.PrimaryTag, err)
	} else if !ok {
		return nil, fmt.Errorf("a SwID tag with the primary tag %s has not been reported: %v", transientInput.PrimaryTag, err)
	}
	var swidBytes []byte

	if swidBytes, err = stub.GetPrivateData(AccountCollection(transientInput.Account), model.SwIDKey(transientInput.PrimaryTag)); err != nil {
		return nil, fmt.Errorf("error getting SwID %s: %v", transientInput.PrimaryTag, err)
	}

	swid := &model.SwID{}
	if err = json.Unmarshal(swidBytes, swid); err != nil {
		return nil, fmt.Errorf("error deserializing SwID tag %s: %v", transientInput.PrimaryTag, err)
	}

	return swid, nil
}

func (b *BlossomSmartContract) GetSwIDsAssociatedWithAsset(stub shim.ChaincodeStubInterface, account string, assetID string) ([]*model.SwID, error) {
	resultsIterator, err := stub.GetPrivateDataByRange(AccountCollection(account), "", "")

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
		if swid.Asset != assetID {
			continue
		}

		swids = append(swids, swid)
	}

	return swids, nil
}
