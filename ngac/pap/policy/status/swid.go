package status

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
)

type SwIDPolicy struct {
	graph pip.Graph
}

func NewSwIDPolicy(graph pip.Graph) SwIDPolicy {
	return SwIDPolicy{graph: graph}
}

func (s SwIDPolicy) ReportSwID(primaryTag string) error {
	err := s.graph.Assign(primaryTag, SwidsOA)
	return errors.Wrap(err, "error assigning swid node to swid container in status policy class")
}
