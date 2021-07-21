package pap

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/hyperledger/fabric/core/chaincode/shim"
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

func NewSwIDAdmin(stub shim.ChaincodeStubInterface) (*SwIDAdmin, error) {
	sa := &SwIDAdmin{}
	err := sa.setup(stub)
	return sa, err
}

func (s *SwIDAdmin) setup(stub shim.ChaincodeStubInterface) error {
	graph, err := ledger.GetGraph(stub)
	if err != nil {
		return errors.Wrap(err, "error retrieving ngac graph from ledger")
	}

	s.graph = graph

	return nil
}

func (s *SwIDAdmin) Graph() pip.Graph {
	return s.graph
}

func (s *SwIDAdmin) ReportSwID(stub shim.ChaincodeStubInterface, swid *model.SwID, agency string) error {
	if err := s.setup(stub); err != nil {
		return errors.Wrapf(err, "error setting up agency admin")
	}

	rbacPolicy := rbacpolicy.NewSwIDPolicy(s.graph)
	if err := rbacPolicy.ReportSwID(swid.PrimaryTag, swid.Asset, swid.License); err != nil {
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

	return ledger.UpdateGraphState(stub, s.graph)
}
