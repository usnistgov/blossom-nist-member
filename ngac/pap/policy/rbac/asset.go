package rbac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	assetpap "github.com/usnistgov/blossom/chaincode/ngac/pap/asset"
)

type AssetPolicy struct {
	graph pip.Graph
}

func NewAssetPolicy(graph pip.Graph) AssetPolicy {
	return AssetPolicy{graph: graph}
}

func (l AssetPolicy) OnboardAsset(asset *model.Asset) (pip.Node, error) {
	// create an object attribute for the asset
	var (
		assetOA pip.Node
		err     error
	)

	if assetOA, err = l.graph.CreateNode(assetpap.ObjectAttribute(asset.ID), pip.ObjectAttribute, nil); err != nil {
		return pip.Node{}, errors.Wrapf(err, "error creating object attribute for asset %s", asset.ID)
	}

	// create objects for the licenses and assign to the asset
	for _, key := range asset.Licenses {
		var keyNode pip.Node
		if keyNode, err = l.graph.CreateNode(assetpap.LicenseObject(asset.ID, key), pip.Object, nil); err != nil {
			return pip.Node{}, errors.Wrapf(err, "error creating object for license %s", key)
		}

		if err = l.graph.Assign(keyNode.Name, assetOA.Name); err != nil {
			return pip.Node{}, errors.Wrapf(err, "error assigning licebse object %s to asset %s object attribute", key, asset.ID)
		}
	}

	// assign the asset OA to the RBAC assets OA
	if err = l.graph.Assign(assetOA.Name, AssetsOA); err != nil {
		return pip.Node{}, errors.Wrapf(err, "error assigning asset %s object attribute to Assets object attribute", asset.ID)
	}

	return assetOA, nil
}

func (l AssetPolicy) OffboardAsset(assetID string) error {
	// delete key objects
	licenseNodes, err := l.graph.GetChildren(assetpap.ObjectAttribute(assetID))
	if err != nil {
		return errors.Wrapf(err, "error getting license nodes of asset %s", assetID)
	}

	for licenseNode := range licenseNodes {
		if err = l.graph.DeleteNode(licenseNode); err != nil {
			return errors.Wrapf(err, "error deleting license node %s", licenseNode)
		}
	}

	// delete license oa
	if err = l.graph.DeleteNode(assetpap.ObjectAttribute(assetID)); err != nil {
		return errors.Wrapf(err, "error deleting asset object attribute")
	}

	return nil
}
