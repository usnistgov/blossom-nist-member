package dac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
)

const (
	PolicyClassName     = "DAC"
	ObjectAttributeName = "DAC_OA"
	UserAttributeName   = "DAC_UA"
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

	return nil
}
