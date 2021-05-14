package pap

import (
	"encoding/json"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy/dac"
	"testing"
	"time"

	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/api/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
	licensepap "github.com/usnistgov/blossom/chaincode/ngac/pap/license"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy/rbac"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy/status"
)

func TestOnboardLicense(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	graphBytes, err := json.Marshal(graph)
	require.NoError(t, err)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.GetStateReturns(graphBytes, nil)

	licenseAdmin, err := NewLicenseAdmin(transactionContext)
	require.NoError(t, err)

	license := &model.License{
		ID:             "test-license-id",
		Name:           "test-license",
		TotalAmount:    5,
		Available:      5,
		Cost:           20,
		OnboardingDate: time.Date(2021, 5, 12, 12, 0, 0, 0, time.Local),
		Expiration:     time.Date(2026, 5, 12, 12, 0, 0, 0, time.Local),
		AllKeys:        []string{"1", "2", "3", "4", "5"},
		AvailableKeys:  []string{"1", "2", "3", "4", "5"},
		CheckedOut:     make(map[string]map[string]time.Time),
	}

	err = licenseAdmin.OnboardLicense(transactionContext, license)
	require.NoError(t, err)

	graph = licenseAdmin.Graph()
	ok, err := graph.Exists(licensepap.LicenseObjectAttribute(license.ID))
	require.NoError(t, err)
	require.True(t, ok)

	parents, err := graph.GetParents(licensepap.LicenseObjectAttribute(license.ID))
	require.NoError(t, err)
	require.Contains(t, parents, rbac.LicensesOA)
	require.Contains(t, parents, dac.LicensesOA)
	require.Contains(t, parents, status.LicensesOA)

	children, err := graph.GetChildren(licensepap.LicenseObjectAttribute(license.ID))
	require.NoError(t, err)
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "1"))
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "2"))
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "3"))
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "4"))
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "5"))
}

func TestOffboardLicense(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	graphBytes, err := json.Marshal(graph)
	require.NoError(t, err)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.GetStateReturns(graphBytes, nil)

	licenseAdmin, err := NewLicenseAdmin(transactionContext)
	require.NoError(t, err)

	license := &model.License{
		ID:             "test-license-id",
		Name:           "test-license",
		TotalAmount:    5,
		Available:      5,
		Cost:           20,
		OnboardingDate: time.Date(2021, 5, 12, 12, 0, 0, 0, time.Local),
		Expiration:     time.Date(2026, 5, 12, 12, 0, 0, 0, time.Local),
		AllKeys:        []string{"1", "2", "3", "4", "5"},
		AvailableKeys:  []string{"1", "2", "3", "4", "5"},
		CheckedOut:     make(map[string]map[string]time.Time),
	}

	err = licenseAdmin.OnboardLicense(transactionContext, license)
	require.NoError(t, err)

	err = licenseAdmin.OffboardLicense(transactionContext, license.ID)
	require.NoError(t, err)

	graph = licenseAdmin.Graph()
	ok, err := graph.Exists(licensepap.LicenseObjectAttribute(license.ID))
	require.NoError(t, err)
	require.False(t, ok)
	ok, err = graph.Exists(licensepap.LicenseKeyObject(license.ID, "1"))
	require.NoError(t, err)
	require.False(t, ok)
	ok, err = graph.Exists(licensepap.LicenseKeyObject(license.ID, "2"))
	require.NoError(t, err)
	require.False(t, ok)
	ok, err = graph.Exists(licensepap.LicenseKeyObject(license.ID, "3"))
	require.NoError(t, err)
	require.False(t, ok)
	ok, err = graph.Exists(licensepap.LicenseKeyObject(license.ID, "4"))
	require.NoError(t, err)
	require.False(t, ok)
	ok, err = graph.Exists(licensepap.LicenseKeyObject(license.ID, "5"))
	require.NoError(t, err)
	require.False(t, ok)
}

func TestCheckoutCheckinLicense(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	graphBytes, err := json.Marshal(graph)
	require.NoError(t, err)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.GetStateReturns(graphBytes, nil)

	licenseAdmin, err := NewLicenseAdmin(transactionContext)
	require.NoError(t, err)

	license := &model.License{
		ID:             "test-license-id",
		Name:           "test-license",
		TotalAmount:    5,
		Available:      5,
		Cost:           20,
		OnboardingDate: time.Date(2021, 5, 12, 12, 0, 0, 0, time.Local),
		Expiration:     time.Date(2026, 5, 12, 12, 0, 0, 0, time.Local),
		AllKeys:        []string{"1", "2", "3", "4", "5"},
		AvailableKeys:  []string{"1", "2", "3", "4", "5"},
		CheckedOut:     make(map[string]map[string]time.Time),
	}

	err = licenseAdmin.OnboardLicense(transactionContext, license)
	require.NoError(t, err)

	graphBytes, err = json.Marshal(licenseAdmin.graph)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(graphBytes, nil)

	// create a new test agency
	agencyAdmin, err := NewAgencyAdmin(transactionContext)
	require.NoError(t, err)

	agency := model.Agency{
		Name:  "Org2",
		ATO:   "",
		MSPID: "Org2MSP",
		Users: model.Users{
			SystemOwner:           "a1_system_owner",
			SystemAdministrator:   "a1_system_admin",
			AcquisitionSpecialist: "a1_acq_spec",
		},
		Status:   "",
		Licenses: nil,
	}

	err = agencyAdmin.RequestAccount(transactionContext, agency)
	require.NoError(t, err)

	restartBytes, err := json.Marshal(agencyAdmin.graph)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(restartBytes, nil)

	licenseAdmin, err = NewLicenseAdmin(transactionContext)
	require.NoError(t, err)
	err = licenseAdmin.CheckoutLicense(transactionContext, agency.Name, license.ID,
		map[string]time.Time{"1": {}, "2": {}, "3": {}})
	require.NoError(t, err)

	graph = licenseAdmin.graph
	children, err := graph.GetChildren(agencypap.ObjectAttributeName("Org2"))
	require.NoError(t, err)
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "1"))
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "2"))
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "3"))

	graphBytes, err = json.Marshal(licenseAdmin.graph)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(graphBytes, nil)

	licenseAdmin, err = NewLicenseAdmin(transactionContext)
	require.NoError(t, err)
	err = licenseAdmin.CheckinLicense(transactionContext, agency.Name, license.ID, []string{"1", "2", "3"})
	require.NoError(t, err)

	graph = licenseAdmin.graph
	children, err = graph.GetChildren(agencypap.ObjectAttributeName("Org2"))
	require.NoError(t, err)
	require.NotContains(t, children, licensepap.LicenseKeyObject(license.ID, "1"))
	require.NotContains(t, children, licensepap.LicenseKeyObject(license.ID, "2"))
	require.NotContains(t, children, licensepap.LicenseKeyObject(license.ID, "3"))

	// test only returning 2 of 3 keys
	chaincodeStub.GetStateReturns(restartBytes, nil)

	licenseAdmin, err = NewLicenseAdmin(transactionContext)
	require.NoError(t, err)
	err = licenseAdmin.CheckoutLicense(transactionContext, agency.Name, license.ID,
		map[string]time.Time{"1": {}, "2": {}, "3": {}})
	require.NoError(t, err)

	graph = licenseAdmin.graph
	children, err = graph.GetChildren(agencypap.ObjectAttributeName("Org2"))
	require.NoError(t, err)
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "1"))
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "2"))
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "3"))

	graphBytes, err = json.Marshal(licenseAdmin.graph)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(graphBytes, nil)

	licenseAdmin, err = NewLicenseAdmin(transactionContext)
	require.NoError(t, err)
	err = licenseAdmin.CheckinLicense(transactionContext, agency.Name, license.ID, []string{"1", "2"})
	require.NoError(t, err)

	graph = licenseAdmin.graph
	children, err = graph.GetChildren(agencypap.ObjectAttributeName("Org2"))
	require.NoError(t, err)
	require.NotContains(t, children, licensepap.LicenseKeyObject(license.ID, "1"))
	require.NotContains(t, children, licensepap.LicenseKeyObject(license.ID, "2"))
	require.Contains(t, children, licensepap.LicenseKeyObject(license.ID, "3"))
}
