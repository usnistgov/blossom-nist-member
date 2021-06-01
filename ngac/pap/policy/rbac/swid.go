package rbac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	assetpap "github.com/usnistgov/blossom/chaincode/ngac/pap/asset"
	swidpap "github.com/usnistgov/blossom/chaincode/ngac/pap/swid"
)

type SwIDPolicy struct {
	graph pip.Graph
}

func NewSwIDPolicy(graph pip.Graph) SwIDPolicy {
	return SwIDPolicy{graph: graph}
}

func (s SwIDPolicy) ReportSwID(primaryTag string, assetID string, license string) error {
	swidNode, err := s.graph.CreateNode(swidpap.ObjectAttributeName(primaryTag), pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating node for swid %s ", primaryTag)
	}

	// assign swid node to swid container
	if err = s.graph.Assign(swidNode.Name, SwIDsOA); err != nil {
		return errors.Wrapf(err, "error assiging swid node to swids container")
	}

	// assign the license key object to the swid node
	if err = s.graph.Assign(assetpap.LicenseObject(assetID, license), swidNode.Name); err != nil {
		return errors.Wrapf(err, "error assigning the license to the swid node")
	}

	return nil
}
