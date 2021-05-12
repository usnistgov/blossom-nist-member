package pdp

import (
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/ledger"
)

type AgencyDecider struct {
	user    string
	decider pdp.Decider
}

// NewAgencyDecider creates a new AgencyDecider with the user from the ctx and a NGAC Decider using the NGAC graph
// from the ledger.
func NewAgencyDecider() AgencyDecider {
	return AgencyDecider{}
}

func (a AgencyDecider) setup(ctx contractapi.TransactionContextInterface) error {
	if a.user == "" {
		user, err := GetUser(ctx)
		if err != nil {
			return errors.Wrapf(err, "error getting user from request")
		}

		a.user = user
	}

	if a.decider == nil {
		graph, err := ledger.GetGraph(ctx)
		if err != nil {
			return errors.Wrap(err, "error retrieving ngac graph from ledger")
		}

		a.decider = pdp.NewDecider(graph)
	}

	return nil
}

func (a AgencyDecider) FilterAgency(ctx contractapi.TransactionContextInterface, agency *model.Agency) error {
	if err := a.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	return a.filterAgency(agency)
}

func (a AgencyDecider) filterAgency(agency *model.Agency) error {
	permissions, err := a.decider.ListPermissions(a.user, agencypap.InfoObjectName(agency.Name))
	if err != nil {
		return errors.Wrapf(err, "error getting permissions for user %s on agency %s", a.user, agency.Name)
	}

	// if the user cannot view agency on the agency info object, return an empty agency
	if !permissions.Contains(operations.ViewAgency) {
		agency = &model.Agency{}
		return nil
	}

	if !permissions.Contains(operations.ViewATO) {
		agency.ATO = ""
	}

	if !permissions.Contains(operations.ViewMSPID) {
		agency.MSPID = ""
	}

	if !permissions.Contains(operations.ViewUsers) {
		agency.Users = model.Users{}
	}

	if !permissions.Contains(operations.ViewStatus) {
		agency.Status = ""
	}

	if !permissions.Contains(operations.ViewAgencyLicenses) {
		agency.Status = ""
	}

	return nil
}

func (a AgencyDecider) FilterAgencies(ctx contractapi.TransactionContextInterface, agencies []*model.Agency) error {
	if err := a.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	for _, agency := range agencies {
		if err := a.filterAgency(agency); err != nil {
			return errors.Wrapf(err, "error filtering agency")
		}
	}

	return nil
}

func (a AgencyDecider) RequestAccount(ctx contractapi.TransactionContextInterface, agency model.Agency) error {
	// any user can create an account
	agencyAdmin := new(pap.AgencyAdmin)
	return agencyAdmin.RequestAccount(ctx, agency)
}

func (a AgencyDecider) UploadATO(ctx contractapi.TransactionContextInterface, agency string, ato string) error {
	if err := a.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	if ok, err := a.decider.HasPermissions(a.user, agencypap.InfoObjectName(agency), operations.UploadATO); err != nil {
		return errors.Wrapf(err, "error checking if user %s can upload an ATO for agency %s", a.user, agency)
	} else if !ok {
		return ErrAccessDenied
	}

	return nil
}

func (a AgencyDecider) UpdateAgencyStatus(ctx contractapi.TransactionContextInterface, agency string, status model.Status) error {
	if err := a.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	if ok, err := a.decider.HasPermissions(a.user, agencypap.InfoObjectName(agency), operations.UpdateAgencyStatus); err != nil {
		return errors.Wrapf(err, "error checking if user %s can update status of agency %s", a.user, agency)
	} else if !ok {
		return ErrAccessDenied
	}

	agencyAdmin := new(pap.AgencyAdmin)
	return agencyAdmin.UpdateAgencyStatus(ctx, agency, status)
}
