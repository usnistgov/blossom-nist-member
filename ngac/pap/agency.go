package pap

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/ledger"
	dacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/policy/dac"
	rbacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/policy/rbac"
	statuspolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/policy/status"
)

type AgencyAdmin struct {
	graph pip.Graph
}

func NewAgencyAdmin(ctx contractapi.TransactionContextInterface) (*AgencyAdmin, error) {
	aa := &AgencyAdmin{}
	err := aa.setup(ctx)
	return aa, err
}

func (a *AgencyAdmin) setup(ctx contractapi.TransactionContextInterface) error {
	graph, err := ledger.GetGraph(ctx)
	if err != nil {
		return errors.Wrap(err, "error retrieving ngac graph from ledger")
	}

	a.graph = graph

	return nil
}

func (a *AgencyAdmin) Graph() pip.Graph {
	return a.graph
}

func (a *AgencyAdmin) RequestAccount(ctx contractapi.TransactionContextInterface, agency model.Agency) error {
	if err := a.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency admin")
	}

	dacPolicy := dacpolicy.NewAgencyPolicy(a.graph)
	if err := dacPolicy.RequestAccount(agency); err != nil {
		return errors.Wrap(err, "error configuring account DAC policy")
	}

	rbacPolicy := rbacpolicy.NewAgencyPolicy(a.graph)
	if err := rbacPolicy.RequestAccount(agency); err != nil {
		return errors.Wrap(err, "error configuring account RBAC policy")
	}

	statusPolicy := statuspolicy.NewAgencyPolicy(a.graph)
	if err := statusPolicy.RequestAccount(agency); err != nil {
		return errors.Wrap(err, "error configuring account Status policy")
	}

	return ledger.UpdateGraphState(ctx, a.graph)
}

func (a *AgencyAdmin) UpdateAgencyStatus(ctx contractapi.TransactionContextInterface, agency string, status model.Status) error {
	if err := a.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency admin")
	}

	statusPolicy := statuspolicy.NewAgencyPolicy(a.graph)
	if err := statusPolicy.UpdateAgencyStatus(agency, status); err != nil {
		return errors.Wrap(err, "error updating agency status")
	}

	return ledger.UpdateGraphState(ctx, a.graph)
}
