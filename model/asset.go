package model

import (
	"fmt"
	"time"
)

type (
	// Asset represents a software asset on the ledger.
	Asset struct {
		// ID is the unique identifier of the asset
		ID string `json:"id"`
		// Name is the common name of the asset
		Name string `json:"name"`
		// TotalAmount is the total number of licenses available to Blossom
		TotalAmount int `json:"total_amount"`
		// Available is the number of licenses that are currently available to be checked out
		Available int `json:"available"`
		// Cost is the cost of obtaining a license
		Cost float64 `json:"cost"`
		// OnboardingDate is the date in which the asset was added to Blossom
		OnboardingDate time.Time `json:"onboarding_date"`
		// Expiration is the date in which the asset will expire from Blossom
		Expiration time.Time `json:"expiration"`
		// Licenses is the complete set of licenses associated with this asset
		Licenses []string `json:"licenses"`
		// AvailableLicenses is the set of licenses that are available to be checked out
		AvailableLicenses []string `json:"available_licenses"`
		// CheckedOut stores the agencies that have checked out this asset and which licenses they have leased
		CheckedOut map[string]map[string]time.Time `json:"checked_out"`
	}
)

const AssetPrefix = "asset:"

// AssetKey returns the key for an asset on the ledger.  Assets are stored with the format: "asset:<asset_id>".
func AssetKey(id string) string {
	return fmt.Sprintf("%s%s", AssetPrefix, id)
}
