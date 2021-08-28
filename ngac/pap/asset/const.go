package asset

import "fmt"

func ObjectAttribute(licenseID string) string {
	return fmt.Sprintf("%s", licenseID)
}

func AssetsObjectAttribute(account string) string {
	return fmt.Sprintf("%s assets", account)
}

func LicenseObject(asset string, license string) string {
	return fmt.Sprintf("%s:%s", asset, license)
}
