package rbac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	licensepap "github.com/usnistgov/blossom/chaincode/ngac/pap/license"
)

type LicensePolicy struct {
	graph pip.Graph
}

func NewLicensePolicy(graph pip.Graph) LicensePolicy {
	return LicensePolicy{graph: graph}
}

func (l LicensePolicy) OnboardLicense(license *model.License) (pip.Node, error) {
	// create an object attribute for the license
	var (
		licenseOA pip.Node
		err       error
	)

	if licenseOA, err = l.graph.CreateNode(licensepap.LicenseObjectAttribute(license.ID), pip.ObjectAttribute, nil); err != nil {
		return pip.Node{}, errors.Wrapf(err, "error creating object attribute for license %s", license.ID)
	}

	// create objects for the keys and assign to the license
	for _, key := range license.AllKeys {
		var keyNode pip.Node
		if keyNode, err = l.graph.CreateNode(licensepap.LicenseKeyObject(license.ID, key), pip.Object, nil); err != nil {
			return pip.Node{}, errors.Wrapf(err, "error creating object for key %s", key)
		}

		if err = l.graph.Assign(keyNode.Name, licenseOA.Name); err != nil {
			return pip.Node{}, errors.Wrapf(err, "error assigning key object %s to license %s object attribute", key, license.ID)
		}
	}

	// assign the license OA to the RBAC licenses OA
	if err = l.graph.Assign(licenseOA.Name, LicensesOA); err != nil {
		return pip.Node{}, errors.Wrapf(err, "error assigning license %s object attribute to Licenses object attribute", license.ID)
	}

	return licenseOA, nil
}

func (l LicensePolicy) OffboardLicense(licenseID string) error {
	// delete key objects
	keyNodes, err := l.graph.GetChildren(licensepap.LicenseObjectAttribute(licenseID))
	if err != nil {
		return errors.Wrapf(err, "error getting key nodes of license %s", licenseID)
	}

	for keyNode := range keyNodes {
		if err = l.graph.DeleteNode(keyNode); err != nil {
			return errors.Wrapf(err, "error deleting key node %s", keyNode)
		}
	}

	// delete license oa
	if err = l.graph.DeleteNode(licensepap.LicenseObjectAttribute(licenseID)); err != nil {
		return errors.Wrapf(err, "error deleting license object attribute")
	}

	return nil
}
