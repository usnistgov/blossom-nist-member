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
	// AccountInterface provides the functions to interact with Accounts in blossom.
	AccountInterface interface {
		// RequestAccount allows accounts to request an account in the Blossom system.  This function will stage the information
		// provided in the Account parameter in a separate structure until the request is accepted or denied.  The account will
		// be identified by the name provided in the request. The MSPID of the account is needed to distinguish users, who may have
		// the same username in a differing MSPs, in the NGAC system.
		RequestAccount(stub shim.ChaincodeStubInterface, account *model.Account) error

		// UploadATO updates the ATO field of the Account with the given name.
		// TODO placeholder function until ATO model is finalized
		UploadATO(stub shim.ChaincodeStubInterface, account string, ato string) error

		// UpdateAccountStatus updates the status of an account in Blossom.
		// Updating the status to Approved allows the account to read and write to blossom.
		// Updating the status to Pending allows the account to read write only account related information such as ATOs.
		// Updating the status to Inactive provides the same NGAC consequences as Pending
		UpdateAccountStatus(stub shim.ChaincodeStubInterface, account string, status model.Status) error

		// Accounts returns a list of all the accounts that are registered with Blossom.  Any account in which the requesting
		// user does not have access to will not be returned.  Likewise, any fields of any account the user does not have access
		// to will not be returned.
		Accounts(stub shim.ChaincodeStubInterface) ([]*model.Account, error)

		// Account returns the account information of the account with the provided name.  Any fields of any account the user
		// does not have access to will not be returned.
		Account(stub shim.ChaincodeStubInterface, account string) (*model.Account, error)
	}
)

func NewAccountContract() AccountInterface {
	return &BlossomSmartContract{}
}

func (b *BlossomSmartContract) accountExists(stub shim.ChaincodeStubInterface, accountName string) (bool, error) {
	data, err := stub.GetState(model.AccountKey(accountName))
	if err != nil {
		return false, errors.Wrapf(err, "error checking if account %q already exists on the ledger", accountName)
	}

	return data != nil, nil
}

func (b *BlossomSmartContract) RequestAccount(stub shim.ChaincodeStubInterface, account *model.Account) error {
	// check that an account doesn't already exist with the same name
	if ok, err := b.accountExists(stub, account.Name); err != nil {
		return errors.Wrapf(err, "error requesting account")
	} else if ok {
		return errors.Errorf("an account with the name %q already exists", account.Name)
	}

	// begin NGAC
	if err := pdp.NewAccountDecider().RequestAccount(stub, account); err != nil {
		return errors.Wrapf(err, "error adding account to NGAC")
	}
	// end NGAC

	// add account to ledger with pending status
	account.Status = model.PendingApproval
	account.Assets = make(map[string]map[string]time.Time)

	// convert account to bytes
	bytes, err := json.Marshal(account)
	if err != nil {
		return errors.Wrapf(err, "error marshaling account %q", account.Name)
	}

	// add account to world state
	if err = stub.PutState(model.AccountKey(account.Name), bytes); err != nil {
		return errors.Wrapf(err, "error adding account to ledger")
	}

	return nil
}

func (b *BlossomSmartContract) UploadATO(stub shim.ChaincodeStubInterface, accountName string, ato string) error {
	if ok, err := b.accountExists(stub, accountName); err != nil {
		return errors.Wrapf(err, "error checking if account %q exists", accountName)
	} else if !ok {
		return errors.Errorf("an account with the name %q does not exist", accountName)
	}

	// begin NGAC
	if err := pdp.NewAccountDecider().UploadATO(stub, accountName); errors.Is(err, pdp.ErrAccessDenied) {
		return err
	} else if err != nil {
		return errors.Wrapf(err, "error checking if user can update ATO")
	}
	// end NGAC

	bytes, err := stub.GetState(model.AccountKey(accountName))
	if err != nil {
		return errors.Wrapf(err, "error getting account %q from world state", accountName)
	}

	ledgerAccount := &model.Account{}
	if err = json.Unmarshal(bytes, ledgerAccount); err != nil {
		return errors.Wrapf(err, "error unmarshaling account %q", accountName)
	}

	// update ATO value
	ledgerAccount.ATO = ato

	// marshal back to json
	if bytes, err = json.Marshal(ledgerAccount); err != nil {
		return errors.Wrapf(err, "error marshaling account %q", accountName)
	}

	// update world state
	if err = stub.PutState(model.AccountKey(accountName), bytes); err != nil {
		return errors.Wrapf(err, "error updating ATO for account %q", accountName)
	}

	return nil
}

func (b *BlossomSmartContract) UpdateAccountStatus(stub shim.ChaincodeStubInterface, accountName string, status model.Status) error {
	if ok, err := b.accountExists(stub, accountName); err != nil {
		return errors.Wrapf(err, "error checking if account %q exists", accountName)
	} else if !ok {
		return errors.Errorf("an account with the name %q does not exist", accountName)
	}

	// begin NGAC
	if err := pdp.NewAccountDecider().UpdateAccountStatus(stub, accountName, status); errors.Is(err, pdp.ErrAccessDenied) {
		return err
	} else if err != nil {
		return errors.Wrapf(err, "error checking if user can update account status")
	}
	// end NGAC

	bytes, err := stub.GetState(model.AccountKey(accountName))
	if err != nil {
		return errors.Wrapf(err, "error getting account %q from world state", accountName)
	}

	ledgerAccount := &model.Account{}
	if err = json.Unmarshal(bytes, ledgerAccount); err != nil {
		return errors.Wrapf(err, "error unmarshaling account %q", accountName)
	}

	// update ATO value
	ledgerAccount.Status = status

	// marshal back to json
	if bytes, err = json.Marshal(ledgerAccount); err != nil {
		return errors.Wrapf(err, "error marshaling account %q", accountName)
	}

	// update world state
	if err = stub.PutState(model.AccountKey(accountName), bytes); err != nil {
		return errors.Wrapf(err, "error updating status of account %q", accountName)
	}

	return nil
}

func (b *BlossomSmartContract) Accounts(stub shim.ChaincodeStubInterface) ([]*model.Account, error) {
	accounts, err := accounts(stub)
	if err != nil {
		return nil, errors.Wrap(err, "error getting accounts")
	}

	// begin NGAC
	if accounts, err = pdp.NewAccountDecider().FilterAccounts(stub, accounts); err != nil {
		return nil, errors.Wrapf(err, "error filtering accounts")
	}
	// end NGAC

	return accounts, nil
}

func accounts(stub shim.ChaincodeStubInterface) ([]*model.Account, error) {
	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	accounts := make([]*model.Account, 0)
	for resultsIterator.HasNext() {
		var queryResponse *queryresult.KV
		if queryResponse, err = resultsIterator.Next(); err != nil {
			return nil, err
		}

		// accounts on the ledger begin with the account prefix -- ignore other assets
		if !strings.HasPrefix(queryResponse.Key, model.AccountPrefix) {
			continue
		}

		account := &model.Account{}
		if err = json.Unmarshal(queryResponse.Value, account); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (b *BlossomSmartContract) Account(stub shim.ChaincodeStubInterface, accountName string) (*model.Account, error) {
	var (
		account = &model.Account{}
		bytes   []byte
		err     error
	)

	if bytes, err = stub.GetState(model.AccountKey(accountName)); err != nil {
		return nil, errors.Wrapf(err, "error getting account from ledger")
	}

	if err = json.Unmarshal(bytes, account); err != nil {
		return nil, errors.Wrapf(err, "error deserializing account")
	}

	// begin NGAC
	// filter account object removing any fields the user does not have access to
	if err = pdp.NewAccountDecider().FilterAccount(stub, account); err != nil {
		return nil, errors.Wrapf(err, "error filtering account %s", accountName)
	}
	// end NGAC

	return account, nil
}
