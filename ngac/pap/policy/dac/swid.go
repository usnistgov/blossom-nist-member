package dac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
)

type SwIDPolicy struct {
	graph pip.Graph
}

func NewSwIDPolicy(graph pip.Graph) SwIDPolicy {
	return SwIDPolicy{graph: graph}
}

func (s SwIDPolicy) ReportSwID(primaryTag string, agencyName string) error {
	// assign the swid object attribute to the agency object attribute
	if err := s.graph.Assign(primaryTag, agencypap.ObjectAttributeName(agencyName)); err != nil {
		return errors.Wrap(err, "error assigning swid node to the agency container")
	}

	return nil
}
