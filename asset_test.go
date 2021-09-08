package main

import (
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
)

func TestCheckoutLicense(t *testing.T) {
	asset := &model.Asset{
		ID:                "123",
		Name:              "my-asset",
		TotalAmount:       3,
		Available:         3,
		Cost:              20,
		OnboardingDate:    "",
		Expiration:        "",
		Licenses:          []string{"1", "2", "3"},
		AvailableLicenses: []string{"1", "2", "3"},
		CheckedOut:        make(map[string]map[string]model.DateTime),
	}

	account := &model.Account{
		Name:   "Account1",
		ATO:    "",
		MSPID:  "Account1MSP",
		Users:  model.Users{},
		Status: "",
		Assets: make(map[string]map[string]model.DateTime),
	}

	licenses, err := checkout(account, asset, 2)
	require.NoError(t, err)

	require.Contains(t, licenses, "1")
	require.Contains(t, licenses, "2")

	require.Equal(t, []string{"3"}, asset.AvailableLicenses)
	require.Equal(t, 1, asset.Available)
	require.Contains(t, asset.CheckedOut, "Account1")
	require.Contains(t, asset.CheckedOut["Account1"], "1")
	require.Contains(t, asset.CheckedOut["Account1"], "2")

	require.Contains(t, account.Assets, "123")
	require.Contains(t, account.Assets["123"], "1")
	require.Contains(t, account.Assets["123"], "2")
}

func TestCheckInLicense(t *testing.T) {
	asset := &model.Asset{
		ID:                "123",
		Name:              "my-asset",
		TotalAmount:       3,
		Available:         3,
		Cost:              20,
		OnboardingDate:    "",
		Expiration:        "",
		Licenses:          []string{"1", "2", "3"},
		AvailableLicenses: []string{"1", "2", "3"},
		CheckedOut:        make(map[string]map[string]model.DateTime),
	}

	account := &model.Account{
		Name:   "Account1",
		ATO:    "",
		MSPID:  "Account1MSP",
		Users:  model.Users{},
		Status: "",
		Assets: make(map[string]map[string]model.DateTime),
	}

	t.Run("test return all licenses", func(t *testing.T) {
		_, err := checkout(account, asset, 2)
		require.NoError(t, err)

		err = checkin(account, asset, []string{"1", "2"})
		require.NoError(t, err)

		require.Equal(t, []string{"3", "1", "2"}, asset.AvailableLicenses)
		require.Equal(t, 3, asset.Available)
		require.NotContains(t, asset.CheckedOut, "Account1")
		require.NotContains(t, account.Assets, "123")
	})

	t.Run("test return 2 of 3 licenses", func(t *testing.T) {
		_, err := checkout(account, asset, 3)
		require.NoError(t, err)

		err = checkin(account, asset, []string{"1", "2"})
		require.NoError(t, err)

		require.Equal(t, []string{"1", "2"}, asset.AvailableLicenses)
		require.Equal(t, 2, asset.Available)
		require.Contains(t, asset.CheckedOut, "Account1")
		require.Contains(t, asset.CheckedOut["Account1"], "3")
		require.Contains(t, account.Assets, "123")
		require.Contains(t, account.Assets["123"], "3")
	})

}
