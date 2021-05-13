package status

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
)

type AgencyPolicy struct {
	graph pip.Graph
}

func NewAgencyPolicy(graph pip.Graph) AgencyPolicy {
	return AgencyPolicy{graph: graph}
}

func (a AgencyPolicy) RequestAccount(agency model.Agency) error {
	// assign users to pending
	agencyUA := agencypap.UserAttributeName(agency.Name)
	if err := a.graph.Assign(agencyUA, PendingUA); err != nil {
		return errors.Wrap(err, "error assigning the agency user attribute to the pending user attribute")
	}

	// assign agency object to agencies oa
	agencyObj := agencypap.InfoObjectName(agency.Name)
	if err := a.graph.Assign(agencyObj, AgenciesOA); err != nil {
		return errors.Wrap(err, "error assigning the agency object to the agencies object attribute")
	}

	return nil
}

func (a AgencyPolicy) UpdateAgencyStatus(agencyName string, status model.Status) error {
	switch status {
	case model.Approved:
		return a.approved(agencyName)
	case model.PendingApproval:
		return a.pending(agencyName)
	case model.PendingATO:
		return a.pending(agencyName)
	case model.PendingDenied:
		return a.pending(agencyName)
	case model.InactiveATO:
		return a.inactive(agencyName)
	case model.InactiveOptOut:
		return a.inactive(agencyName)
	case model.InactiveSecurityRisk:
		return a.inactive(agencyName)
	case model.InactiveRulesOfEngagement:
		return a.inactive(agencyName)
	}

	return nil
}

func (a AgencyPolicy) approved(agencyName string) error {
	agencyUA := agencypap.UserAttributeName(agencyName)

	if err := a.graph.Deassign(agencyUA, PendingUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the pending user attribute", agencyUA)
	}

	if err := a.graph.Deassign(agencyUA, InactiveUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the inactive user attribute", agencyUA)
	}

	if err := a.graph.Assign(agencyUA, ActiveUA); err != nil {
		return errors.Wrapf(err, "error assigning %q to the active user attribute", agencyUA)
	}

	return nil
}

func (a AgencyPolicy) pending(agencyName string) error {
	agencyUA := agencypap.UserAttributeName(agencyName)

	if err := a.graph.Assign(agencyUA, PendingUA); err != nil {
		return errors.Wrapf(err, "error assigning %q to the pending user attribute", agencyUA)
	}

	if err := a.graph.Deassign(agencyUA, InactiveUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the inactive user attribute", agencyUA)
	}

	if err := a.graph.Deassign(agencyUA, ActiveUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the active user attribute", agencyUA)
	}

	return nil
}

func (a AgencyPolicy) inactive(agencyName string) error {
	agencyUA := agencypap.UserAttributeName(agencyName)

	if err := a.graph.Deassign(agencyUA, PendingUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the pending user attribute", agencyUA)
	}

	if err := a.graph.Assign(agencyUA, InactiveUA); err != nil {
		return errors.Wrapf(err, "error assigning %q to the inactive user attribute", agencyUA)
	}

	if err := a.graph.Deassign(agencyUA, ActiveUA); err != nil {
		return errors.Wrapf(err, "error removing %q from the active user attribute", agencyUA)
	}

	return nil
}
