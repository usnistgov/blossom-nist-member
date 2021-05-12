package rbac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
)

var SystemOwnerPermissions = pip.ToOps(
	operations.ViewAgency,
	operations.ViewLicense,
	operations.UploadATO,
	operations.ViewATO,
	operations.ViewMSPID,
	operations.ViewUsers,
	operations.ViewStatus)

var SystemAdminPermissions = pip.ToOps(
	operations.ViewLicense,
	operations.CheckOutLicense,
	operations.CheckInLicense)

var AcqSpecPermissions = pip.ToOps(
	operations.ViewLicense,
	operations.ViewAgency,
	operations.ViewStatus)

type AgencyPolicy struct {
	graph pip.Graph
}

func NewAgencyPolicy(graph pip.Graph) AgencyPolicy {
	return AgencyPolicy{graph: graph}
}

func (a AgencyPolicy) RequestAccount(ctx contractapi.TransactionContextInterface, agency model.Agency) error {
	// assign the agency object to the agencies attribute
	if err := a.graph.Assign(agencypap.InfoObjectName(agency.Name), AgenciesOA); err != nil {
		return errors.Wrapf(err, "error assigning agency %q to agencies attribute", agency.Name)
	}

	// assign the users to their attributes
	if err := a.graph.Assign(ngac.FormatUsername(agency.Users.SystemOwner, agency.MSPID), SystemOwnerUA); err != nil {
		return errors.Wrapf(err, "error assigning system owner to SystemOwner user attribute")
	}

	if err := a.graph.Assign(ngac.FormatUsername(agency.Users.AcquisitionSpecialist, agency.MSPID), AcquisitionSpecialistUA); err != nil {
		return errors.Wrapf(err, "error assigning acquisition specialist to AcquisitionSpecialist user attribute")
	}

	if err := a.graph.Assign(ngac.FormatUsername(agency.Users.SystemAdministrator, agency.MSPID), SystemAdministratorUA); err != nil {
		return errors.Wrapf(err, "error assigning system owner to SystemOwner user attribute")
	}

	return nil
}
