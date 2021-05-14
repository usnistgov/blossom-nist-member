package dac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
)

const (
	PolicyClassName     = "DAC"
	ObjectAttributeName = "DAC_OA"
	UserAttributeName   = "DAC_UA"
	LicensesOA          = "dac_licenses"
)

func Configure(graph pip.Graph, adminUA string) error {
	// create DAC policy class node
	dacPC, err := graph.CreateNode(PolicyClassName, pip.PolicyClass, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating DAC policy class node")
	}

	// DAC default nodes
	dacUA, err := graph.CreateNode(UserAttributeName, pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating DAC user attribute node")
	}

	if err = graph.Assign(dacUA.Name, dacPC.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", dacUA.Name, dacPC.Name)
	}

	dacOA, err := graph.CreateNode(ObjectAttributeName, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating DAC object attribute node")
	}

	if err = graph.Assign(dacOA.Name, dacPC.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", dacUA.Name, dacPC.Name)
	}

	if err = graph.Associate(adminUA, dacUA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", adminUA, dacUA.Name)
	}

	if err = graph.Associate(adminUA, dacOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", adminUA, dacOA.Name)
	}

	// create a container for all licenses
	licensesOA, err := graph.CreateNode(LicensesOA, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating licenses container in DAC policy")
	}

	// assign the licenses container to the dac OA
	if err = graph.Assign(licensesOA.Name, dacOA.Name); err != nil {
		return errors.Wrapf(err, "error assignign licenses container to DAC object attribute")
	}

	// associate the org1 admin ua with * on licenses container
	if err = graph.Associate(adminUA, licensesOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating admin user attribute with licenses container in DAC policy")
	}

	// associate dac UA with container with checkin/checkout permissions
	if err = graph.Associate(dacUA.Name, licensesOA.Name, pip.ToOps(operations.ViewLicense, operations.CheckOutLicense,
		operations.CheckInLicense)); err != nil {
		return errors.Wrapf(err, "error associating dac user attribute with license container in DAC policy")
	}

	return nil
}
