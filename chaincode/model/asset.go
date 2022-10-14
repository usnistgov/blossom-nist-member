package model

import (
	"fmt"
)

type (
	AssetPrivate struct {
		// TotalAmount is the total number of licenses available to Blossom
		TotalAmount int `json:"total_amount"`
		// Licenses is the complete set of licenses associated with this asset
		Licenses map[string]string `json:"licenses"`
		// AvailableLicenses is the set of licenses that are available to be checked out
		AvailableLicenses []string `json:"available_licenses"`
		// CheckedOut stores the accounts that have checked out this asset, which licenses they have leased and the
		// expiration for each license
		CheckedOut map[string]map[string]string `json:"checked_out"`
	}

	// AssetPublic represents the public info for software asset on the ledger.
	AssetPublic struct {
		// ID is the unique identifier of the asset
		ID string `json:"id"`
		// Name is the common name of the asset
		Name string `json:"name"`
		// Available is the number of licenses that are currently available to be checked out
		Available int `json:"available"`
		// OnboardingDate is the date in which the asset was added to Blossom
		OnboardingDate string `json:"onboarding_date"`
		// Expiration is the date in which the asset will expire from Blossom
		Expiration string `json:"expiration"`
	}

	Asset struct {
		// ID is the unique identifier of the asset
		ID string `json:"id"`
		// Name is the common name of the asset
		Name string `json:"name"`
		// Available is the number of licenses that are currently available to be checked out
		Available int `json:"available"`
		// OnboardingDate is the date in which the asset was added to Blossom
		OnboardingDate string `json:"onboarding_date"`
		// Expiration is the date in which the asset will expire from Blossom
		Expiration string `json:"expiration"`
		// TotalAmount is the total number of licenses available to Blossom
		TotalAmount int `json:"total_amount"`
		// Licenses is the complete set of licenses associated with this asset
		Licenses map[string]string `json:"licenses"`
		// AvailableLicenses is the set of licenses that are available to be checked out
		AvailableLicenses []string `json:"available_licenses"`
		// CheckedOut stores the accounts that have checked out this asset, which licenses they have leased and the
		// expiration for each license
		CheckedOut map[string]map[string]string `json:"checked_out"`
	}

	License struct {
		LicenseID  string `json:"license_id,omitempty"`
		Expiration string `json:"expiration,omitempty"`
	}
)

const AssetPrefix = "asset:"

// AssetKey returns the key for an asset on the ledger.  Assets are stored with the format: "asset:<asset_id>".
func AssetKey(id string) string {
	return fmt.Sprintf("%s%s", AssetPrefix, id)
}

func NewAssetPublic() *AssetPublic {
	return &AssetPublic{
		ID:             "",
		Name:           "",
		Available:      0,
		OnboardingDate: "",
		Expiration:     "",
	}
}

func NewAssetPrivate() *AssetPrivate {
	return &AssetPrivate{
		TotalAmount:       0,
		Licenses:          make(map[string]string),
		AvailableLicenses: make([]string, 0),
		CheckedOut:        make(map[string]map[string]string),
	}
}

func NewAsset() *Asset {
	return &Asset{
		ID:                "",
		Name:              "",
		Available:         0,
		OnboardingDate:    "",
		Expiration:        "",
		TotalAmount:       0,
		Licenses:          make(map[string]string),
		AvailableLicenses: make([]string, 0),
		CheckedOut:        make(map[string]map[string]string),
	}
}
