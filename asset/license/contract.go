package license

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

// Contract for managing licenses
type Contract struct {
	contractapi.Contract
}

// OnboardLicense adds a new license to Blossom.  This will create a new license object on the ledger and in the NGAC graph.
// Licenses are identified by the LicenseNumber field. The user performing the request will need to have permission
// to add a license to the ledger/NGAC.
func (c Contract) OnboardLicense(ctx contractapi.TransactionContextInterface, license License) error {
	return nil
}

// OffboardLicense removes an existing license in Blossom.  This will remove the license from the ledger and from NGAC.
// TODO this should revoke any outstanding license keys
func (c Contract) OffboardLicense(ctx contractapi.TransactionContextInterface, licenseID string) error {
	return nil
}

// Licenses returns all licenses in Blossom.
func (c Contract) Licenses(ctx contractapi.TransactionContextInterface) ([]License, error) {
	return nil, nil
}

// LicenseInfo returns the info for the license with the given license ID.
func (c Contract) LicenseInfo(ctx contractapi.TransactionContextInterface, licenseID string) (License, error) {
	return License{}, nil
}

// LicenseKeys retrieves all of the agencies and their assigned keys for the given license. The returned map will store
// the agency names as the keys and an array of license keys as the values.
func (c Contract) LicenseKeys(ctx contractapi.TransactionContextInterface, licenseID string) (map[string][]string, error) {
	return nil, nil
}

// AgencyLicenseKeys retrieves all of the licenses the agency has checkout and the associated keys.  The returned map
// will store the license IDs as the keys and an array of license keys as the values.
func (c Contract) AgencyLicenseKeys(ctx contractapi.TransactionContextInterface, agency string) (map[int][]string, error) {
	return nil, nil
}

// CheckoutLicense requests a software license for an agency.  The requesting user must have permission to request
// (i.e. System Administrator). The amount parameter is the amount of software license keys the agency is requesting.
// This number is subtracted from the total available for the license. Return the set of keys that are now assigned to
// the agency.
func (c Contract) CheckoutLicense(ctx contractapi.TransactionContextInterface, licenseID string, agency string, amount int) ([]string, error) {
	return nil, nil
}

// CheckinLicense returns the license keys to Blossom.  The return of these keys is reflected in the amount available for
// the license, and the keys assigned to the agency on the ledger.
func (c Contract) CheckinLicense(ctx contractapi.TransactionContextInterface, licenseID string) error {
	return nil
}
