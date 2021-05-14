package pdp

import (
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/api/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy"
	"testing"
	"time"
)

func TestOnboardLicense(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	err := initLicenseTestGraph(t, transactionContext, chaincodeStub)
	require.NoError(t, err)

	license := &model.License{
		ID:             "",
		Name:           "",
		TotalAmount:    0,
		Available:      0,
		Cost:           0,
		OnboardingDate: time.Time{},
		Expiration:     time.Time{},
		AllKeys:        make([]string, 0),
		AvailableKeys:  make([]string, 0),
		CheckedOut:     make(map[string]map[string]time.Time),
	}

	decider := NewLicenseDecider()

	t.Run("test org1 admin", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org1MSP", nil)
		clientIdentity.GetX509CertificateReturns(Org1AdminCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		err = decider.OnboardLicense(transactionContext, license)
		require.NoError(t, err)
	})

	t.Run("test a1 system owner", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org2MSP", nil)
		clientIdentity.GetX509CertificateReturns(A1SystemOwnerCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		err = decider.OnboardLicense(transactionContext, license)
		require.Error(t, err)
	})
}

func TestOffboardLicense(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	err := initLicenseTestGraph(t, transactionContext, chaincodeStub)
	require.NoError(t, err)

	decider := NewLicenseDecider()
	require.NoError(t, err)

	t.Run("test org1 admin", func(t *testing.T) {
		SetUser(transactionContext, Org1AdminCert(), "Org1MSP")

		err = decider.OffboardLicense(transactionContext, "test-license-id")
		require.NoError(t, err)
	})

	t.Run("test a1 system owner", func(t *testing.T) {
		SetUser(transactionContext, A1SystemOwnerCert(), "Org2MSP")

		err = decider.OffboardLicense(transactionContext, "test-license-id")
		require.Error(t, err)
	})
}

func TestCheckoutLicense(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	t.Run("test checkout active", func(t *testing.T) {
		// initialize the test graph with an onboarded license
		err := initLicenseTestGraph(t, transactionContext, chaincodeStub)
		require.NoError(t, err)

		// request account
		agency := model.Agency{
			Name:  "Org2",
			ATO:   "ato",
			MSPID: "Org2MSP",
			Users: model.Users{
				SystemOwner:           "a1_system_owner",
				SystemAdministrator:   "a1_system_admin",
				AcquisitionSpecialist: "a1_acq_spec",
			},
			Status:   "status",
			Licenses: make(map[string]map[string]time.Time),
		}

		SetUser(transactionContext, A1SystemOwnerCert(), "Org2MSP")

		agencyDecider := NewAgencyDecider()
		err = agencyDecider.RequestAccount(transactionContext, agency)
		require.NoError(t, err)

		SetGraphState(t, chaincodeStub, agencyDecider.pap.Graph())

		// approve agency
		SetUser(transactionContext, Org1AdminCert(), "Org1MSP")
		agencyDecider = NewAgencyDecider()
		err = agencyDecider.UpdateAgencyStatus(transactionContext, agency.Name, model.Approved)
		require.NoError(t, err)

		SetGraphState(t, chaincodeStub, agencyDecider.pap.Graph())

		SetUser(transactionContext, A1SystemAdminCert(), "Org2MSP")
		licenseDecider := NewLicenseDecider()
		err = licenseDecider.CheckoutLicense(transactionContext, "Org2", "test-license-id",
			map[string]time.Time{"1": time.Now()})
		require.NoError(t, err)
	})

	t.Run("test checkout inactive", func(t *testing.T) {
		// initialize the test graph with an onboarded license
		err := initLicenseTestGraph(t, transactionContext, chaincodeStub)
		require.NoError(t, err)

		// request account
		agency := model.Agency{
			Name:  "Org2",
			ATO:   "ato",
			MSPID: "Org2MSP",
			Users: model.Users{
				SystemOwner:           "a1_system_owner",
				SystemAdministrator:   "a1_system_admin",
				AcquisitionSpecialist: "a1_acq_spec",
			},
			Status:   "status",
			Licenses: make(map[string]map[string]time.Time),
		}

		SetUser(transactionContext, A1SystemOwnerCert(), "Org2MSP")

		agencyDecider := NewAgencyDecider()
		err = agencyDecider.RequestAccount(transactionContext, agency)
		require.NoError(t, err)

		SetGraphState(t, chaincodeStub, agencyDecider.pap.Graph())

		// do not approve agency

		// checkout license as pending
		SetUser(transactionContext, A1SystemAdminCert(), "Org2MSP")
		licenseDecider := NewLicenseDecider()
		err = licenseDecider.CheckoutLicense(transactionContext, "Org2", "test-license-id",
			map[string]time.Time{"1": time.Now()})
		require.Error(t, err)
	})
}

// initialize a test graph to be used by test methods
func initLicenseTestGraph(t *testing.T, ctx *mocks.TransactionContext, stub *mocks.ChaincodeStub) error {
	// create a new ngac graph and configure it using the blossom policy
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	// set the ngac graph as the result of get state
	// later when OnboardLicense is called, this graph will be used to determine if
	// the org1 admin has permission to onboard a license
	SetGraphState(t, stub, graph)

	// set up the mock identity as the org1 admin
	SetUser(ctx, Org1AdminCert(), "Org1MSP")

	// create a test license
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

	// onboard the license as the org1 admin
	licenseDecider := NewLicenseDecider()
	err = licenseDecider.OnboardLicense(ctx, license)
	require.NoError(t, err)

	// re marshal the graph bytes using the PAP's graph
	// this will have the graph that includes the onboarded license
	// this is the graph the tests will operate on
	SetGraphState(t, stub, licenseDecider.pap.Graph())
	return nil
}
