package model

import (
	"fmt"
	"time"
)

type (
	// License represents a license on the ledger.
	License struct {
		// ID is the unique identifier of the license
		ID string `json:"id"`
		// Name is the common name of the license
		Name string `json:"name"`
		// TotalAmount is the total number of keys available to Blossom
		TotalAmount int `json:"total_amount"`
		// Available is the number of license keys that are currently available to be checked out
		Available int `json:"available"`
		// Cost is the cost of obtaining the license
		Cost float64 `json:"cost"`
		// OnboardingDate is the date in which the license was added to Blossom
		OnboardingDate time.Time `json:"onboarding_date"`
		// Expiration is the date in which the license will expire from Blossom
		Expiration time.Time `json:"expiration"`
		// AllKeys is the complete set of license keys associated with this license
		AllKeys []string `json:"all_keys"`
		// AvailableKeys is the set of keys that are available to be checked out
		AvailableKeys []string `json:"available_keys"`
		// CheckedOut stores the agencies that have checked out this license and which license keys they have leased
		CheckedOut map[string]map[string]time.Time `json:"checked_out"`
	}
)

const (
	LicensePrefix = "license:"
)

// LicenseKey returns the key for a license on the ledger.  Licenses are stored with the format: "license:<license_id>".
func LicenseKey(name string) string {
	return fmt.Sprintf("%s%s", LicensePrefix, name)
}
