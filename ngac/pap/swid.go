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

type SwIDAdmin struct {
	graph pip.Graph
}

func NewSwIDAdmin(ctx contractapi.TransactionContextInterface) (*SwIDAdmin, error) {
	sa := &SwIDAdmin{}
	err := sa.setup(ctx)
	return sa, err
}

func (s *SwIDAdmin) setup(ctx contractapi.TransactionContextInterface) error {
	graph, err := ledger.GetGraph(ctx)
	if err != nil {
		return errors.Wrap(err, "error retrieving ngac graph from ledger")
	}

	s.graph = graph

	return nil
}

func (s *SwIDAdmin) Graph() pip.Graph {
	return s.graph
}

func (s *SwIDAdmin) ReportSwID(ctx contractapi.TransactionContextInterface, swid *model.SwID, agency string) error {
	if err := s.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency admin")
	}

	rbacPolicy := rbacpolicy.NewSwIDPolicy(s.graph)
	if err := rbacPolicy.ReportSwID(swid.PrimaryTag, swid.License, swid.LicenseKey); err != nil {
		return errors.Wrap(err, "error configuring swid RBAC policy")
	}

	dacPolicy := dacpolicy.NewSwIDPolicy(s.graph)
	if err := dacPolicy.ReportSwID(swid.PrimaryTag, agency); err != nil {
		return errors.Wrap(err, "error configuring swid DAC policy")
	}

	statusPolicy := statuspolicy.NewSwIDPolicy(s.graph)
	if err := statusPolicy.ReportSwID(swid.PrimaryTag); err != nil {
		return errors.Wrap(err, "error configuring swid Status policy")
	}

	return ledger.UpdateGraphState(ctx, s.graph)
}
