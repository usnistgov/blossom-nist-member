package pdp

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/api/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	assetpap "github.com/usnistgov/blossom/chaincode/ngac/pap/asset"
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

	asset := &model.Asset{
		ID:                "",
		Name:              "",
		TotalAmount:       0,
		Available:         0,
		Cost:              0,
		OnboardingDate:    time.Time{},
		Expiration:        time.Time{},
		Licenses:          make([]string, 0),
		AvailableLicenses: make([]string, 0),
		CheckedOut:        make(map[string]map[string]time.Time),
	}

	decider := NewAssetDecider()

	t.Run("test org1 admin", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org1MSP", nil)
		clientIdentity.GetX509CertificateReturns(Org1AdminCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		err = decider.OnboardAsset(transactionContext, asset)
		require.NoError(t, err)
	})

	t.Run("test a1 system owner", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org2MSP", nil)
		clientIdentity.GetX509CertificateReturns(A1SystemOwnerCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		err = decider.OnboardAsset(transactionContext, asset)
		require.Error(t, err)
	})
}

func TestOffboardLicense(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	err := initLicenseTestGraph(t, transactionContext, chaincodeStub)
	require.NoError(t, err)

	decider := NewAssetDecider()
	require.NoError(t, err)

	t.Run("test org1 admin", func(t *testing.T) {
		SetUser(transactionContext, Org1AdminCert(), "Org1MSP")

		err = decider.OffboardAsset(transactionContext, "test-license-id")
		require.NoError(t, err)
	})

	t.Run("test a1 system owner", func(t *testing.T) {
		SetUser(transactionContext, A1SystemOwnerCert(), "Org2MSP")

		err = decider.OffboardAsset(transactionContext, "test-license-id")
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
			Status: "status",
			Assets: make(map[string]map[string]time.Time),
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
		licenseDecider := NewAssetDecider()
		err = licenseDecider.Checkout(transactionContext, "Org2", "test-license-id",
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
			Status: "status",
			Assets: make(map[string]map[string]time.Time),
		}

		SetUser(transactionContext, A1SystemOwnerCert(), "Org2MSP")

		agencyDecider := NewAgencyDecider()
		err = agencyDecider.RequestAccount(transactionContext, agency)
		require.NoError(t, err)

		SetGraphState(t, chaincodeStub, agencyDecider.pap.Graph())

		// do not approve agency

		// checkout license as pending
		SetUser(transactionContext, A1SystemAdminCert(), "Org2MSP")
		licenseDecider := NewAssetDecider()
		err = licenseDecider.Checkout(transactionContext, "Org2", "test-license-id",
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
	// later when OnboardAsset is called, this graph will be used to determine if
	// the org1 admin has permission to onboard a license
	SetGraphState(t, stub, graph)

	// set up the mock identity as the org1 admin
	SetUser(ctx, Org1AdminCert(), "Org1MSP")

	// create a test license
	asset := &model.Asset{
		ID:                "test-license-id",
		Name:              "test-license",
		TotalAmount:       5,
		Available:         5,
		Cost:              20,
		OnboardingDate:    time.Date(2021, 5, 12, 12, 0, 0, 0, time.Local),
		Expiration:        time.Date(2026, 5, 12, 12, 0, 0, 0, time.Local),
		Licenses:          []string{"1", "2", "3", "4", "5"},
		AvailableLicenses: []string{"1", "2", "3", "4", "5"},
		CheckedOut:        make(map[string]map[string]time.Time),
	}

	// onboard the license as the org1 admin
	licenseDecider := NewAssetDecider()
	err = licenseDecider.OnboardAsset(ctx, asset)
	require.NoError(t, err)

	// re marshal the graph bytes using the PAP's graph
	// this will have the graph that includes the onboarded license
	// this is the graph the tests will operate on
	SetGraphState(t, stub, licenseDecider.pap.Graph())
	return nil
}

func TestFilterLicense(t *testing.T) {
	graph := memory.NewGraph()
	pcNode, err := graph.CreateNode("pc1", pip.PolicyClass, nil)
	require.NoError(t, err)
	oa1, err := graph.CreateNode("oa1", pip.ObjectAttribute, nil)
	require.NoError(t, err)
	oa2, err := graph.CreateNode("oa2", pip.ObjectAttribute, nil)
	require.NoError(t, err)
	l1, err := graph.CreateNode(assetpap.ObjectAttribute("test-license-1"), pip.ObjectAttribute, nil)
	require.NoError(t, err)
	l2, err := graph.CreateNode(assetpap.ObjectAttribute("test-license-2"), pip.ObjectAttribute, nil)
	require.NoError(t, err)
	err = graph.Assign(oa1.Name, pcNode.Name)
	require.NoError(t, err)
	err = graph.Assign(oa2.Name, pcNode.Name)
	require.NoError(t, err)
	err = graph.Assign(l1.Name, oa1.Name)
	require.NoError(t, err)
	err = graph.Assign(l2.Name, oa2.Name)
	require.NoError(t, err)

	ua1, err := graph.CreateNode("ua1", pip.UserAttribute, nil)
	require.NoError(t, err)
	u1, err := graph.CreateNode("Org1 Admin:Org1MSP", pip.User, nil)
	require.NoError(t, err)
	err = graph.Assign(u1.Name, ua1.Name)
	require.NoError(t, err)
	err = graph.Assign(ua1.Name, pcNode.Name)
	require.NoError(t, err)

	err = graph.Associate("ua1", "oa1", pip.ToOps(operations.ViewAsset))
	require.NoError(t, err)
	err = graph.Associate("ua1", "oa2", pip.ToOps(operations.ViewAsset, operations.ViewAllLicenses,
		operations.ViewAvailableLicenses, operations.ViewCheckedOut))
	require.NoError(t, err)

	testTime := time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
	assets := []*model.Asset{
		{
			ID:                "test-license-1",
			Name:              "test-license-1",
			TotalAmount:       5,
			Available:         4,
			Cost:              99,
			OnboardingDate:    testTime,
			Expiration:        testTime,
			Licenses:          []string{"1", "2", "3", "4", "5"},
			AvailableLicenses: []string{"2", "3", "4", "5"},
			CheckedOut: map[string]map[string]time.Time{
				"agency1": {"test-license-1": testTime},
			},
		},
		{
			ID:                "test-license-2",
			Name:              "test-license-2",
			TotalAmount:       5,
			Available:         4,
			Cost:              99,
			OnboardingDate:    time.Time{},
			Expiration:        time.Time{},
			Licenses:          []string{"1", "2", "3", "4", "5"},
			AvailableLicenses: []string{"2", "3", "4", "5"},
			CheckedOut: map[string]map[string]time.Time{
				"agency1": {"test-license-2": testTime},
			},
		},
	}

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	SetGraphState(t, chaincodeStub, graph)
	SetUser(transactionContext, Org1AdminCert(), "Org1MSP")

	assets, err = NewAssetDecider().FilterAssets(transactionContext, assets)
	require.NoError(t, err)
	require.Equal(t, 2, len(assets))

	license := assets[0]
	require.Equal(t, "test-license-1", license.ID)
	require.Equal(t, "test-license-1", license.Name)
	require.Equal(t, 5, license.TotalAmount)
	require.Equal(t, 4, license.Available)
	require.Equal(t, float64(99), license.Cost)
	require.Equal(t, []string{}, license.Licenses)
	require.Equal(t, []string{}, license.AvailableLicenses)
	require.Equal(t, map[string]map[string]time.Time{}, license.CheckedOut)

	license = assets[1]
	require.Equal(t, "test-license-2", license.ID)
	require.Equal(t, "test-license-2", license.Name)
	require.Equal(t, 5, license.TotalAmount)
	require.Equal(t, 4, license.Available)
	require.Equal(t, float64(99), license.Cost)
	require.Equal(t, []string{"1", "2", "3", "4", "5"}, license.Licenses)
	require.Equal(t, []string{"2", "3", "4", "5"}, license.AvailableLicenses)
	require.Equal(t, map[string]map[string]time.Time{
		"agency1": {"test-license-2": testTime},
	}, license.CheckedOut)
}
