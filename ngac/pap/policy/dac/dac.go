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
	AssetsOA            = "dac_assets"
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

	// create a container for all assets
	assetsOA, err := graph.CreateNode(AssetsOA, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating licenses container in DAC policy")
	}

	// assign the assets container to the dac OA
	if err = graph.Assign(assetsOA.Name, dacOA.Name); err != nil {
		return errors.Wrapf(err, "error assigning assets container to DAC object attribute")
	}

	// associate the org1 admin ua with * on licenses container
	if err = graph.Associate(adminUA, assetsOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating admin user attribute with assets container in DAC policy")
	}

	// associate dac UA with assets container with checkin/checkout permissions
	if err = graph.Associate(dacUA.Name, assetsOA.Name, pip.ToOps(operations.ViewAsset, operations.CheckOut,
		operations.CheckIn)); err != nil {
		return errors.Wrapf(err, "error associating dac user attribute with assets container in DAC policy")
	}

	return nil
}
