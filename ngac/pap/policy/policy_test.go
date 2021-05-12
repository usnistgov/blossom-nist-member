package policy

import (
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
	dacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/dac"
	rbacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/rbac"
	statuspolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/status"
	"testing"
)

func TestConfigure(t *testing.T) {
	graph := memory.NewGraph()
	if err := Configure(graph); err != nil {
		t.Fatal(err)
	}

	a := model.Agency{
		Name:  "Org2",
		MSPID: "Org2MSP",
		Users: model.Users{
			SystemOwner:           "system_owner",
			SystemAdministrator:   "sys_admin",
			AcquisitionSpecialist: "acq_spec",
		},
	}

	dacPolicy := dacpolicy.NewAgencyPolicy(graph)
	if err := dacPolicy.RequestAccount(nil, a); err != nil {
		t.Fatal(errors.Wrap(err, "error configuring account DAC policy"))
	}

	rbacPolicy := rbacpolicy.NewAgencyPolicy(graph)
	if err := rbacPolicy.RequestAccount(nil, a); err != nil {
		t.Fatal(errors.Wrap(err, "error configuring account RBAC policy"))
	}

	statusPolicy := statuspolicy.NewAgencyPolicy(graph)
	if err := statusPolicy.RequestAccount(nil, a); err != nil {
		t.Fatal(errors.Wrap(err, "error configuring account Status policy"))
	}

	decider := pdp.NewDecider(graph)
	if ok, err := decider.HasPermissions("system_owner:Org2MSP", agencypap.InfoObjectName("Org2"), operations.ViewAgency); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Fatal("no")
	}
}
