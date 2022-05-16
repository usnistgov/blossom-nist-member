package api

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/usnistgov/blossom/chaincode/collections"
	events "github.com/usnistgov/blossom/chaincode/ngac/epp"
	decider "github.com/usnistgov/blossom/chaincode/ngac/pdp"
	"strings"
	"time"

	"github.com/usnistgov/blossom/chaincode/model"
)

func NewAccountContract() AccountInterface {
	return &BlossomSmartContract{}
}

func accountExists(ctx contractapi.TransactionContextInterface, accountName string) (bool, error) {
	data, err := ctx.GetStub().GetState(model.AccountKey(accountName))
	if err != nil {
		return false, fmt.Errorf("error checking if account %q already exists on the ledger: %w", accountName, err)
	}

	return data != nil, nil
}

func accountName(ctx contractapi.TransactionContextInterface) (string, error) {
	return ctx.GetClientIdentity().GetMSPID()
}

func (b *BlossomSmartContract) RequestAccount(ctx contractapi.TransactionContextInterface) error {
	/*
		removing for now until it's clear users can have the admin attribute
		attr, _, err := cid.GetAttributeValue(stub, "hf.Type")
		if err != nil {
			return err
		}

		// check if requesting user is an admin
		if attr != "admin" {
			return fmt.Errorf("only org admins can request accounts")
		}*/

	transientInput, err := getAccountTransientInput(ctx)
	if err != nil {
		return fmt.Errorf("error getting transient input: %w", err)
	}

	accountName, err := accountName(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving MSPID from ctx: %w", err)
	}

	// check that an account doesn't already exist with the same name
	if ok, err := accountExists(ctx, accountName); err != nil {
		return fmt.Errorf("error requesting account: %w", err)
	} else if ok {
		return fmt.Errorf("an account with the name %q already exists", accountName)
	}

	mspid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("error getting mspid: %w", err)
	}

	// account public goes on public ledger
	acctPub := model.AccountPublic{
		Name:   accountName,
		MSPID:  mspid,
		Status: model.PendingApproval,
	}

	// account private goes on private data collection for the msp
	acctPvt := model.AccountPrivate{
		ATO: "",
		Users: model.Users{
			SystemOwner:           transientInput.SystemOwner,
			SystemAdministrator:   transientInput.SystemAdmin,
			AcquisitionSpecialist: transientInput.AcquisitionSpecialist,
		},
		Assets: make(map[string]map[string]string),
	}

	// add account public to world state
	pubBytes, err := json.Marshal(acctPub)
	if err != nil {
		return fmt.Errorf("error marshaling private account details for %q: %w", accountName, err)
	}

	if err = ctx.GetStub().PutState(model.AccountKey(accountName), pubBytes); err != nil {
		return fmt.Errorf("error adding account to ledger: %w", err)
	}

	// add account private to pdc
	pvtBytes, err := json.Marshal(acctPvt)
	if err != nil {
		return fmt.Errorf("error marshaling private account details for %q: %w", accountName, err)
	}

	collection := collections.Account(accountName)

	if err = ctx.GetStub().PutPrivateData(collection, model.AccountKey(accountName), pvtBytes); err != nil {
		return fmt.Errorf("error putting private data: %w", err)
	}

	return nil
}

func (b *BlossomSmartContract) ApproveAccount(ctx contractapi.TransactionContextInterface, account string) error {
	var (
		acctPvt model.AccountPrivate
		bytes   []byte
		err     error
	)

	// get account private details from PDC to add users to NGAC graph
	if bytes, err = ctx.GetStub().GetPrivateData(collections.Account(account), model.AccountKey(account)); err != nil {
		return fmt.Errorf("error getting private data: %w", err)
	} else {
		if err = json.Unmarshal(bytes, &acctPvt); err != nil {
			return fmt.Errorf("error deserializing account private info: %w", err)
		}
	}

	if err = decider.CanApproveAccount(ctx); err != nil {
		return fmt.Errorf("error approving account in NGAC: %w", err)
	}

	// update account status
	status := model.PendingATO
	bytes, err = ctx.GetStub().GetState(model.AccountKey(account))
	if err != nil {
		return fmt.Errorf("error getting account %q from world state: %w", account, err)
	}

	acctPub := &model.AccountPublic{}
	if err = json.Unmarshal(bytes, acctPub); err != nil {
		return fmt.Errorf("error unmarshaling account %q: %w", account, err)
	}

	// update status
	acctPub.Status = status

	// marshal back to json
	if bytes, err = json.Marshal(acctPub); err != nil {
		return fmt.Errorf("error marshaling account %q: %w", account, err)
	}

	// update world state
	if err = ctx.GetStub().PutState(model.AccountKey(account), bytes); err != nil {
		return fmt.Errorf("error updating status of account %q: %w", account, err)
	}

	return events.ProcessApproveAccount(ctx, account)
}

func (b *BlossomSmartContract) UploadATO(ctx contractapi.TransactionContextInterface) error {
	transientInput, err := getUploadATOTransientInput(ctx)
	if err != nil {
		return fmt.Errorf("error getting transient input: %w", err)
	}

	accountName, err := accountName(ctx)
	if err != nil {
		return fmt.Errorf("error getting mspid: %w", err)
	}

	if ok, err := accountExists(ctx, accountName); err != nil {
		return fmt.Errorf("error checking if account %q exists: %w", accountName, err)
	} else if !ok {
		return fmt.Errorf("an account with the name %q does not exist", accountName)
	}

	collection := collections.Account(accountName)

	// ngac check
	if err := decider.CanUploadATO(ctx, accountName); err != nil {
		return fmt.Errorf("error uploading ATO for account %s: %w", accountName, err)
	}

	bytes, err := ctx.GetStub().GetPrivateData(collection, model.AccountKey(accountName))
	if err != nil {
		return fmt.Errorf("error getting account %q from world state: %w", accountName, err)
	}

	acctPvt := &model.AccountPrivate{}
	if err = json.Unmarshal(bytes, acctPvt); err != nil {
		return fmt.Errorf("error unmarshaling account %q: %w", accountName, err)
	}

	// update ATO value
	acctPvt.ATO = transientInput.ATO

	// marshal back to json
	if bytes, err = json.Marshal(acctPvt); err != nil {
		return fmt.Errorf("error marshaling account %q: %w", accountName, err)
	}

	// update pdc
	if err = ctx.GetStub().PutPrivateData(collection, model.AccountKey(accountName), bytes); err != nil {
		return fmt.Errorf("error updating ATO for account %q: %w", accountName, err)
	}

	return nil
}

func (b *BlossomSmartContract) UpdateAccountStatus(ctx contractapi.TransactionContextInterface, accountName, statusStr string) error {
	status, err := model.GetStatusUpdate(statusStr)
	if err != nil {
		return err
	}

	if ok, err := accountExists(ctx, accountName); err != nil {
		return fmt.Errorf("error checking if account %q exists: %w", accountName, err)
	} else if !ok {
		return fmt.Errorf("an account with the name %q does not exist", accountName)
	}

	// ngac check
	if err = decider.CanUpdateAccountStatus(ctx, accountName); err != nil {
		return fmt.Errorf("error updating account status for account %s: %w", accountName, err)
	}

	bytes, err := ctx.GetStub().GetState(model.AccountKey(accountName))
	if err != nil {
		return fmt.Errorf("error getting account %q from world state: %w", accountName, err)
	}

	acctPub := &model.AccountPublic{}
	if err = json.Unmarshal(bytes, acctPub); err != nil {
		return fmt.Errorf("error unmarshaling account %q: %w", accountName, err)
	}

	// update status
	acctPub.Status = status

	// marshal back to json
	if bytes, err = json.Marshal(acctPub); err != nil {
		return fmt.Errorf("error marshaling account %q: %w", accountName, err)
	}

	// update world state
	if err = ctx.GetStub().PutState(model.AccountKey(accountName), bytes); err != nil {
		return fmt.Errorf("error updating status of account %q: %w", accountName, err)
	}

	// process event
	return events.UpdateAccountStatusEvent(ctx, accountName, collections.Catalog(), status)
}

func (b *BlossomSmartContract) GetAccounts(ctx contractapi.TransactionContextInterface) ([]*model.AccountPublic, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	accounts := make([]*model.AccountPublic, 0)
	for resultsIterator.HasNext() {
		var queryResponse *queryresult.KV
		if queryResponse, err = resultsIterator.Next(); err != nil {
			return nil, err
		}

		if !strings.HasPrefix(queryResponse.Key, model.AccountPrefix) {
			continue
		}

		acctPub := &model.AccountPublic{}
		if err = json.Unmarshal(queryResponse.Value, acctPub); err != nil {
			return nil, err
		}

		accounts = append(accounts, acctPub)
	}

	return accounts, nil
}

func (b *BlossomSmartContract) GetAccount(ctx contractapi.TransactionContextInterface, accountName string) (*model.Account, error) {
	var (
		acctPub = &model.AccountPublic{}
		acctPvt = &model.AccountPrivate{}
		bytes   []byte
		err     error
	)

	if ok, err := accountExists(ctx, accountName); err != nil {
		return nil, fmt.Errorf("error checking if account %q exists: %w", accountName, err)
	} else if !ok {
		return nil, fmt.Errorf("an account with the name %q does not exist", accountName)
	}

	if bytes, err = ctx.GetStub().GetState(model.AccountKey(accountName)); err != nil {
		return nil, fmt.Errorf("error getting account public info from ledger: %w", err)
	}

	if err = json.Unmarshal(bytes, acctPub); err != nil {
		return nil, fmt.Errorf("error deserializing account public info: %w", err)
	}

	if bytes, err = ctx.GetStub().GetPrivateData(collections.Account(accountName), model.AccountKey(accountName)); err != nil {
		// ignore error if a user does not have access to the private data collection of the account
		// they can still have access to the public info
		fmt.Printf("error occurred reading pvtdata: %v\n", err)
	} else {
		if err = json.Unmarshal(bytes, acctPvt); err != nil {
			return nil, fmt.Errorf("error deserializing account private info: %w", err)
		}
	}

	return &model.Account{
		Name:   acctPub.Name,
		MSPID:  acctPub.MSPID,
		Status: acctPub.Status,
		ATO:    acctPvt.ATO,
		Users:  acctPvt.Users,
		Assets: acctPvt.Assets,
	}, nil
}

func (b *BlossomSmartContract) GetHistory(ctx contractapi.TransactionContextInterface, account string) ([]model.HistorySnapshot, error) {
	history := []model.HistorySnapshot{}

	iter, err := ctx.GetStub().GetHistoryForKey(model.AccountKey(account))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	// Getting history for private data is not yet supported, so info will be limited to public account info

	for iter.HasNext() {
		result, err := iter.Next()
		if err != nil {
			return nil, err
		}

		snapshot := model.HistorySnapshot{
			TxId:      result.TxId,
			Timestamp: time.Unix(result.Timestamp.Seconds, int64(result.Timestamp.Nanos)),
		}

		err = json.Unmarshal(result.Value, &snapshot.Value)
		if err != nil {
			return nil, err
		}

		history = append(history, snapshot)
	}

	return history, nil
}
