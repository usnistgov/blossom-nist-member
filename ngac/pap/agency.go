package pap

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	dacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/dac"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/ledger"
	rbacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/rbac"
	statuspolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/status"
)

type AgencyAdmin struct {
}

func (a AgencyAdmin) RequestAccount(ctx contractapi.TransactionContextInterface, agency model.Agency) error {
	graph, err := ledger.GetGraph(ctx)
	if err != nil {
		return errors.Wrap(err, "error retrieving ngac graph from ledger")
	}

	dacPolicy := dacpolicy.NewAgencyPolicy(graph)
	if err := dacPolicy.RequestAccount(ctx, agency); err != nil {
		return errors.Wrap(err, "error configuring account DAC policy")
	}

	rbacPolicy := rbacpolicy.NewAgencyPolicy(graph)
	if err := rbacPolicy.RequestAccount(ctx, agency); err != nil {
		return errors.Wrap(err, "error configuring account RBAC policy")
	}

	statusPolicy := statuspolicy.NewAgencyPolicy(graph)
	if err := statusPolicy.RequestAccount(ctx, agency); err != nil {
		return errors.Wrap(err, "error configuring account Status policy")
	}

	return ledger.UpdateGraphState(ctx, graph)
}

func (a AgencyAdmin) UpdateAgencyStatus(ctx contractapi.TransactionContextInterface, agency string, status model.Status) error {
	statusPolicy := new(statuspolicy.AgencyPolicy)
	if err := statusPolicy.UpdateAgencyStatus(ctx, agency, status); err != nil {
		return errors.Wrap(err, "error updating agency status")
	}

	return nil
}
