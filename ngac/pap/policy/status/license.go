package status

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
)

type LicensePolicy struct {
	graph pip.Graph
}

func NewLicensePolicy(graph pip.Graph) LicensePolicy {
	return LicensePolicy{graph: graph}
}

func (l LicensePolicy) OnboardLicense(licenseNode pip.Node) error {
	// assign the license object attribute to the status licenses object attribute
	if err := l.graph.Assign(licenseNode.Name, LicensesOA); err != nil {
		return errors.Wrap(err, "error assigning the license node to the Status policy's licenses object attribute")
	}

	return nil
}
