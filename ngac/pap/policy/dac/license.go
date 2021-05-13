package dac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/license"
)

type LicensePolicy struct {
	graph pip.Graph
}

func NewLicensePolicy(graph pip.Graph) LicensePolicy {
	return LicensePolicy{graph: graph}
}

func (l LicensePolicy) CheckoutLicense(agencyName string, licenseID string, keys []string) error {
	// assign the objects representing the keys to the agency making the request's DAC object attribute
	for _, key := range keys {
		if err := l.graph.Assign(license.LicenseKeyObject(licenseID, key), agencypap.ObjectAttributeName(agencyName)); err != nil {
			return errors.Wrapf(err, "error assigning key %s to agency %s", key, agencyName)
		}
	}

	return nil
}

func (l LicensePolicy) CheckinLicense(agencyName string, licenseID string, keys []string) error {
	// deassign the objects representing the keys from the agency's DAC object attribute
	for _, key := range keys {
		if err := l.graph.Deassign(license.LicenseKeyObject(licenseID, key), agencypap.ObjectAttributeName(agencyName)); err != nil {
			return errors.Wrapf(err, "error assigning key %s to agency %s", key, agencyName)
		}
	}

	return nil
}
