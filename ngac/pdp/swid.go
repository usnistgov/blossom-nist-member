package pdp

import (
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
	assetpap "github.com/usnistgov/blossom/chaincode/ngac/pap/asset"
	swidpap "github.com/usnistgov/blossom/chaincode/ngac/pap/swid"
	"time"
)

// SwIDDecider is the Policy Decision Point (PDP) for the SwID API
type SwIDDecider struct {
	// user is the user that is currently executing a function
	user string
	// pap is the policy administration point for agencies
	pap *pap.SwIDAdmin
	// decider is the NGAC decider used to make decisions
	decider pdp.Decider
}

// NewSwIDDecider creates a new SwIDDecider with the user from the ctx and a NGAC Decider using the NGAC graph
// from the ledger.
func NewSwIDDecider() *SwIDDecider {
	return &SwIDDecider{}
}

func (s *SwIDDecider) setup(ctx contractapi.TransactionContextInterface) error {
	user, err := GetUser(ctx)
	if err != nil {
		return errors.Wrapf(err, "error getting user from request")
	}

	s.user = user

	// initialize the agency policy administration point
	s.pap, err = pap.NewSwIDAdmin(ctx)
	if err != nil {
		return errors.Wrapf(err, "error initializing agency administraion point")
	}

	s.decider = pdp.NewDecider(s.pap.Graph())

	return nil
}

func (s *SwIDDecider) ReportSwID(ctx contractapi.TransactionContextInterface, swid *model.SwID, agency string) error {
	if err := s.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up swid decider")
	}

	// check user can assign report swid on the license key object
	if ok, err := s.decider.HasPermissions(s.user, assetpap.LicenseObject(swid.Asset, swid.License), operations.ReportSwid); err != nil {
		return errors.Wrapf(err, "error checking if user has permission to report swid")
	} else if !ok {
		return ErrAccessDenied
	}

	err := s.pap.ReportSwID(ctx, swid, agency)
	return errors.Wrapf(err, "error reporting swid %s", swid.PrimaryTag)
}

func (s *SwIDDecider) FilterSwID(ctx contractapi.TransactionContextInterface, swid *model.SwID) error {
	if err := s.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up swid decider")
	}

	return s.filterSwID(swid)
}

func (s *SwIDDecider) FilterSwIDs(ctx contractapi.TransactionContextInterface, swids []*model.SwID) ([]*model.SwID, error) {
	if err := s.setup(ctx); err != nil {
		return nil, errors.Wrapf(err, "error setting up swid decider")
	}

	filteredSwids := make([]*model.SwID, 0)
	for _, swid := range swids {
		if err := s.filterSwID(swid); err != nil {
			return nil, errors.Wrapf(err, "error filtering swids")
		}

		if swid.PrimaryTag == "" {
			continue
		}

		filteredSwids = append(filteredSwids, swid)
	}

	return filteredSwids, nil
}

func (s *SwIDDecider) filterSwID(swid *model.SwID) error {
	permissions, err := s.decider.ListPermissions(s.user, swidpap.ObjectAttributeName(swid.PrimaryTag))
	if err != nil {
		return errors.Wrapf(err, "error getting permissions for user %s on swid %s", s.user, swid.PrimaryTag)
	}

	if !permissions.Contains(operations.ViewSwID) {
		swid.PrimaryTag = ""
		swid.XML = ""
		swid.Asset = ""
		swid.License = ""
		swid.LeaseExpiration = time.Time{}
	}

	return nil
}
