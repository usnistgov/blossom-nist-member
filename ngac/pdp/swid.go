package pdp

import (
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
	licensepap "github.com/usnistgov/blossom/chaincode/ngac/pap/license"
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
	if ok, err := s.decider.HasPermissions(s.user, licensepap.LicenseKeyObject(swid.License, swid.LicenseKey), operations.ReportSwid); err != nil {
		return errors.Wrapf(err, "error checking if user has permission to report swid")
	} else if !ok {
		return ErrAccessDenied
	}

	err := s.pap.ReportSwID(ctx, swid, agency)
	return errors.Wrapf(err, "error reporting swid %s", swid.PrimaryTag)
}
