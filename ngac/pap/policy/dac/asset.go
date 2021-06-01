package dac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/asset"
	"time"
)

type AssetPolicy struct {
	graph pip.Graph
}

func NewAssetPolicy(graph pip.Graph) AssetPolicy {
	return AssetPolicy{graph: graph}
}

func (l AssetPolicy) OnboardAsset(assetOA pip.Node) error {
	// assign the license OA to the dac licenses container
	if err := l.graph.Assign(assetOA.Name, AssetsOA); err != nil {
		return errors.Wrapf(err, "error assigning the asset to the dac assets container")
	}

	return nil
}

func (l AssetPolicy) Checkout(agencyName string, assetID string, licenses map[string]time.Time) error {
	// assign the objects representing the licenses to the agency making the request's DAC object attribute
	for license := range licenses {
		if err := l.graph.Assign(asset.LicenseObject(assetID, license), agencypap.ObjectAttributeName(agencyName)); err != nil {
			return errors.Wrapf(err, "error assigning key %s to agency %s", license, agencyName)
		}
	}

	return nil
}

func (l AssetPolicy) Checkin(agencyName string, assetID string, licenses []string) error {
	// deassign the objects representing the licenses from the agency's DAC object attribute
	for _, license := range licenses {
		if err := l.graph.Deassign(asset.LicenseObject(assetID, license), agencypap.ObjectAttributeName(agencyName)); err != nil {
			return errors.Wrapf(err, "error assigning key %s to agency %s", license, agencyName)
		}
	}

	return nil
}
