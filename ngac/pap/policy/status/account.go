package status

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	accountpap "github.com/usnistgov/blossom/chaincode/ngac/pap/account"
)

type AccountPolicy struct {
	graph pip.Graph
}

func NewAccountPolicy(graph pip.Graph) AccountPolicy {
	return AccountPolicy{graph: graph}
}

func (a AccountPolicy) RequestAccount(account *model.Account) error {
	// assign users to pending
	accountUA := accountpap.UserAttributeName(account.Name)
	if err := a.graph.Assign(accountUA, PendingUA); err != nil {
		return errors.Wrap(err, "error assigning the account user attribute to the pending user attribute")
	}

	// assign account object to accounts oa
	accountObj := accountpap.InfoObjectName(account.Name)
	if err := a.graph.Assign(accountObj, AccountsOA); err != nil {
		return errors.Wrap(err, "error assigning the account object to the accounts object attribute")
	}

	return nil
}

func (a AccountPolicy) UpdateAccountStatus(accountName string, status model.Status) error {
	switch status {
	case model.Approved:
		return a.approved(accountName)
	case model.PendingApproval:
		return a.pending(accountName)
	case model.PendingATO:
		return a.pending(accountName)
	case model.PendingDenied:
		return a.pending(accountName)
	case model.InactiveATO:
		return a.inactive(accountName)
	case model.InactiveOptOut:
		return a.inactive(accountName)
	case model.InactiveSecurityRisk:
		return a.inactive(accountName)
	case model.InactiveRulesOfEngagement:
		return a.inactive(accountName)
	}

	return nil
}

func (a AccountPolicy) approved(accountName string) error {
	accountUA := accountpap.UserAttributeName(accountName)

	if err := a.graph.Deassign(accountUA, PendingUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the pending user attribute", accountUA)
	}

	if err := a.graph.Deassign(accountUA, InactiveUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the inactive user attribute", accountUA)
	}

	if err := a.graph.Assign(accountUA, ActiveUA); err != nil {
		return errors.Wrapf(err, "error assigning %q to the active user attribute", accountUA)
	}

	return nil
}

func (a AccountPolicy) pending(accountName string) error {
	accountUA := accountpap.UserAttributeName(accountName)

	if err := a.graph.Assign(accountUA, PendingUA); err != nil {
		return errors.Wrapf(err, "error assigning %q to the pending user attribute", accountUA)
	}

	if err := a.graph.Deassign(accountUA, InactiveUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the inactive user attribute", accountUA)
	}

	if err := a.graph.Deassign(accountUA, ActiveUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the active user attribute", accountUA)
	}

	return nil
}

func (a AccountPolicy) inactive(accountName string) error {
	accountUA := accountpap.UserAttributeName(accountName)

	if err := a.graph.Deassign(accountUA, PendingUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the pending user attribute", accountUA)
	}

	if err := a.graph.Assign(accountUA, InactiveUA); err != nil {
		return errors.Wrapf(err, "error assigning %q to the inactive user attribute", accountUA)
	}

	if err := a.graph.Deassign(accountUA, ActiveUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the active user attribute", accountUA)
	}

	return nil
}
