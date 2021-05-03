package asset

import (
	"github.com/usnistgov/blossom/chaincode/asset/operations"
	"testing"

	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/PM-Master/policy-machine-go/pip"
)

func TestInitNGAC(t *testing.T) {
	graph, err := initGraph()
	if err != nil {
		t.Fatal(err)
	}

	decider := pdp.NewDecider(graph)
	if ok, err := decider.Decide("A0admin", "licenses", pip.AllOps); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Fatalf("expected user A0admin to have %s but it did not", operations.ViewLicense)
	}
}
