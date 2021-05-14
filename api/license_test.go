package api

import (
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
	"time"
)

var (
	license = &model.License{
		ID:             "123",
		Name:           "my-license",
		TotalAmount:    3,
		Available:      3,
		Cost:           20,
		OnboardingDate: time.Time{},
		Expiration:     time.Time{},
		AllKeys:        []string{"1", "2", "3"},
		AvailableKeys:  []string{"1", "2", "3"},
		CheckedOut:     make(map[string]map[string]time.Time),
	}

	agency = &model.Agency{
		Name:     "Agency1",
		ATO:      "",
		MSPID:    "Agency1MSP",
		Users:    model.Users{},
		Status:   "",
		Licenses: make(map[string]map[string]time.Time),
	}
)

func TestCheckoutLicense(t *testing.T) {
	keys, err := checkoutLicense(agency, license, 2)
	require.NoError(t, err)

	require.Contains(t, keys, "1")
	require.Contains(t, keys, "2")

	require.Equal(t, []string{"3"}, license.AvailableKeys)
	require.Equal(t, 1, license.Available)
	require.Contains(t, license.CheckedOut, "Agency1")
	require.Contains(t, license.CheckedOut["Agency1"], "1")
	require.Contains(t, license.CheckedOut["Agency1"], "2")

	require.Contains(t, agency.Licenses, "123")
	require.Contains(t, agency.Licenses["123"], "1")
	require.Contains(t, agency.Licenses["123"], "2")
}

func TestCheckInLicense(t *testing.T) {
	t.Run("test return all keys", func(t *testing.T) {
		_, err := checkoutLicense(agency, license, 2)
		require.NoError(t, err)

		err = checkinLicense(agency, license, []string{"1", "2"})
		require.NoError(t, err)

		require.Equal(t, []string{"3", "1", "2"}, license.AvailableKeys)
		require.Equal(t, 3, license.Available)
		require.NotContains(t, license.CheckedOut, "Agency1")
		require.NotContains(t, agency.Licenses, "123")
	})

	t.Run("test return 2 of 3 keys", func(t *testing.T) {
		_, err := checkoutLicense(agency, license, 3)
		require.NoError(t, err)

		err = checkinLicense(agency, license, []string{"1", "2"})
		require.NoError(t, err)

		require.Equal(t, []string{"1", "2"}, license.AvailableKeys)
		require.Equal(t, 2, license.Available)
		require.Contains(t, license.CheckedOut, "Agency1")
		require.Contains(t, license.CheckedOut["Agency1"], "3")
		require.Contains(t, agency.Licenses, "123")
		require.Contains(t, agency.Licenses["123"], "3")
	})

}
