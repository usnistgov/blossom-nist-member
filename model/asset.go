package model

import (
	"fmt"
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
		OnboardingDate DateTime `json:"onboarding_date"`
		// Expiration is the date in which the asset will expire from Blossom
		Expiration DateTime `json:"expiration"`
		// Licenses is the complete set of licenses associated with this asset
		Licenses []string `json:"licenses"`
		// AvailableLicenses is the set of licenses that are available to be checked out
		AvailableLicenses []string `json:"available_licenses"`
		// CheckedOut stores the accounts that have checked out this asset, which licenses they have leased and the
		// expiration for each license
		CheckedOut map[string]map[string]DateTime `json:"checked_out"`
	}

	DateTime string
)

const AssetPrefix = "asset:"

// AssetKey returns the key for an asset on the ledger.  Assets are stored with the format: "asset:<asset_id>".
func AssetKey(id string) string {
	return fmt.Sprintf("%s%s", AssetPrefix, id)
}
