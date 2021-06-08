package pap

import (
	"encoding/json"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	mocks2 "github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	assetpap "github.com/usnistgov/blossom/chaincode/ngac/pap/asset"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy/rbac"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy/status"
	"testing"
	"time"
)

func TestReportSwID(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	graphBytes, err := json.Marshal(graph)
	require.NoError(t, err)

	chaincodeStub := &mocks2.ChaincodeStub{}
	transactionContext := &mocks2.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.GetStateReturns(graphBytes, nil)

	assetAdmin, err := NewAssetAdmin(transactionContext)
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

	err = assetAdmin.OnboardAsset(transactionContext, asset)
	require.NoError(t, err)

	graphBytes, err = json.Marshal(assetAdmin.graph)
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
		Status: "",
		Assets: nil,
	}

	err = agencyAdmin.RequestAccount(transactionContext, agency)
	require.NoError(t, err)

	graphBytes, err = json.Marshal(agencyAdmin.graph)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(graphBytes, nil)

	swidAdmin, err := NewSwIDAdmin(transactionContext)
	require.NoError(t, err)

	swid := &model.SwID{
		PrimaryTag:      "pt1",
		XML:             "xml",
		Asset:           "test-asset-id",
		License:         "1",
		LeaseExpiration: time.Time{},
	}

	err = swidAdmin.ReportSwID(transactionContext, swid, "Org2")
	require.NoError(t, err)

	graph = swidAdmin.Graph()
	ok, err := graph.Exists("pt1")
	require.NoError(t, err)
	require.True(t, ok)

	children, err := graph.GetChildren("pt1")
	require.NoError(t, err)
	require.Contains(t, children, assetpap.LicenseObject("test-asset-id", "1"))

	parents, err := graph.GetParents("pt1")
	require.NoError(t, err)
	require.Contains(t, parents, "Org2_OA")
	require.Contains(t, parents, rbac.SwIDsOA)
	require.Contains(t, parents, status.SwIDsOA)
}
