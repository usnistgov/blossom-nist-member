package pdp

import (
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
	accountpap "github.com/usnistgov/blossom/chaincode/ngac/pap/account"
)

// AccountDecider is the Policy Decision Point (PDP) for the Account API
type AccountDecider struct {
	// user is the user that is currently executing a function
	user string
	// pap is the policy administration point for accounts
	pap *pap.AccountAdmin
	// decider is the NGAC decider used to make decisions
	decider pdp.Decider
}

// NewAccountDecider creates a new AccountDecider with the user from the stub and a NGAC Decider using the NGAC graph
// from the ledger.
func NewAccountDecider() *AccountDecider {
	return &AccountDecider{}
}

func (a *AccountDecider) setup(stub shim.ChaincodeStubInterface) error {
	user, err := GetUser(stub)
	if err != nil {
		return errors.Wrapf(err, "error getting user from request")
	}

	a.user = user

	// initialize the account policy administration point
	a.pap, err = pap.NewAccountAdmin(stub)
	if err != nil {
		return errors.Wrapf(err, "error initializing account administraion point")
	}

	a.decider = pdp.NewDecider(a.pap.Graph())

	return nil
}

func (a *AccountDecider) FilterAccounts(stub shim.ChaincodeStubInterface, accounts []*model.Account) ([]*model.Account, error) {
	if err := a.setup(stub); err != nil {
		return nil, errors.Wrapf(err, "error setting up account decider")
	}

	filteredAccounts := make([]*model.Account, 0)
	for _, account := range accounts {
		if err := a.filterAccount(account); err != nil {
			return nil, errors.Wrapf(err, "error filtering account")
		}

		if account.Name == "" {
			continue
		}

		filteredAccounts = append(filteredAccounts, account)
	}

	return filteredAccounts, nil
}

func (a *AccountDecider) FilterAccount(stub shim.ChaincodeStubInterface, account *model.Account) error {
	if err := a.setup(stub); err != nil {
		return errors.Wrapf(err, "error setting up account decider")
	}

	return a.filterAccount(account)
}

func (a *AccountDecider) filterAccount(account *model.Account) error {
	permissions, err := a.decider.ListPermissions(a.user, accountpap.InfoObjectName(account.Name))
	if err != nil {
		return errors.Wrapf(err, "error getting permissions for user %s on account %s", a.user, account.Name)
	}

	// if the user cannot view account on the account info object, return an empty account
	if !permissions.Contains(operations.ViewAccount) {
		account.Assets = make(map[string]map[string]model.DateTime)
		account.Status = ""
		account.ATO = ""
		account.Users = model.Users{}
		account.MSPID = ""
		return nil
	}

	if !permissions.Contains(operations.ViewATO) {
		account.ATO = ""
	}

	if !permissions.Contains(operations.ViewMSPID) {
		account.MSPID = ""
	}

	if !permissions.Contains(operations.ViewUsers) {
		account.Users = model.Users{}
	}

	if !permissions.Contains(operations.ViewStatus) {
		account.Status = ""
	}

	if !permissions.Contains(operations.ViewAccountLicenses) {
		account.Assets = make(map[string]map[string]model.DateTime)
	}

	return nil
}

func (a *AccountDecider) RequestAccount(stub shim.ChaincodeStubInterface, account *model.Account) error {
	if err := a.setup(stub); err != nil {
		return errors.Wrapf(err, "error setting up account decider")
	}

	// any user can create an account
	return a.pap.RequestAccount(stub, account)
}

func (a *AccountDecider) UploadATO(stub shim.ChaincodeStubInterface, account string) error {
	if err := a.setup(stub); err != nil {
		return errors.Wrapf(err, "error setting up account decider")
	}

	if ok, err := a.decider.HasPermissions(a.user, accountpap.InfoObjectName(account), operations.UploadATO); err != nil {
		return errors.Wrapf(err, "error checking if user %s can upload an ATO for account %s", a.user, account)
	} else if !ok {
		return ErrAccessDenied
	}

	// nothing to update in the account admin
	return nil
}

func (a *AccountDecider) UpdateAccountStatus(stub shim.ChaincodeStubInterface, account string, status model.Status) error {
	if err := a.setup(stub); err != nil {
		return errors.Wrapf(err, "error setting up account decider")
	}

	if ok, err := a.decider.HasPermissions(a.user, accountpap.InfoObjectName(account), operations.UpdateAccountStatus); err != nil {
		return errors.Wrapf(err, "error checking if user %s can update status of account %s", a.user, account)
	} else if !ok {
		return ErrAccessDenied
	}

	return a.pap.UpdateAccountStatus(stub, account, status)
}
