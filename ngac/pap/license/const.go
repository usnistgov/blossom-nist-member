package license

import "fmt"

func LicenseObjectAttribute(licenseID string) string {
	return fmt.Sprintf("%s", licenseID)
}

func LicensesObjectAttribute(agency string) string {
	return fmt.Sprintf("%s licenses", agency)
}

func LicenseKeyObject(licenseID string, key string) string {
	return fmt.Sprintf("%s:%s", licenseID, key)
}
