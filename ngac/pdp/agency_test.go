package pdp

import (
	"encoding/json"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/api/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
	"testing"
	"time"
)

func TestUploadATO(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	graphBytes, err := testGraph(t)
	require.NoError(t, err)

	decider := NewAgencyDecider()

	t.Run("test a1 system owner", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org2MSP", nil)
		clientIdentity.GetX509CertificateReturns(A1SystemOwnerCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		chaincodeStub.GetStateReturns(graphBytes, nil)

		err = decider.setup(transactionContext)
		require.NoError(t, err)

		err = decider.UploadATO(transactionContext, "Org2", "test ato")
		require.NoError(t, err)
	})

	t.Run("test a1 system admin", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org2MSP", nil)
		clientIdentity.GetX509CertificateReturns(A1SystemAdminCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		chaincodeStub.GetStateReturns(graphBytes, nil)

		err = decider.setup(transactionContext)
		require.NoError(t, err)

		err = decider.UploadATO(transactionContext, "Org2", "test ato")
		require.Error(t, err)
	})
}

func TestUpdateAgencyStatus(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	graphBytes, err := testGraph(t)
	require.NoError(t, err)

	decider := NewAgencyDecider()

	t.Run("test a1 system owner", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org2MSP", nil)
		clientIdentity.GetX509CertificateReturns(A1SystemOwnerCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		chaincodeStub.GetStateReturns(graphBytes, nil)

		err = decider.setup(transactionContext)
		require.NoError(t, err)

		err = decider.UpdateAgencyStatus(transactionContext, "Org2", "test")
		require.Error(t, err)
	})

	t.Run("test a1 system admin", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org2MSP", nil)
		clientIdentity.GetX509CertificateReturns(A1SystemAdminCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		chaincodeStub.GetStateReturns(graphBytes, nil)

		err = decider.setup(transactionContext)
		require.NoError(t, err)

		err = decider.UpdateAgencyStatus(transactionContext, "Org2", "test")
		require.Error(t, err)
	})

	t.Run("test Org1 Admin", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org1MSP", nil)
		clientIdentity.GetX509CertificateReturns(Org1AdminCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		chaincodeStub.GetStateReturns(graphBytes, nil)

		err = decider.setup(transactionContext)
		require.NoError(t, err)

		err = decider.UpdateAgencyStatus(transactionContext, "Org2", "test")
		require.NoError(t, err)
	})
}

func TestFilterAgency(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	graphBytes, err := testGraph(t)
	require.NoError(t, err)

	decider := NewAgencyDecider()
	exp := time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)

	t.Run("test a1 system owner", func(t *testing.T) {
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org2MSP", nil)
		clientIdentity.GetX509CertificateReturns(A1SystemOwnerCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		chaincodeStub.GetStateReturns(graphBytes, nil)

		err = decider.setup(transactionContext)
		require.NoError(t, err)

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
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org2MSP", nil)
		clientIdentity.GetX509CertificateReturns(A1SystemAdminCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		chaincodeStub.GetStateReturns(graphBytes, nil)

		err = decider.setup(transactionContext)
		require.NoError(t, err)

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
		clientIdentity := &mocks.ClientIdentity{}
		clientIdentity.GetMSPIDReturns("Org1MSP", nil)
		clientIdentity.GetX509CertificateReturns(Org1AdminCert(), nil)
		transactionContext.GetClientIdentityReturns(clientIdentity)

		chaincodeStub.GetStateReturns(graphBytes, nil)

		err = decider.setup(transactionContext)
		require.NoError(t, err)

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

func testGraph(t *testing.T) ([]byte, error) {
	graph := memory.NewGraph()

	pc1, err := graph.CreateNode("pc1", pip.PolicyClass, nil)
	require.NoError(t, err)
	oa1, err := graph.CreateNode("oa1", pip.ObjectAttribute, nil)
	require.NoError(t, err)
	agencyInfoObj, err := graph.CreateNode(agencypap.InfoObjectName("Org2"), pip.Object, nil)
	require.NoError(t, err)
	adminUA, err := graph.CreateNode("adminUA", pip.UserAttribute, nil)
	require.NoError(t, err)
	nonAdminUA, err := graph.CreateNode("nonAdminUA", pip.UserAttribute, nil)
	require.NoError(t, err)
	superUA, err := graph.CreateNode("superUA", pip.UserAttribute, nil)
	require.NoError(t, err)
	a1SystemOwner, err := graph.CreateNode("a1_system_owner:Org2MSP", pip.User, nil)
	require.NoError(t, err)
	a1SystemAdmin, err := graph.CreateNode("a1_system_admin:Org2MSP", pip.User, nil)
	require.NoError(t, err)
	superUser, err := graph.CreateNode("Org1 Admin:Org1MSP", pip.User, nil)
	require.NoError(t, err)

	err = graph.Assign(oa1.Name, pc1.Name)
	require.NoError(t, err)
	err = graph.Assign(agencyInfoObj.Name, oa1.Name)
	require.NoError(t, err)
	err = graph.Assign(a1SystemOwner.Name, adminUA.Name)
	require.NoError(t, err)
	err = graph.Assign(a1SystemAdmin.Name, nonAdminUA.Name)
	require.NoError(t, err)
	err = graph.Assign(superUser.Name, superUA.Name)
	require.NoError(t, err)

	err = graph.Associate(adminUA.Name, oa1.Name, pip.ToOps(operations.UploadATO, operations.ViewAgency,
		operations.ViewATO, operations.ViewMSPID, operations.ViewUsers, operations.ViewStatus, operations.ViewAgencyLicenses))
	require.NoError(t, err)
	err = graph.Associate(nonAdminUA.Name, oa1.Name, pip.ToOps(operations.ViewAgency, operations.ViewAgencyLicenses))
	require.NoError(t, err)
	err = graph.Associate(superUA.Name, oa1.Name, pip.ToOps(pip.AllOps))
	require.NoError(t, err)

	return json.Marshal(graph)
}
