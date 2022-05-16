package api

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/usnistgov/blossom/chaincode/collections"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
	"strings"
)

func NewSwIDContract() SwIDInterface {
	return &BlossomSmartContract{}
}

func (b *BlossomSmartContract) swidExists(ctx contractapi.TransactionContextInterface, account, primaryTag string) (bool, error) {
	data, err := ctx.GetStub().GetPrivateData(collections.Account(account), model.SwIDKey(primaryTag))
	if err != nil {
		return false, fmt.Errorf("error checking if SwID with primary tag %q already exists on the ledger: %w", primaryTag, err)
	}

	return data != nil, nil
}

func (b *BlossomSmartContract) ReportSwID(ctx contractapi.TransactionContextInterface) error {
	transientInput, err := getReportSwIDTransientInput(ctx)
	if err != nil {
		return fmt.Errorf("error getting transient input: %w", err)
	}

	account, err := accountName(ctx)
	if err != nil {
		return fmt.Errorf("error getting account name from stub: %w", err)
	}

	if ok, err := b.swidExists(ctx, account, transientInput.PrimaryTag); err != nil {
		return fmt.Errorf("error checking if SwID with primary tag %s already exists: %w", transientInput.PrimaryTag, err)
	} else if ok {
		return fmt.Errorf("a SwID tag with the primary tag %s has already been reported", transientInput.PrimaryTag)
	}

	collection := collections.Account(account)

	// check if this account did indeed checkout the license in the request
	licenses, err := b.GetLicenses(ctx, account, transientInput.Asset)
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
	if err = pdp.CanReportSwID(ctx, account); err != nil {
		return fmt.Errorf("ngac check failed: %w", err)
	}

	swid := &model.SwID{
		PrimaryTag: transientInput.PrimaryTag,
		XML:        transientInput.Xml,
		Asset:      transientInput.Asset,
		License:    transientInput.License,
	}

	swidBytes, err := json.Marshal(swid)
	if err != nil {
		return fmt.Errorf("error serializing swid tag: %w", err)
	}

	if err = ctx.GetStub().PutPrivateData(collection, model.SwIDKey(swid.PrimaryTag), swidBytes); err != nil {
		return fmt.Errorf("error updating SwID %s: %w", swid.PrimaryTag, err)
	}

	return nil
}

func (b *BlossomSmartContract) DeleteSwID(ctx contractapi.TransactionContextInterface) error {
	transientInput, err := getGetSwIDTransientInput(ctx)
	if err != nil {
		return fmt.Errorf("error getting transient input: %w", err)
	}

	if ok, err := b.swidExists(ctx, transientInput.Account, transientInput.PrimaryTag); err != nil {
		return fmt.Errorf("error checking if SwID with primary tag %s already exists: %w", transientInput.PrimaryTag, err)
	} else if !ok {
		return fmt.Errorf("a SwID tag with the primary tag %s has not been reported", transientInput.PrimaryTag)
	}

	// ngac check
	if err = pdp.CanReportSwID(ctx, transientInput.Account); err != nil {
		return fmt.Errorf("ngac check failed: %w", err)
	}

	if err = ctx.GetStub().DelPrivateData(collections.Account(transientInput.Account), model.SwIDKey(transientInput.PrimaryTag)); err != nil {
		return fmt.Errorf("error getting SwID %s: %w", transientInput.PrimaryTag, err)
	}

	return nil
}

func (b *BlossomSmartContract) GetSwID(ctx contractapi.TransactionContextInterface) (*model.SwID, error) {
	transientInput, err := getGetSwIDTransientInput(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting transient input: %w", err)
	}

	if ok, err := b.swidExists(ctx, transientInput.Account, transientInput.PrimaryTag); err != nil {
		return nil, fmt.Errorf("error checking if SwID with primary tag %s already exists: %w", transientInput.PrimaryTag, err)
	} else if !ok {
		return nil, fmt.Errorf("a SwID tag with the primary tag %s has not been reported: %w", transientInput.PrimaryTag, err)
	}
	var swidBytes []byte

	if swidBytes, err = ctx.GetStub().GetPrivateData(collections.Account(transientInput.Account), model.SwIDKey(transientInput.PrimaryTag)); err != nil {
		return nil, fmt.Errorf("error getting SwID %s: %w", transientInput.PrimaryTag, err)
	}

	swid := &model.SwID{}
	if err = json.Unmarshal(swidBytes, swid); err != nil {
		return nil, fmt.Errorf("error deserializing SwID tag %s: %w", transientInput.PrimaryTag, err)
	}

	return swid, nil
}

func (b *BlossomSmartContract) GetSwIDsAssociatedWithAsset(ctx contractapi.TransactionContextInterface, account string, assetID string) ([]*model.SwID, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collections.Account(account), "", "")

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
