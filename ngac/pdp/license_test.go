package pdp

import (
	"encoding/json"
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

	graphBytes, err := initLicenseTestGraph(t, transactionContext, chaincodeStub)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(graphBytes, nil)

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

		chaincodeStub.GetStateReturns(graphBytes, nil)

		err = decider.setup(transactionContext)
		require.NoError(t, err)

		err = decider.OnboardLicense(transactionContext, license)
		require.NoError(t, err)
	})

	t.Run("test a1 system owner", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org2MSP", nil)
		clientIdentity.GetX509CertificateReturns(A1SystemOwnerCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		chaincodeStub.GetStateReturns(graphBytes, nil)

		err = decider.setup(transactionContext)
		require.NoError(t, err)

		err = decider.OnboardLicense(transactionContext, license)
		require.Error(t, err)
	})
}

func TestOffboardLicense(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	_, err := initLicenseTestGraph(t, transactionContext, chaincodeStub)
	require.NoError(t, err)

	decider := NewLicenseDecider()

	t.Run("test org1 admin", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org1MSP", nil)
		clientIdentity.GetX509CertificateReturns(Org1AdminCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		err = decider.OffboardLicense(transactionContext, "test-license-id")
		require.NoError(t, err)
	})

	t.Run("test a1 system owner", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org2MSP", nil)
		clientIdentity.GetX509CertificateReturns(A1SystemOwnerCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		err = decider.OffboardLicense(transactionContext, "test-license-id")
		require.Error(t, err)
	})
}

// initialize a test graph to be used by test methods
func initLicenseTestGraph(t *testing.T, ctx *mocks.TransactionContext, stub *mocks.ChaincodeStub) ([]byte, error) {
	// create a new ngac graph and configure it using the blossom policy
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	// set the ngac graph as the result of get state
	// later when OnboardLicense is called, this graph will be used to determine if
	// the org1 admin has permission to onboard a license
	graphBytes, err := json.Marshal(graph)
	require.NoError(t, err)
	stub.GetStateReturns(graphBytes, nil)

	// set up the mock identity as the org1 admin
	clientIdentity := &mocks.ClientIdentity{}
	clientIdentity.GetMSPIDReturns("Org1MSP", nil)
	clientIdentity.GetX509CertificateReturns(Org1AdminCert(), nil)
	ctx.GetClientIdentityReturns(clientIdentity)

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
	graphBytes, err = json.Marshal(licenseDecider.pap.Graph())
	require.NoError(t, err)
	stub.GetStateReturns(graphBytes, nil)
	return graphBytes, nil
}
