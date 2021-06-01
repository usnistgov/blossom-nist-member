package status

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
)

type AssetPolicy struct {
	graph pip.Graph
}

func NewAssetPolicy(graph pip.Graph) AssetPolicy {
	return AssetPolicy{graph: graph}
}

func (l AssetPolicy) OnboardAsset(licenseNode pip.Node) error {
	// assign the license object attribute to the status licenses object attribute
	if err := l.graph.Assign(licenseNode.Name, AssetsOA); err != nil {
		return errors.Wrap(err, "error assigning the license node to the Status policy's licenses object attribute")
	}

	return nil
}
