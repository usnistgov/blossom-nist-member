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

type AccountAdmin struct {
	graph pip.Graph
}

func NewAccountAdmin(stub shim.ChaincodeStubInterface) (*AccountAdmin, error) {
	aa := &AccountAdmin{}
	err := aa.setup(stub)
	return aa, err
}

func (a *AccountAdmin) setup(stub shim.ChaincodeStubInterface) error {
	graph, err := ledger.GetGraph(stub)
	if err != nil {
		return errors.Wrap(err, "error retrieving ngac graph from ledger")
	}

	a.graph = graph

	return nil
}

func (a *AccountAdmin) Graph() pip.Graph {
	return a.graph
}

func (a *AccountAdmin) RequestAccount(stub shim.ChaincodeStubInterface, account *model.Account) error {
	if err := a.setup(stub); err != nil {
		return errors.Wrapf(err, "error setting up account admin")
	}

	dacPolicy := dacpolicy.NewAccountPolicy(a.graph)
	if err := dacPolicy.RequestAccount(account); err != nil {
		return errors.Wrap(err, "error configuring account DAC policy")
	}

	rbacPolicy := rbacpolicy.NewAccountPolicy(a.graph)
	if err := rbacPolicy.RequestAccount(account); err != nil {
		return errors.Wrap(err, "error configuring account RBAC policy")
	}

	statusPolicy := statuspolicy.NewAccountPolicy(a.graph)
	if err := statusPolicy.RequestAccount(account); err != nil {
		return errors.Wrap(err, "error configuring account Status policy")
	}

	return ledger.UpdateGraphState(stub, a.graph)
}

func (a *AccountAdmin) UpdateAccountStatus(stub shim.ChaincodeStubInterface, account string, status model.Status) error {
	if err := a.setup(stub); err != nil {
		return errors.Wrapf(err, "error setting up account admin")
	}

	statusPolicy := statuspolicy.NewAccountPolicy(a.graph)
	if err := statusPolicy.UpdateAccountStatus(account, status); err != nil {
		return errors.Wrap(err, "error updating account status")
	}

	return ledger.UpdateGraphState(stub, a.graph)
}
