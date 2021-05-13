package pap

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/ledger"
)

type LicenseAdmin struct {
	graph pip.Graph
}

func NewLicenseAdmin() *LicenseAdmin {
	return &LicenseAdmin{}
}

func (l *LicenseAdmin) setup(ctx contractapi.TransactionContextInterface) error {
	graph, err := ledger.GetGraph(ctx)
	if err != nil {
		return errors.Wrap(err, "error retrieving ngac graph from ledger")
	}

	l.graph = graph

	return nil
}

func (l *LicenseAdmin) OnboardLicense(ctx contractapi.TransactionContextInterface, license *model.License) error {
	return nil
}
