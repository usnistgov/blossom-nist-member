package status

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
