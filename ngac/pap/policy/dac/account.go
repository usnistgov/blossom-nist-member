package dac

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac"
	accountpap "github.com/usnistgov/blossom/chaincode/ngac/pap/account"
	assetpap "github.com/usnistgov/blossom/chaincode/ngac/pap/asset"
)

type AccountPolicy struct {
	graph pip.Graph
}

func NewAccountPolicy(graph pip.Graph) AccountPolicy {
	return AccountPolicy{graph: graph}
}

func (a AccountPolicy) RequestAccount(account *model.Account) error {
	// TODO using the PAP to avoid permission check for now as obligations are not yet implemented.  Once obligations are
	// implemented the admin can create one to create the system owners and when executed will be executed on behalf of
	// the admin not the system owner requesting an account

	if account.Users.SystemOwner == "" {
		return errors.Errorf("request missing system owner")
	}

	if account.Users.AcquisitionSpecialist == "" {
		return errors.Errorf("request missing acquisition specialist")
	}

	if account.Users.SystemAdministrator == "" {
		return errors.Errorf("request missing system administrator")
	}

	// create the system owner user
	systemOwnerNode, err := a.graph.CreateNode(ngac.FormatUsername(account.Users.SystemOwner, account.MSPID), pip.User, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating system owner user")
	}

	// create the acquisition specialist user
	acqSpecNode, err := a.graph.CreateNode(ngac.FormatUsername(account.Users.AcquisitionSpecialist, account.MSPID), pip.User, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating acquisition specialist user")
	}

	// create the system administrator user
	sysAdminNode, err := a.graph.CreateNode(ngac.FormatUsername(account.Users.SystemAdministrator, account.MSPID), pip.User, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating system administrator user")
	}

	// create a DAC UA for this account
	accountUA, err := a.graph.CreateNode(accountpap.UserAttributeName(account.Name), pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating DAC user attribute")
	}

	// assign system owner to account ua
	if err = a.graph.Assign(systemOwnerNode.Name, accountUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning SystemOwner to account DAC user attribute")
	}

	// assign acquisition specialist to account ua
	if err = a.graph.Assign(acqSpecNode.Name, accountUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning AcquisitionSpecialist to account DAC user attribute")
	}

	// assign system administrator to account ua
	if err = a.graph.Assign(sysAdminNode.Name, accountUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning SystemAdministrator to account DAC user attribute")
	}

	// assign account ua to policy class
	if err = a.graph.Assign(accountUA.Name, UserAttributeName); err != nil {
		return errors.Wrapf(err, "error assigning account user attribute to DAC user attribute")
	}

	// create an object attribute for the account container
	var accountOA pip.Node
	if accountOA, err = a.graph.CreateNode(accountpap.ObjectAttributeName(account.Name), pip.ObjectAttribute, nil); err != nil {
		return fmt.Errorf("error creating account info object attribute in NGAC: %w", err)
	}

	// assign the account oa to the dac oa
	if err = a.graph.Assign(accountOA.Name, ObjectAttributeName); err != nil {
		return fmt.Errorf("error assigning the account object attribute to the DAC object attribute")
	}

	// create an object to represent the account
	var accountInfo pip.Node
	if accountInfo, err = a.graph.CreateNode(accountpap.InfoObjectName(account.Name), pip.Object,
		map[string]string{"account": account.Name, "type": "account"}); err != nil {
		return fmt.Errorf("error creating account info object attribute in NGAC: %w", err)
	}

	// create an object attribute for the account licenses
	var licensesOA pip.Node
	if licensesOA, err = a.graph.CreateNode(assetpap.AssetsObjectAttribute(account.Name), pip.ObjectAttribute, nil); err != nil {
		return errors.Wrapf(err, "error creating licenses object attribute for account %s", account.Name)
	}

	// assign the account oa to the dac oa
	if err = a.graph.Assign(accountInfo.Name, accountOA.Name); err != nil {
		return fmt.Errorf("error assigning the account info object to the account object attribute")
	}

	// assign the licenses OA to the account oa
	if err = a.graph.Assign(licensesOA.Name, accountOA.Name); err != nil {
		return fmt.Errorf("error assigning the account licenses object attribute to the account object attribute")
	}

	// associate the account ua and oa with all ops because RBAC will enforce permissions for each user's role
	if err = a.graph.Associate(accountUA.Name, accountOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return fmt.Errorf("error associating account user attribute and account object attribute")
	}

	return nil
}
