package asset

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

type (
	// LicenseInterface provides the functions to interact with Licenses in fabric.
	LicenseInterface interface {
		// OnboardLicense adds a new license to Blossom.  This will create a new license object on the ledger and in the
		// NGAC graph. Licenses are identified by the LicenseNumber field. The user performing the request will need to
		// have permission to add a license to the ledger/NGAC. The license will be an object attribute in NGAC and the
		// license keys will be objects that are assigned to the license.
		OnboardLicense(ctx contractapi.TransactionContextInterface, license *License) error

		// OffboardLicense removes an existing license in Blossom.  This will remove the license from the ledger
		// and from NGAC.
		// TODO this should revoke any outstanding license keys
		OffboardLicense(ctx contractapi.TransactionContextInterface, licenseID string) error

		// Licenses returns all licenses in Blossom. This information does not include which agencies have keys for each
		// license
		Licenses(ctx contractapi.TransactionContextInterface) ([]License, error)

		// LicenseInfo returns the info for the license with the given license ID.
		LicenseInfo(ctx contractapi.TransactionContextInterface, licenseID string) (License, error)

		// LicenseKeys retrieves all of the agencies and their assigned keys for the given license. The returned map will
		// store the agency names as the keys and an array of license keys as the values.
		LicenseKeys(ctx contractapi.TransactionContextInterface, licenseID string) (map[string][]string, error)

		// AgencyLicenseKeys retrieves all of the licenses the agency has checkout and the associated keys.  The
		// returned map will store the license IDs as the keys and an array of license keys as the values.
		AgencyLicenseKeys(ctx contractapi.TransactionContextInterface, agency string) (map[int][]string, error)

		// CheckoutLicense requests a software license for an agency.  The requesting user must have permission to request
		// (i.e. System Administrator). The amount parameter is the amount of software license keys the agency is requesting.
		// This number is subtracted from the total available for the license. Return the set of keys that are now assigned to
		// the agency.
		CheckoutLicense(ctx contractapi.TransactionContextInterface, licenseID string, agency string, amount int) ([]string, error)

		// CheckinLicense returns the license keys to Blossom.  The return of these keys is reflected in the amount available for
		// the license, and the keys assigned to the agency on the ledger.
		CheckinLicense(ctx contractapi.TransactionContextInterface, licenseID string) error
	}
)

func NewLicenseContract() LicenseInterface {
	return &BlossomContract{}
}

func (b *BlossomContract) OnboardLicense(ctx contractapi.TransactionContextInterface, license *License) error {
	return nil
}

func (b *BlossomContract) OffboardLicense(ctx contractapi.TransactionContextInterface, licenseID string) error {
	return nil
}

func (b *BlossomContract) Licenses(ctx contractapi.TransactionContextInterface) ([]License, error) {
	return nil, nil
}

func (b *BlossomContract) LicenseInfo(ctx contractapi.TransactionContextInterface, licenseID string) (License, error) {
	return License{}, nil
}

func (b *BlossomContract) LicenseKeys(ctx contractapi.TransactionContextInterface, licenseID string) (map[string][]string, error) {
	return nil, nil
}

func (b *BlossomContract) AgencyLicenseKeys(ctx contractapi.TransactionContextInterface, agency string) (map[int][]string, error) {
	return nil, nil
}

func (b *BlossomContract) CheckoutLicense(ctx contractapi.TransactionContextInterface, licenseID string, agency string, amount int) ([]string, error) {
	return nil, nil
}

func (b *BlossomContract) CheckinLicense(ctx contractapi.TransactionContextInterface, licenseID string) error {
	return nil
}
