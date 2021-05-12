package dac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/usnistgov/blossom/chaincode/model"
)

type LicensePolicy struct {
	graph pip.Graph
}

func NewLicensePolicy(graph pip.Graph) LicensePolicy {
	return LicensePolicy{graph: graph}
}

func (l LicensePolicy) OnboardLicense(ctx contractapi.TransactionContextInterface, license *model.License) error {
	return nil
}

func (l LicensePolicy) OffboardLicense(ctx contractapi.TransactionContextInterface, licenseID string) error {
	return nil
}

func (l LicensePolicy) Licenses(ctx contractapi.TransactionContextInterface) ([]*model.License, error) {
	return nil, nil
}

func (l LicensePolicy) LicenseInfo(ctx contractapi.TransactionContextInterface, licenseID string) (*model.License, error) {
	return nil, nil
}

func (l LicensePolicy) LicenseKeys(ctx contractapi.TransactionContextInterface, licenseID string) (map[string][]string, error) {
	return nil, nil
}

func (l LicensePolicy) AgencyLicenseKeys(ctx contractapi.TransactionContextInterface, agency string) (map[int][]string, error) {
	return nil, nil
}

func (l LicensePolicy) CheckoutLicense(ctx contractapi.TransactionContextInterface, licenseID string, agency string, amount int) ([]string, error) {
	// get the agency OA
	/*agencyOA, err := l.graph.GetNode(agencypap.ObjectAttributeName(agency))
	if err != nil {
		return nil, errors.Wrap(err, "error getting agency object attribute")
	}

	// get the license object attribute
	licenseOA, err := l.graph.GetNode(licenseID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting license object attribute")
	}*/

	// get the children of the license object attribute
	return nil, nil
}

func (l LicensePolicy) CheckinLicense(ctx contractapi.TransactionContextInterface, licenseID string) error {
	return nil
}
