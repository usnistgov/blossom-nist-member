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

func TestUploadATO(t *testing.T) {
	decider := NewAgencyDecider()

	t.Run("test a1 system owner", func(t *testing.T) {
		chaincodeStub := &mocks.ChaincodeStub{}
		transactionContext := &mocks.TransactionContext{}
		transactionContext.GetStubReturns(chaincodeStub)

		err := initAgencyTestGraph(t, transactionContext, chaincodeStub)
		require.NoError(t, err)

		SetUser(transactionContext, A1SystemOwnerCert(), "Org2MSP")

		err = decider.UploadATO(transactionContext, "Org2")
		require.NoError(t, err)
	})

	t.Run("test a1 system admin", func(t *testing.T) {
		chaincodeStub := &mocks.ChaincodeStub{}
		transactionContext := &mocks.TransactionContext{}
		transactionContext.GetStubReturns(chaincodeStub)

		err := initAgencyTestGraph(t, transactionContext, chaincodeStub)
		require.NoError(t, err)

		SetUser(transactionContext, A1SystemAdminCert(), "Org2MSP")

		err = decider.UploadATO(transactionContext, "Org2")
		require.Error(t, err)
	})
}

func TestUpdateAgencyStatus(t *testing.T) {
	decider := NewAgencyDecider()

	t.Run("test a1 system owner", func(t *testing.T) {
		chaincodeStub := &mocks.ChaincodeStub{}
		transactionContext := &mocks.TransactionContext{}
		transactionContext.GetStubReturns(chaincodeStub)

		err := initAgencyTestGraph(t, transactionContext, chaincodeStub)
		require.NoError(t, err)

		SetUser(transactionContext, A1SystemOwnerCert(), "Org2MSP")

		err = decider.UpdateAgencyStatus(transactionContext, "Org2", "test")
		require.Error(t, err)
	})

	t.Run("test a1 system admin", func(t *testing.T) {
		chaincodeStub := &mocks.ChaincodeStub{}
		transactionContext := &mocks.TransactionContext{}
		transactionContext.GetStubReturns(chaincodeStub)

		err := initAgencyTestGraph(t, transactionContext, chaincodeStub)
		require.NoError(t, err)

		SetUser(transactionContext, A1SystemAdminCert(), "Org2MSP")

		err = decider.UpdateAgencyStatus(transactionContext, "Org2", "test")
		require.Error(t, err)
	})

	t.Run("test Org1 Admin", func(t *testing.T) {
		chaincodeStub := &mocks.ChaincodeStub{}
		transactionContext := &mocks.TransactionContext{}
		transactionContext.GetStubReturns(chaincodeStub)

		err := initAgencyTestGraph(t, transactionContext, chaincodeStub)
		require.NoError(t, err)

		SetUser(transactionContext, Org1AdminCert(), "Org1MSP")

		err = decider.UpdateAgencyStatus(transactionContext, "Org2", "test")
		require.NoError(t, err)
	})
}

func TestFilterAgency(t *testing.T) {
	decider := NewAgencyDecider()

	t.Run("test a1 system owner", func(t *testing.T) {
		chaincodeStub := &mocks.ChaincodeStub{}
		transactionContext := &mocks.TransactionContext{}
		transactionContext.GetStubReturns(chaincodeStub)

		err := initAgencyTestGraph(t, transactionContext, chaincodeStub)
		require.NoError(t, err)

		exp := time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
		SetUser(transactionContext, A1SystemOwnerCert(), "Org2MSP")

		agency := &model.Agency{
			Name:  "Org2",
			ATO:   "ato",
			MSPID: "Org2MSP",
			Users: model.Users{
				SystemOwner:           "a1_system_owner",
				SystemAdministrator:   "a1_system_admin",
				AcquisitionSpecialist: "a1_acq_spec",
			},
			Status: "status",
			Licenses: map[string]map[string]time.Time{
				"license1": {
					"k1": exp,
					"k2": exp,
				},
			},
		}

		err = decider.setup(transactionContext)
		require.NoError(t, err)
		err = decider.filterAgency(agency)
		require.NoError(t, err)
		require.Equal(t, "Org2", agency.Name)
		require.Equal(t, "ato", agency.ATO)
		require.Equal(t, "Org2MSP", agency.MSPID)
		require.Equal(t, model.Status("status"), agency.Status)
		require.Equal(t, model.Users{
			SystemOwner:           "a1_system_owner",
			SystemAdministrator:   "a1_system_admin",
			AcquisitionSpecialist: "a1_acq_spec",
		}, agency.Users)
		require.Equal(t, map[string]map[string]time.Time{
			"license1": {
				"k1": exp,
				"k2": exp,
			},
		}, agency.Licenses)
	})

	t.Run("test a1 system admin", func(t *testing.T) {
		chaincodeStub := &mocks.ChaincodeStub{}
		transactionContext := &mocks.TransactionContext{}
		transactionContext.GetStubReturns(chaincodeStub)

		err := initAgencyTestGraph(t, transactionContext, chaincodeStub)
		require.NoError(t, err)

		exp := time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
		SetUser(transactionContext, A1SystemAdminCert(), "Org2MSP")

		agency := &model.Agency{
			Name:  "Org2",
			ATO:   "ato",
			MSPID: "Org2MSP",
			Users: model.Users{
				SystemOwner:           "a1_system_owner",
				SystemAdministrator:   "a1_system_admin",
				AcquisitionSpecialist: "a1_acq_spec",
			},
			Status: "status",
			Licenses: map[string]map[string]time.Time{
				"license1": {
					"k1": exp,
					"k2": exp,
				},
			},
		}

		err = decider.setup(transactionContext)
		require.NoError(t, err)
		err = decider.filterAgency(agency)
		require.NoError(t, err)
		require.Equal(t, "Org2", agency.Name)
		require.Equal(t, "", agency.ATO)
		require.Equal(t, "", agency.MSPID)
		require.Equal(t, model.Status(""), agency.Status)
		require.Equal(t, model.Users{}, agency.Users)
		require.Equal(t, map[string]map[string]time.Time{
			"license1": {
				"k1": exp,
				"k2": exp,
			},
		}, agency.Licenses)
	})

	t.Run("test Org1 Admin", func(t *testing.T) {
		chaincodeStub := &mocks.ChaincodeStub{}
		transactionContext := &mocks.TransactionContext{}
		transactionContext.GetStubReturns(chaincodeStub)

		err := initAgencyTestGraph(t, transactionContext, chaincodeStub)
		require.NoError(t, err)

		exp := time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
		SetUser(transactionContext, Org1AdminCert(), "Org1MSP")

		agency := &model.Agency{
			Name:  "Org2",
			ATO:   "ato",
			MSPID: "Org2MSP",
			Users: model.Users{
				SystemOwner:           "a1_system_owner",
				SystemAdministrator:   "a1_system_admin",
				AcquisitionSpecialist: "a1_acq_spec",
			},
			Status: "status",
			Licenses: map[string]map[string]time.Time{
				"license1": {
					"k1": exp,
					"k2": exp,
				},
			},
		}

		err = decider.setup(transactionContext)
		require.NoError(t, err)
		err = decider.filterAgency(agency)
		require.NoError(t, err)
		require.Equal(t, "Org2", agency.Name)
		require.Equal(t, "ato", agency.ATO)
		require.Equal(t, "Org2MSP", agency.MSPID)
		require.Equal(t, model.Status("status"), agency.Status)
		require.Equal(t, model.Users{
			SystemOwner:           "a1_system_owner",
			SystemAdministrator:   "a1_system_admin",
			AcquisitionSpecialist: "a1_acq_spec",
		}, agency.Users)
		require.Equal(t, map[string]map[string]time.Time{
			"license1": {
				"k1": exp,
				"k2": exp,
			},
		}, agency.Licenses)
	})
}

func initAgencyTestGraph(t *testing.T, ctx *mocks.TransactionContext, stub *mocks.ChaincodeStub) error {
	graph := memory.NewGraph()

	// configure the policy
	err := policy.Configure(graph)
	require.NoError(t, err)

	// add an account
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

	SetGraphState(t, stub, graph)

	// add account as the a1 system owner
	SetUser(ctx, A1SystemOwnerCert(), "Org2MSP")
	agencyDecider := NewAgencyDecider()
	err = agencyDecider.RequestAccount(ctx, agency)
	require.NoError(t, err)

	SetGraphState(t, stub, agencyDecider.pap.Graph())

	return nil
}
