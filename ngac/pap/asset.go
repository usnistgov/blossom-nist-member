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
	"time"
)

type AssetAdmin struct {
	graph pip.Graph
}

func NewAssetAdmin(ctx contractapi.TransactionContextInterface) (*AssetAdmin, error) {
	la := &AssetAdmin{}
	err := la.setup(ctx)
	return la, err
}

func (l *AssetAdmin) setup(ctx contractapi.TransactionContextInterface) error {
	graph, err := ledger.GetGraph(ctx)
	if err != nil {
		return errors.Wrap(err, "error retrieving ngac graph from ledger")
	}

	l.graph = graph

	return nil
}

func (l *AssetAdmin) Graph() pip.Graph {
	return l.graph
}

func (l *AssetAdmin) OnboardAsset(ctx contractapi.TransactionContextInterface, asset *model.Asset) error {
	var (
		assetOA pip.Node
		err     error
	)

	rbacPolicy := rbacpolicy.NewAssetPolicy(l.graph)
	if assetOA, err = rbacPolicy.OnboardAsset(asset); err != nil {
		return errors.Wrap(err, "error configuring asset onboard RBAC policy")
	}

	dacPolicy := dacpolicy.NewAssetPolicy(l.graph)
	if err = dacPolicy.OnboardAsset(assetOA); err != nil {
		return errors.Wrap(err, "error configuring asset onboard RBAC policy")
	}

	statusPolicy := statuspolicy.NewAssetPolicy(l.graph)
	if err = statusPolicy.OnboardAsset(assetOA); err != nil {
		return errors.Wrap(err, "error configuring asset onboard Status policy")
	}

	return ledger.UpdateGraphState(ctx, l.graph)
}

func (l *AssetAdmin) OffboardAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	rbacPolicy := rbacpolicy.NewAssetPolicy(l.graph)
	if err := rbacPolicy.OffboardAsset(assetID); err != nil {
		return errors.Wrap(err, "error configuring asset offboard RBAC policy")
	}

	return ledger.UpdateGraphState(ctx, l.graph)
}

func (l *AssetAdmin) Checkout(ctx contractapi.TransactionContextInterface, agencyName string, assetID string,
	licenses map[string]time.Time) error {
	dacPolicy := dacpolicy.NewAssetPolicy(l.graph)
	if err := dacPolicy.Checkout(agencyName, assetID, licenses); err != nil {
		return errors.Wrap(err, "error checking out asset under the DAC policy")
	}

	return ledger.UpdateGraphState(ctx, l.graph)
}

func (l *AssetAdmin) Checkin(ctx contractapi.TransactionContextInterface, agencyName string, assetID string,
	licenses []string) error {
	dacPolicy := dacpolicy.NewAssetPolicy(l.graph)
	if err := dacPolicy.Checkin(agencyName, assetID, licenses); err != nil {
		return errors.Wrap(err, "error checking in asset under the DAC policy")
	}

	return ledger.UpdateGraphState(ctx, l.graph)
}
