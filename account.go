package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	events "github.com/usnistgov/blossom/chaincode/ngac/epp"
	decider "github.com/usnistgov/blossom/chaincode/ngac/pdp"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/usnistgov/blossom/chaincode/model"
)

type (
	// AccountInterface provides the functions to interact with Accounts in blossom.
	AccountInterface interface {
		// RequestAccount allows accounts to request an account in the Blossom system. The systemOwner, systemAdmin, and
		// acqSpec will be added as users to the NGAC graph, and given the appropriate permissions on the account. The
		// ato can be empty and uploaded via UploadATO later. The name of the acount is the MSPID of the requesting
		// user's member.
		// TRANSIENT MAP: export ACCOUNT=$(echo -n "{\"system_owner\":\"\",\"system_admin\":\"\",\"acquisition_specialist\":\"\",\"ato\":\"\"}" | base64 | tr -d \\n)
		RequestAccount(stub shim.ChaincodeStubInterface) error

		// UploadATO updates the ATO field of the Account with the given name.
		// TRANSIENT MAP: export ATO=$(echo -n "{\"ato\":\"\"}" | base64 | tr -d \\n)
		UploadATO(stub shim.ChaincodeStubInterface) error

		// UpdateAccountStatus updates the status of an account in Blossom. The status is one of:
		//		"PENDING_APPROVAL",
		//		"PENDING_ATO",
		//		"ACTIVE",
		//		"INACTIVE_DENIED",
		//		"INACTIVE_ATO",
		//		"INACTIVE_OPTOUT",
		//		"INACTIVE_SECURITY_RISK",
		//		"INACTIVE_ROB"
		// Updating the status to Active allows the account to read and write to blossom.
		// Updating the status to Pending allows the account to read write only account related information such as ATOs.
		// Updating the status to Inactive provides the same NGAC consequences as Pending
		UpdateAccountStatus(stub shim.ChaincodeStubInterface, account string, status string) error

		// Accounts returns the public info of all accounts that are registered with Blossom.
		Accounts(stub shim.ChaincodeStubInterface) ([]*model.AccountPublic, error)

		// Account returns the account information of the account with the provided name.  Any fields of any account the user
		// does not have access to will not be returned.
		Account(stub shim.ChaincodeStubInterface, account string) (*model.Account, error)

		// GetHistory returns the transaction history of the account.
		GetHistory(stub shim.ChaincodeStubInterface, account string) ([]model.HistorySnapshot, error)
	}
)

func NewAccountContract() AccountInterface {
	return &BlossomSmartContract{}
}

func accountExists(stub shim.ChaincodeStubInterface, accountName string) (bool, error) {
	data, err := stub.GetState(model.AccountKey(accountName))
	if err != nil {
		return false, fmt.Errorf("error checking if account %q already exists on the ledger: %v", accountName, err)
	}

	return data != nil, nil
}

func accountName(stub shim.ChaincodeStubInterface) (string, error) {
	return cid.GetMSPID(stub)
}

func (b *BlossomSmartContract) RequestAccount(stub shim.ChaincodeStubInterface) error {
	transientInput, err := getAccountTransientInput(stub)
	if err != nil {
		return fmt.Errorf("error getting transient input: %v", err)
	}

	accountName, err := accountName(stub)
	if err != nil {
		return fmt.Errorf("error retrieving MSPID from stub: %v", err)
	}

	// check that an account doesn't already exist with the same name
	if ok, err := accountExists(stub, accountName); err != nil {
		return fmt.Errorf("error requesting account: %v", err)
	} else if ok {
		return fmt.Errorf("an account with the name %q already exists", accountName)
	}

	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return fmt.Errorf("error getting mspid: %v", err)
	}

	// account public goes on public ledger
	acctPub := model.AccountPublic{
		Name:   accountName,
		MSPID:  mspid,
		Status: model.PendingATO,
	}

	if transientInput.ATO != "" {
		acctPub.Status = model.PendingApproval
	}

	// account private goes on private data collection for the msp
	acctPvt := model.AccountPrivate{
		ATO: transientInput.ATO,
		Users: model.Users{
			SystemOwner:           transientInput.SystemOwner,
			SystemAdministrator:   transientInput.SystemAdmin,
			AcquisitionSpecialist: transientInput.AcquisitionSpecialist,
		},
		Assets: make(map[string]map[string]model.DateTime),
	}

	// add account public to world state
	pubBytes, err := json.Marshal(acctPub)
	if err != nil {
		return fmt.Errorf("error marshaling private account details for %q: %v", accountName, err)
	}

	if err = stub.PutState(model.AccountKey(accountName), pubBytes); err != nil {
		return fmt.Errorf("error adding account to ledger: %v", err)
	}

	// add account private to pdc
	pvtBytes, err := json.Marshal(acctPvt)
	if err != nil {
		return fmt.Errorf("error marshaling private account details for %q: %v", accountName, err)
	}

	collection := AccountCollectionName(accountName)

	if err = stub.PutPrivateData(collection, model.AccountKey(accountName), pvtBytes); err != nil {
		return fmt.Errorf("error putting private data: %v", err)
	}

	// process the request_account event in NGAC
	if err = events.ProcessRequestAccount(stub, collection, accountName,
		transientInput.SystemOwner, transientInput.SystemAdmin, transientInput.AcquisitionSpecialist); err != nil {
		return fmt.Errorf("error processing request_account event: %v", err)
	}

	return nil
}

func (b *BlossomSmartContract) UploadATO(stub shim.ChaincodeStubInterface) error {
	transientInput, err := getUploadATOTransientInput(stub)
	if err != nil {
		return fmt.Errorf("error getting transient input: %v", err)
	}

	accountName, err := accountName(stub)
	if err != nil {
		return fmt.Errorf("error getting mspid: %v", err)
	}

	if ok, err := accountExists(stub, accountName); err != nil {
		return fmt.Errorf("error checking if account %q exists: %v", accountName, err)
	} else if !ok {
		return fmt.Errorf("an account with the name %q does not exist", accountName)
	}

	collection := AccountCollectionName(accountName)

	// ngac check
	if err := decider.CanUploadATO(stub, collection, accountName); err != nil {
		return fmt.Errorf("error uploading ATO for account %s: %v", accountName, err)
	}

	bytes, err := stub.GetPrivateData(collection, model.AccountKey(accountName))
	if err != nil {
		return fmt.Errorf("error getting account %q from world state: %v", accountName, err)
	}

	acctPvt := &model.AccountPrivate{}
	if err = json.Unmarshal(bytes, acctPvt); err != nil {
		return fmt.Errorf("error unmarshaling account %q: %v", accountName, err)
	}

	// update ATO value
	acctPvt.ATO = transientInput.ATO

	// marshal back to json
	if bytes, err = json.Marshal(acctPvt); err != nil {
		return fmt.Errorf("error marshaling account %q: %v", accountName, err)
	}

	// update pdc
	if err = stub.PutPrivateData(collection, model.AccountKey(accountName), bytes); err != nil {
		return fmt.Errorf("error updating ATO for account %q: %v", accountName, err)
	}

	return nil
}

func (b *BlossomSmartContract) UpdateAccountStatus(stub shim.ChaincodeStubInterface, accountName, statusStr string) error {
	status, err := model.GetStatusUpdate(statusStr)
	if err != nil {
		return err
	}

	if ok, err := accountExists(stub, accountName); err != nil {
		return fmt.Errorf("error checking if account %q exists: %v", accountName, err)
	} else if !ok {
		return fmt.Errorf("an account with the name %q does not exist", accountName)
	}

	// ngac check
	if err = decider.CanUpdateAccountStatus(stub, AccountCollectionName(accountName), accountName); err != nil {
		return fmt.Errorf("error updating account status for account %s: %v", accountName, err)
	}

	bytes, err := stub.GetState(model.AccountKey(accountName))
	if err != nil {
		return fmt.Errorf("error getting account %q from world state: %v", accountName, err)
	}

	acctPub := &model.AccountPublic{}
	if err = json.Unmarshal(bytes, acctPub); err != nil {
		return fmt.Errorf("error unmarshaling account %q: %v", accountName, err)
	}

	// update ATO value
	acctPub.Status = status

	// marshal back to json
	if bytes, err = json.Marshal(acctPub); err != nil {
		return fmt.Errorf("error marshaling account %q: %v", accountName, err)
	}

	// update world state
	if err = stub.PutState(model.AccountKey(accountName), bytes); err != nil {
		return fmt.Errorf("error updating status of account %q: %v", accountName, err)
	}

	// process event
	return handleUpdateAccountStatusEvent(stub, accountName, status)
}

func handleUpdateAccountStatusEvent(stub shim.ChaincodeStubInterface, accountName string, status model.Status) error {
	var f func(shim.ChaincodeStubInterface, string, string) error
	switch status {
	case model.PendingApproval:
		f = events.ProcessSetAccountPending
	case model.PendingATO:
		f = events.ProcessSetAccountPending
	case model.Active:
		f = events.ProcessSetAccountActive
	case model.InactiveDenied:
		f = events.ProcessSetAccountInactive
	case model.InactiveATO:
		f = events.ProcessSetAccountInactive
	case model.InactiveOptOut:
		f = events.ProcessSetAccountInactive
	case model.InactiveSecurityRisk:
		f = events.ProcessSetAccountInactive
	case model.InactiveROB:
		f = events.ProcessSetAccountInactive
	default:
		return fmt.Errorf("unknown status: %s", status)
	}

	return f(stub, AccountCollectionName(accountName), accountName)
}

func (b *BlossomSmartContract) Accounts(stub shim.ChaincodeStubInterface) ([]*model.AccountPublic, error) {
	resultsIterator, err := stub.GetStateByRange(string([]byte{model.AccountPrefix}), string([]byte{model.AccountPrefix + 1}))
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

		acctPub := &model.AccountPublic{}
		if err = json.Unmarshal(queryResponse.Value, acctPub); err != nil {
			return nil, err
		}

		accounts = append(accounts, acctPub)
	}

	return accounts, nil
}

func (b *BlossomSmartContract) Account(stub shim.ChaincodeStubInterface, accountName string) (*model.Account, error) {
	var (
		acctPub = &model.AccountPublic{}
		acctPvt = &model.AccountPrivate{}
		bytes   []byte
		err     error
	)

	if bytes, err = stub.GetState(model.AccountKey(accountName)); err != nil {
		return nil, fmt.Errorf("error getting account public info from ledger: %v", err)
	}

	if err = json.Unmarshal(bytes, acctPub); err != nil {
		return nil, fmt.Errorf("error deserializing account public info: %v", err)
	}

	if bytes, err = stub.GetPrivateData(AccountCollectionName(accountName), model.AccountKey(accountName)); err != nil {
		// ignore error if a user does not have access to the private data collection of the account
		// they can still have access to the public info
		fmt.Printf("error occurred reading pvtdata: %v\n", err)
	} else {
		if err = json.Unmarshal(bytes, acctPvt); err != nil {
			return nil, fmt.Errorf("error deserializing account private info: %v", err)
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

func (b *BlossomSmartContract) GetHistory(stub shim.ChaincodeStubInterface, account string) ([]model.HistorySnapshot, error) {
	history := []model.HistorySnapshot{}

	iter, err := stub.GetHistoryForKey(model.AccountKey(account))
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
