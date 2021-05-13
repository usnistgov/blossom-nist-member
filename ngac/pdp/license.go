package pdp

import (
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/ledger"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/rbac"
)

type LicenseDecider struct {
	user    string
	decider pdp.Decider
}

// NewLicenseDecider creates a new LicenseDecider with the user from the ctx and a NGAC Decider using the NGAC graph
// from the ledger.
func NewLicenseDecider() LicenseDecider {
	return LicenseDecider{}
}

func (l LicenseDecider) setup(ctx contractapi.TransactionContextInterface) error {
	user, err := GetUser(ctx)
	if err != nil {
		return errors.Wrapf(err, "error getting user from request")
	}

	l.user = user

	graph, err := ledger.GetGraph(ctx)
	if err != nil {
		return errors.Wrap(err, "error retrieving ngac graph from ledger")
	}

	l.decider = pdp.NewDecider(graph)

	return nil
}

func (l LicenseDecider) FilterLicense(ctx contractapi.TransactionContextInterface, license *model.License) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	return l.filterLicense(license)
}

func (l LicenseDecider) filterLicense(license *model.License) error {
	permissions, err := l.decider.ListPermissions(l.user, license.ID)
	if err != nil {
		return errors.Wrapf(err, "error getting permissions for user %s on license %s", l.user, license.Name)
	}

	// if the user cannot view license on the license object attribute, return an empty license
	if !permissions.Contains(operations.ViewLicense) {
		license = &model.License{}
		return nil
	}

	if !permissions.Contains(operations.ViewAllKeys) {
		license.AllKeys = make([]string, 0)
	}

	if !permissions.Contains(operations.ViewAvailableKeys) {
		license.AvailableKeys = make([]string, 0)
	}

	return nil
}

func (l LicenseDecider) FilterLicenses(ctx contractapi.TransactionContextInterface, licenses []*model.License) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	for _, license := range licenses {
		if err := l.filterLicense(license); err != nil {
			return errors.Wrapf(err, "error filtering license")
		}
	}

	return nil
}

func (l LicenseDecider) OnboardLicense(ctx contractapi.TransactionContextInterface, license *model.License) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	// check user can onboard license
	if ok, err := l.decider.HasPermissions(l.user, rbac.LicensesOA, operations.OnboardLicense); err != nil {
		return errors.Wrapf(err, "error checking if user %s can onboard a license", l.user)
	} else if !ok {
		return ErrAccessDenied
	}

	licenseAdmin := pap.NewLicenseAdmin()
	return licenseAdmin.OnboardLicense(ctx, license)
}

func (l LicenseDecider) OffboardLicense(ctx contractapi.TransactionContextInterface, licenseID string) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	// check user can onboard license
	if ok, err := l.decider.HasPermissions(l.user, licenseID, operations.OffboardLicense); err != nil {
		return errors.Wrapf(err, "error checking if user %s can offboard a license", l.user)
	} else if !ok {
		return ErrAccessDenied
	}

	return nil
}

func (l LicenseDecider) CheckoutLicense(ctx contractapi.TransactionContextInterface, licenseID string) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	// check user can checkout license
	if ok, err := l.decider.HasPermissions(l.user, licenseID, operations.OffboardLicense); err != nil {
		return errors.Wrapf(err, "error checking if user %s can offboard a license", l.user)
	} else if !ok {
		return ErrAccessDenied
	}

	return nil
}

func (l LicenseDecider) CheckinLicense(ctx contractapi.TransactionContextInterface, licenseID string) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	return nil
}
