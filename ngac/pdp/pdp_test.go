package pdp

import (
	"fmt"
	"testing"

	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
)

func TestUpdateGraph(t *testing.T) {
	pdp := new(PDP)

	ledgerGraph := memory.NewGraph()
	ledgerGraph.CreateNode("pc1", pip.PolicyClass, nil)
	jsonGraph := memory.NewGraph()
	jsonGraph.CreateNode("pc2", pip.PolicyClass, nil)

	pdp.UpdateGraph(nil, ledgerGraph, jsonGraph)
	fmt.Println(ledgerGraph.GetNodes())
}
