package api

import (
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
)

type (
	// SwIDInterface provides the functions to interact with SwID tags in fabric.
	SwIDInterface interface {
		// ReportSwID is used by Agencies to report to Blossom when a software user has installed a piece of software associated
		// with a license that agency has checked out. This function will invoke NGAc chaincode to add the SwID to the NGAC graph.
		ReportSwID(ctx contractapi.TransactionContextInterface, swid *model.SwID) error

		// GetSwID returns the SwID object including the XML that matches the provided primaryTag parameter.
		GetSwID(ctx contractapi.TransactionContextInterface, primaryTag string) (*model.SwID, error)

		// GetLicenseSwIDs returns the primary tags of the SwIDs that are associated with the given license ID.
		GetLicenseSwIDs(ctx contractapi.TransactionContextInterface) ([]string, error)
	}
)

func NewSwIDContract() SwIDInterface {
	return &BlossomSmartContract{}
}

func (b *BlossomSmartContract) swidExists(ctx contractapi.TransactionContextInterface, primaryTag string) (bool, error) {
	data, err := ctx.GetStub().GetState(model.SwIDKey(primaryTag))
	if err != nil {
		return false, errors.Wrapf(err, "error checking if SwID with primary tag %q already exists on the ledger", primaryTag)
	}

	return data != nil, nil
}

func (b *BlossomSmartContract) ReportSwID(ctx contractapi.TransactionContextInterface, swid *model.SwID) error {
	if ok, err := b.swidExists(ctx, swid.PrimaryTag); err != nil {
		return errors.Wrapf(err, "error checking if SwID with primary tag %s already exists", swid.PrimaryTag)
	} else if ok {
		return errors.Errorf("a SwID tag with the primary tag %s has already been reported", swid.PrimaryTag)
	}

	swidBytes, err := json.Marshal(swid)
	if err != nil {
		return errors.Wrapf(err, "error serializing swid tag")
	}

	if err = ctx.GetStub().PutState(model.SwIDKey(swid.PrimaryTag), swidBytes); err != nil {
		return errors.Wrapf(err, "error updating SwID %s", swid.PrimaryTag)
	}

	return nil
}

func (b *BlossomSmartContract) GetSwID(ctx contractapi.TransactionContextInterface, primaryTag string) (*model.SwID, error) {
	if ok, err := b.swidExists(ctx, primaryTag); err != nil {
		return nil, errors.Wrapf(err, "error checking if SwID with primary tag %s already exists", primaryTag)
	} else if ok {
		return nil, errors.Errorf("a SwID tag with the primary tag %s has already been reported", primaryTag)
	}

	var (
		swidBytes []byte
		err       error
	)

	if swidBytes, err = ctx.GetStub().GetState(model.SwIDKey(primaryTag)); err != nil {
		return nil, errors.Wrapf(err, "error getting SwID %s", primaryTag)
	}

	swid := &model.SwID{}
	if err = json.Unmarshal(swidBytes, swid); err != nil {
		return nil, errors.Wrapf(err, "error deserializing SwID tag %s", primaryTag)
	}

	return &model.SwID{}, nil
}

func (b *BlossomSmartContract) GetLicenseSwIDs(ctx contractapi.TransactionContextInterface) ([]string, error) {
	return nil, nil
}
