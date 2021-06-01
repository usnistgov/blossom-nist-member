package asset

import "fmt"

func ObjectAttribute(licenseID string) string {
	return fmt.Sprintf("%s", licenseID)
}

func AssetsObjectAttribute(agency string) string {
	return fmt.Sprintf("%s assets", agency)
}

func LicenseObject(asset string, license string) string {
	return fmt.Sprintf("%s:%s", asset, license)
}
