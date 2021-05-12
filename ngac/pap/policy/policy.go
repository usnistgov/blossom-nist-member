package policy

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	dac "github.com/usnistgov/blossom/chaincode/ngac/pap/dac"
	rbac "github.com/usnistgov/blossom/chaincode/ngac/pap/rbac"
	status "github.com/usnistgov/blossom/chaincode/ngac/pap/status"
)

const (
	BlossomObject  = "blossom"
	BlossomOA      = "blossom_OA"
	BlossomAdmin   = "Org1 Admin:Org1MSP"
	BlossomAdminUA = "Org1 Admin:Org1MSP_UA"
)

func Configure(graph pip.Graph) error {
	if err := configureSuperPolicy(graph); err != nil {
		return errors.Wrapf(err, "error configuring super policy")
	}

	// configure RBAC policy class
	if err := rbac.Configure(graph, BlossomAdminUA); err != nil {
		return err
	}

	// configure the DAC policy class
	if err := dac.Configure(graph, BlossomAdminUA); err != nil {
		return err
	}

	// configure the status policy class
	if err := status.Configure(graph, BlossomAdminUA); err != nil {
		return err
	}

	return nil
}

func configureSuperPolicy(graph pip.Graph) error {
	blossomPC, err := graph.CreateNode("blossom_PC", pip.PolicyClass, nil)
	if err != nil {
		return errors.Wrap(err, "error creating blossom policy class")
	}

	// create the admin user
	if _, err := graph.CreateNode(BlossomAdmin, pip.User, nil); err != nil {
		return errors.Wrapf(err, "error creating admin user node")
	}

	// create a UA for the admin user
	// this is the node that will be used to set the admin user's policies
	if _, err := graph.CreateNode(BlossomAdminUA, pip.UserAttribute, nil); err != nil {
		return errors.Wrapf(err, "error creating admin user attribute")
	}

	if err = graph.Assign(BlossomAdmin, BlossomAdminUA); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", BlossomAdmin, BlossomAdminUA)
	}

	if err = graph.Assign(BlossomAdminUA, blossomPC.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", BlossomAdmin, BlossomAdminUA)
	}

	// create blossom object and oa
	if _, err = graph.CreateNode(BlossomObject, pip.Object, nil); err != nil {
		return errors.Wrap(err, "error creating blossom object")
	}

	if _, err = graph.CreateNode(BlossomOA, pip.ObjectAttribute, nil); err != nil {
		return errors.Wrap(err, "error creating admin user attribute")
	}

	if err = graph.Assign(BlossomObject, BlossomOA); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", BlossomObject, BlossomOA)
	}

	if err = graph.Assign(BlossomOA, blossomPC.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", BlossomOA, blossomPC.Name)
	}

	if err = graph.Associate(BlossomAdminUA, BlossomOA, pip.ToOps(operations.InitBlossom)); err != nil {
		return errors.Wrapf(err, "error associating bossom admin with blossom object attribute")
	}

	return nil

}
