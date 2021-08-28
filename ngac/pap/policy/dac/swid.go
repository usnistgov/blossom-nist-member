package dac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	accountpap "github.com/usnistgov/blossom/chaincode/ngac/pap/account"
	swidpap "github.com/usnistgov/blossom/chaincode/ngac/pap/swid"
)

type SwIDPolicy struct {
	graph pip.Graph
}

func NewSwIDPolicy(graph pip.Graph) SwIDPolicy {
	return SwIDPolicy{graph: graph}
}

func (s SwIDPolicy) ReportSwID(primaryTag string, accountName string) error {
	// assign the swid object attribute to the account object attribute
	if err := s.graph.Assign(swidpap.ObjectAttributeName(primaryTag), accountpap.ObjectAttributeName(accountName)); err != nil {
		return errors.Wrap(err, "error assigning swid node to the account container")
	}

	return nil
}
