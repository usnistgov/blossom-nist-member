package api

import (
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
	"time"
)

func TestCheckoutLicense(t *testing.T) {
	asset := &model.Asset{
		ID:                "123",
		Name:              "my-asset",
		TotalAmount:       3,
		Available:         3,
		Cost:              20,
		OnboardingDate:    time.Time{},
		Expiration:        time.Time{},
		Licenses:          []string{"1", "2", "3"},
		AvailableLicenses: []string{"1", "2", "3"},
		CheckedOut:        make(map[string]map[string]time.Time),
	}

	agency := &model.Agency{
		Name:   "Agency1",
		ATO:    "",
		MSPID:  "Agency1MSP",
		Users:  model.Users{},
		Status: "",
		Assets: make(map[string]map[string]time.Time),
	}

	licenses, err := checkout(agency, asset, 2)
	require.NoError(t, err)

	require.Contains(t, licenses, "1")
	require.Contains(t, licenses, "2")

	require.Equal(t, []string{"3"}, asset.AvailableLicenses)
	require.Equal(t, 1, asset.Available)
	require.Contains(t, asset.CheckedOut, "Agency1")
	require.Contains(t, asset.CheckedOut["Agency1"], "1")
	require.Contains(t, asset.CheckedOut["Agency1"], "2")

	require.Contains(t, agency.Assets, "123")
	require.Contains(t, agency.Assets["123"], "1")
	require.Contains(t, agency.Assets["123"], "2")
}

func TestCheckInLicense(t *testing.T) {
	asset := &model.Asset{
		ID:                "123",
		Name:              "my-asset",
		TotalAmount:       3,
		Available:         3,
		Cost:              20,
		OnboardingDate:    time.Time{},
		Expiration:        time.Time{},
		Licenses:          []string{"1", "2", "3"},
		AvailableLicenses: []string{"1", "2", "3"},
		CheckedOut:        make(map[string]map[string]time.Time),
	}

	agency := &model.Agency{
		Name:   "Agency1",
		ATO:    "",
		MSPID:  "Agency1MSP",
		Users:  model.Users{},
		Status: "",
		Assets: make(map[string]map[string]time.Time),
	}

	t.Run("test return all licenses", func(t *testing.T) {
		_, err := checkout(agency, asset, 2)
		require.NoError(t, err)

		err = checkin(agency, asset, []string{"1", "2"})
		require.NoError(t, err)

		require.Equal(t, []string{"3", "1", "2"}, asset.AvailableLicenses)
		require.Equal(t, 3, asset.Available)
		require.NotContains(t, asset.CheckedOut, "Agency1")
		require.NotContains(t, agency.Assets, "123")
	})

	t.Run("test return 2 of 3 licenses", func(t *testing.T) {
		_, err := checkout(agency, asset, 3)
		require.NoError(t, err)

		err = checkin(agency, asset, []string{"1", "2"})
		require.NoError(t, err)

		require.Equal(t, []string{"1", "2"}, asset.AvailableLicenses)
		require.Equal(t, 2, asset.Available)
		require.Contains(t, asset.CheckedOut, "Agency1")
		require.Contains(t, asset.CheckedOut["Agency1"], "3")
		require.Contains(t, agency.Assets, "123")
		require.Contains(t, agency.Assets["123"], "3")
	})

}
