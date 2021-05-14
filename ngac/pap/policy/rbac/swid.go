package rbac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	licensepap "github.com/usnistgov/blossom/chaincode/ngac/pap/license"
)

type SwIDPolicy struct {
	graph pip.Graph
}

func NewSwIDPolicy(graph pip.Graph) SwIDPolicy {
	return SwIDPolicy{graph: graph}
}

func (s SwIDPolicy) ReportSwID(primaryTag string, licenseID string, licenseKey string) error {
	swidNode, err := s.graph.CreateNode(primaryTag, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating node for swid %s ", primaryTag)
	}

	// assign swid node to swid container
	if err = s.graph.Assign(swidNode.Name, SwIDsOA); err != nil {
		return errors.Wrapf(err, "error assiging swid node to swids container")
	}

	// assign the license key object to the swid node
	if err = s.graph.Assign(licensepap.LicenseKeyObject(licenseID, licenseKey), swidNode.Name); err != nil {
		return errors.Wrapf(err, "error assigning the license key to the swid node")
	}

	return nil
}
