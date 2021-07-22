package dac

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
	assetpap "github.com/usnistgov/blossom/chaincode/ngac/pap/asset"
)

type AgencyPolicy struct {
	graph pip.Graph
}

func NewAgencyPolicy(graph pip.Graph) AgencyPolicy {
	return AgencyPolicy{graph: graph}
}

func (a AgencyPolicy) RequestAccount(agency *model.Agency) error {
	// TODO using the PAP to avoid permission check for now as obligations are not yet implemented.  Once obligations are
	// implemented the admin can create one to create the system owners and when executed will be executed on behalf of
	// the admin not the system owner requesting an account

	if agency.Users.SystemOwner == "" {
		return errors.Errorf("request missing system owner")
	}

	if agency.Users.AcquisitionSpecialist == "" {
		return errors.Errorf("request missing acquisition specialist")
	}

	if agency.Users.SystemAdministrator == "" {
		return errors.Errorf("request missing system administrator")
	}

	// create the system owner user
	systemOwnerNode, err := a.graph.CreateNode(ngac.FormatUsername(agency.Users.SystemOwner, agency.MSPID), pip.User, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating system owner user")
	}

	// create the acquisition specialist user
	acqSpecNode, err := a.graph.CreateNode(ngac.FormatUsername(agency.Users.AcquisitionSpecialist, agency.MSPID), pip.User, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating acquisition specialist user")
	}

	// create the system administrator user
	sysAdminNode, err := a.graph.CreateNode(ngac.FormatUsername(agency.Users.SystemAdministrator, agency.MSPID), pip.User, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating system administrator user")
	}

	// create a DAC UA for this agency
	agencyUA, err := a.graph.CreateNode(agencypap.UserAttributeName(agency.Name), pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating DAC user attribute")
	}

	// assign system owner to agency ua
	if err = a.graph.Assign(systemOwnerNode.Name, agencyUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning SystemOwner to agency DAC user attribute")
	}

	// assign acquisition specialist to agency ua
	if err = a.graph.Assign(acqSpecNode.Name, agencyUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning AcquisitionSpecialist to agency DAC user attribute")
	}

	// assign system administrator to agency ua
	if err = a.graph.Assign(sysAdminNode.Name, agencyUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning SystemAdministrator to agency DAC user attribute")
	}

	// assign agency ua to policy class
	if err = a.graph.Assign(agencyUA.Name, UserAttributeName); err != nil {
		return errors.Wrapf(err, "error assigning agency user attribute to DAC user attribute")
	}

	// create an object attribute for the agency container
	var agencyOA pip.Node
	if agencyOA, err = a.graph.CreateNode(agencypap.ObjectAttributeName(agency.Name), pip.ObjectAttribute, nil); err != nil {
		return fmt.Errorf("error creating agency info object attribute in NGAC: %w", err)
	}

	// assign the agency oa to the dac oa
	if err = a.graph.Assign(agencyOA.Name, ObjectAttributeName); err != nil {
		return fmt.Errorf("error assigning the agency object attribute to the DAC object attribute")
	}

	// create an object to represent the agency
	var agencyInfo pip.Node
	if agencyInfo, err = a.graph.CreateNode(agencypap.InfoObjectName(agency.Name), pip.Object,
		map[string]string{"agency": agency.Name, "type": "agency"}); err != nil {
		return fmt.Errorf("error creating agency info object attribute in NGAC: %w", err)
	}

	// create an object attribute for the agency licenses
	var licensesOA pip.Node
	if licensesOA, err = a.graph.CreateNode(assetpap.AssetsObjectAttribute(agency.Name), pip.ObjectAttribute, nil); err != nil {
		return errors.Wrapf(err, "error creating licenses object attribute for agency %s", agency.Name)
	}

	// assign the agency oa to the dac oa
	if err = a.graph.Assign(agencyInfo.Name, agencyOA.Name); err != nil {
		return fmt.Errorf("error assigning the agency info object to the agency object attribute")
	}

	// assign the licenses OA to the agency oa
	if err = a.graph.Assign(licensesOA.Name, agencyOA.Name); err != nil {
		return fmt.Errorf("error assigning the agency licenses object attribute to the agency object attribute")
	}

	// associate the agency ua and oa with all ops because RBAC will enforce permissions for each user's role
	if err = a.graph.Associate(agencyUA.Name, agencyOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return fmt.Errorf("error associating agency user attribute and agency object attribute")
	}

	return nil
}
