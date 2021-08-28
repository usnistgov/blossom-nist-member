package rbac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac"
	accountpap "github.com/usnistgov/blossom/chaincode/ngac/pap/account"
)

type AccountPolicy struct {
	graph pip.Graph
}

func NewAccountPolicy(graph pip.Graph) AccountPolicy {
	return AccountPolicy{graph: graph}
}

func (a AccountPolicy) RequestAccount(account *model.Account) error {
	// assign the account object to the accounts attribute
	if err := a.graph.Assign(accountpap.InfoObjectName(account.Name), AccountsOA); err != nil {
		return errors.Wrapf(err, "error assigning account %q to accounts attribute", account.Name)
	}

	// assign the users to their attributes
	if err := a.graph.Assign(ngac.FormatUsername(account.Users.SystemOwner, account.MSPID), SystemOwnerUA); err != nil {
		return errors.Wrapf(err, "error assigning system owner to SystemOwner user attribute")
	}

	if err := a.graph.Assign(ngac.FormatUsername(account.Users.AcquisitionSpecialist, account.MSPID), AcquisitionSpecialistUA); err != nil {
		return errors.Wrapf(err, "error assigning acquisition specialist to AcquisitionSpecialist user attribute")
	}

	if err := a.graph.Assign(ngac.FormatUsername(account.Users.SystemAdministrator, account.MSPID), SystemAdministratorUA); err != nil {
		return errors.Wrapf(err, "error assigning system owner to SystemOwner user attribute")
	}

	return nil
}
