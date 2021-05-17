package pdp

import (
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
	licensepap "github.com/usnistgov/blossom/chaincode/ngac/pap/license"
	rbacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/policy/rbac"
	"time"
)

type LicenseDecider struct {
	// user is the user that is currently executing a function
	user string
	// pap is the policy administration point for licenses
	pap *pap.LicenseAdmin
	// decider is the NGAC decider used to make decisions
	decider pdp.Decider
}

// NewLicenseDecider creates a new LicenseDecider with the user from the ctx and a NGAC Decider using the NGAC graph
// from the ledger.
func NewLicenseDecider() *LicenseDecider {
	return &LicenseDecider{}
}

func (l *LicenseDecider) setup(ctx contractapi.TransactionContextInterface) error {
	user, err := GetUser(ctx)
	if err != nil {
		return errors.Wrapf(err, "error getting user from request")
	}

	l.user = user

	// initialize the license policy administration point
	l.pap, err = pap.NewLicenseAdmin(ctx)
	if err != nil {
		return errors.Wrapf(err, "error initializing agency administraion point")
	}

	l.decider = pdp.NewDecider(l.pap.Graph())

	return nil
}

func (l *LicenseDecider) FilterLicense(ctx contractapi.TransactionContextInterface, license *model.License) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up license decider")
	}

	return l.filterLicense(license)
}

func (l *LicenseDecider) filterLicense(license *model.License) error {
	permissions, err := l.decider.ListPermissions(l.user, licensepap.LicenseObjectAttribute(license.ID))
	if err != nil {
		return errors.Wrapf(err, "error getting permissions for user %s on license %s", l.user, license.Name)
	}

	if !permissions.Contains(operations.ViewLicense) {
		// if the user cannot view license on the license object attribute, return an empty license
		// initialize array and map values to avoid fabric schema errors
		license = &model.License{
			ID:             "",
			Name:           "",
			TotalAmount:    0,
			Available:      0,
			Cost:           0,
			OnboardingDate: time.Time{},
			Expiration:     time.Time{},
			AllKeys:        make([]string, 0),
			AvailableKeys:  make([]string, 0),
			CheckedOut:     make(map[string]map[string]time.Time),
		}
		return nil
	}

	if !permissions.Contains(operations.ViewAllKeys) {
		license.AllKeys = make([]string, 0)
	}

	if !permissions.Contains(operations.ViewAvailableKeys) {
		license.AvailableKeys = make([]string, 0)
	}

	if !permissions.Contains(operations.ViewCheckedOut) {
		license.CheckedOut = make(map[string]map[string]time.Time)
	}

	return nil
}

func (l *LicenseDecider) FilterLicenses(ctx contractapi.TransactionContextInterface, licenses []*model.License) ([]*model.License, error) {
	if err := l.setup(ctx); err != nil {
		return nil, errors.Wrapf(err, "error setting up agency decider")
	}

	filteredLicenses := make([]*model.License, 0)
	for _, license := range licenses {
		if err := l.filterLicense(license); err != nil {
			return nil, errors.Wrapf(err, "error filtering license")
		}

		if license.ID == "" {
			continue
		}

		filteredLicenses = append(filteredLicenses, license)
	}

	return filteredLicenses, nil
}

func (l *LicenseDecider) OnboardLicense(ctx contractapi.TransactionContextInterface, license *model.License) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	// check user can onboard license
	if ok, err := l.decider.HasPermissions(l.user, rbacpolicy.LicensesOA, operations.OnboardLicense); err != nil {
		return errors.Wrapf(err, "error checking if user %s can onboard a license", l.user)
	} else if !ok {
		return ErrAccessDenied
	}

	return l.pap.OnboardLicense(ctx, license)
}

func (l *LicenseDecider) OffboardLicense(ctx contractapi.TransactionContextInterface, licenseID string) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	// check user can onboard license
	if ok, err := l.decider.HasPermissions(l.user, licenseID, operations.OffboardLicense); err != nil {
		return errors.Wrapf(err, "error checking if user %s can offboard a license", l.user)
	} else if !ok {
		return ErrAccessDenied
	}

	return l.pap.OffboardLicense(ctx, licenseID)
}

func (l *LicenseDecider) CheckoutLicense(ctx contractapi.TransactionContextInterface, agencyName string, licenseID string,
	keys map[string]time.Time) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	// check user can checkout license
	if ok, err := l.decider.HasPermissions(l.user, licenseID, operations.CheckOutLicense); err != nil {
		return errors.Wrapf(err, "error checking if user %s can checkout a license", l.user)
	} else if !ok {
		return ErrAccessDenied
	}

	return l.pap.CheckoutLicense(ctx, agencyName, licenseID, keys)
}

func (l *LicenseDecider) CheckinLicense(ctx contractapi.TransactionContextInterface, agencyName string, licenseID string,
	keys []string) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	// check user can checkin license
	if ok, err := l.decider.HasPermissions(l.user, licenseID, operations.CheckInLicense); err != nil {
		return errors.Wrapf(err, "error checking if user %s can checkin a license", l.user)
	} else if !ok {
		return ErrAccessDenied
	}

	return l.pap.CheckinLicense(ctx, agencyName, licenseID, keys)
}
