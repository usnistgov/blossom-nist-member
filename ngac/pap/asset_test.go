package pap

import (
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy/dac"
	"testing"
	"time"

	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/model"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
	assetpap "github.com/usnistgov/blossom/chaincode/ngac/pap/asset"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy/rbac"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy/status"
)

func TestOnboardLicense(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	mock := mocks.New()
	mock.SetGraphState(graph)

	assetAdmin, err := NewAssetAdmin(mock.Stub)
	require.NoError(t, err)

	asset := &model.Asset{
		ID:                "test-asset-id",
		Name:              "test-asset",
		TotalAmount:       5,
		Available:         5,
		Cost:              20,
		OnboardingDate:    time.Date(2021, 5, 12, 12, 0, 0, 0, time.Local),
		Expiration:        time.Date(2026, 5, 12, 12, 0, 0, 0, time.Local),
		Licenses:          []string{"1", "2", "3", "4", "5"},
		AvailableLicenses: []string{"1", "2", "3", "4", "5"},
		CheckedOut:        make(map[string]map[string]time.Time),
	}

	err = assetAdmin.OnboardAsset(mock.Stub, asset)
	require.NoError(t, err)

	graph = assetAdmin.Graph()
	ok, err := graph.Exists(assetpap.ObjectAttribute(asset.ID))
	require.NoError(t, err)
	require.True(t, ok)

	parents, err := graph.GetParents(assetpap.ObjectAttribute(asset.ID))
	require.NoError(t, err)
	require.Contains(t, parents, rbac.AssetsOA)
	require.Contains(t, parents, dac.AssetsOA)
	require.Contains(t, parents, status.AssetsOA)

	children, err := graph.GetChildren(assetpap.ObjectAttribute(asset.ID))
	require.NoError(t, err)
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "1"))
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "2"))
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "3"))
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "4"))
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "5"))
}

func TestOffboardLicense(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	mock := mocks.New()
	mock.SetGraphState(graph)

	assetAdmin, err := NewAssetAdmin(mock.Stub)
	require.NoError(t, err)

	asset := &model.Asset{
		ID:                "test-asset-id",
		Name:              "test-asset",
		TotalAmount:       5,
		Available:         5,
		Cost:              20,
		OnboardingDate:    time.Date(2021, 5, 12, 12, 0, 0, 0, time.Local),
		Expiration:        time.Date(2026, 5, 12, 12, 0, 0, 0, time.Local),
		Licenses:          []string{"1", "2", "3", "4", "5"},
		AvailableLicenses: []string{"1", "2", "3", "4", "5"},
		CheckedOut:        make(map[string]map[string]time.Time),
	}

	err = assetAdmin.OnboardAsset(mock.Stub, asset)
	require.NoError(t, err)

	err = assetAdmin.OffboardAsset(mock.Stub, asset.ID)
	require.NoError(t, err)

	graph = assetAdmin.Graph()
	ok, err := graph.Exists(assetpap.ObjectAttribute(asset.ID))
	require.NoError(t, err)
	require.False(t, ok)
	ok, err = graph.Exists(assetpap.LicenseObject(asset.ID, "1"))
	require.NoError(t, err)
	require.False(t, ok)
	ok, err = graph.Exists(assetpap.LicenseObject(asset.ID, "2"))
	require.NoError(t, err)
	require.False(t, ok)
	ok, err = graph.Exists(assetpap.LicenseObject(asset.ID, "3"))
	require.NoError(t, err)
	require.False(t, ok)
	ok, err = graph.Exists(assetpap.LicenseObject(asset.ID, "4"))
	require.NoError(t, err)
	require.False(t, ok)
	ok, err = graph.Exists(assetpap.LicenseObject(asset.ID, "5"))
	require.NoError(t, err)
	require.False(t, ok)
}

func TestCheckoutCheckinLicense(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	mock := mocks.New()
	mock.SetGraphState(graph)

	assetAdmin, err := NewAssetAdmin(mock.Stub)
	require.NoError(t, err)

	asset := &model.Asset{
		ID:                "test-asset-id",
		Name:              "test-asset",
		TotalAmount:       5,
		Available:         5,
		Cost:              20,
		OnboardingDate:    time.Date(2021, 5, 12, 12, 0, 0, 0, time.Local),
		Expiration:        time.Date(2026, 5, 12, 12, 0, 0, 0, time.Local),
		Licenses:          []string{"1", "2", "3", "4", "5"},
		AvailableLicenses: []string{"1", "2", "3", "4", "5"},
		CheckedOut:        make(map[string]map[string]time.Time),
	}

	err = assetAdmin.OnboardAsset(mock.Stub, asset)
	require.NoError(t, err)

	mock.SetGraphState(assetAdmin.graph)

	// create a new test agency
	agencyAdmin, err := NewAgencyAdmin(mock.Stub)
	require.NoError(t, err)

	agency := &model.Agency{
		Name:  "Org2",
		ATO:   "",
		MSPID: "Org2MSP",
		Users: model.Users{
			SystemOwner:           "a1_system_owner",
			SystemAdministrator:   "a1_system_admin",
			AcquisitionSpecialist: "a1_acq_spec",
		},
		Status: "",
		Assets: nil,
	}

	err = agencyAdmin.RequestAccount(mock.Stub, agency)
	require.NoError(t, err)

	restartGraph := agencyAdmin.graph
	mock.SetGraphState(agencyAdmin.graph)

	assetAdmin, err = NewAssetAdmin(mock.Stub)
	require.NoError(t, err)
	err = assetAdmin.Checkout(mock.Stub, agency.Name, asset.ID,
		map[string]time.Time{"1": {}, "2": {}, "3": {}})
	require.NoError(t, err)

	graph = assetAdmin.graph
	children, err := graph.GetChildren(agencypap.ObjectAttributeName("Org2"))
	require.NoError(t, err)
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "1"))
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "2"))
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "3"))

	mock.SetGraphState(assetAdmin.graph)

	assetAdmin, err = NewAssetAdmin(mock.Stub)
	require.NoError(t, err)
	err = assetAdmin.Checkin(mock.Stub, agency.Name, asset.ID, []string{"1", "2", "3"})
	require.NoError(t, err)

	graph = assetAdmin.graph
	children, err = graph.GetChildren(agencypap.ObjectAttributeName("Org2"))
	require.NoError(t, err)
	require.NotContains(t, children, assetpap.LicenseObject(asset.ID, "1"))
	require.NotContains(t, children, assetpap.LicenseObject(asset.ID, "2"))
	require.NotContains(t, children, assetpap.LicenseObject(asset.ID, "3"))

	// test only returning 2 of 3 keys
	mock.SetGraphState(restartGraph)

	assetAdmin, err = NewAssetAdmin(mock.Stub)
	require.NoError(t, err)
	err = assetAdmin.Checkout(mock.Stub, agency.Name, asset.ID,
		map[string]time.Time{"1": {}, "2": {}, "3": {}})
	require.NoError(t, err)

	graph = assetAdmin.graph
	children, err = graph.GetChildren(agencypap.ObjectAttributeName("Org2"))
	require.NoError(t, err)
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "1"))
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "2"))
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "3"))

	mock.SetGraphState(assetAdmin.graph)

	assetAdmin, err = NewAssetAdmin(mock.Stub)
	require.NoError(t, err)
	err = assetAdmin.Checkin(mock.Stub, agency.Name, asset.ID, []string{"1", "2"})
	require.NoError(t, err)

	graph = assetAdmin.graph
	children, err = graph.GetChildren(agencypap.ObjectAttributeName("Org2"))
	require.NoError(t, err)
	require.NotContains(t, children, assetpap.LicenseObject(asset.ID, "1"))
	require.NotContains(t, children, assetpap.LicenseObject(asset.ID, "2"))
	require.Contains(t, children, assetpap.LicenseObject(asset.ID, "3"))
}
