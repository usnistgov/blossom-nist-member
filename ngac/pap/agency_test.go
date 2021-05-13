package pap

import (
	"encoding/json"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/api/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/dac"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/rbac"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/status"
	"testing"
)

func TestRequestAccount(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)
	graphBytes, err := json.Marshal(graph)
	require.NoError(t, err)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.GetStateReturns(graphBytes, nil)

	require.NoError(t, err)

	agencyAdmin := NewAgencyAdmin()
	agency := model.Agency{
		Name:  "Org2",
		ATO:   "",
		MSPID: "",
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

	graph = agencyAdmin.graph

	t.Run("test DAC", func(t *testing.T) {
		ok, err := graph.Exists(ngac.FormatUsername(agency.Users.SystemOwner, agency.MSPID))
		require.NoError(t, err)
		require.True(t, ok)
		ok, err = graph.Exists(ngac.FormatUsername(agency.Users.SystemAdministrator, agency.MSPID))
		require.NoError(t, err)
		require.True(t, ok)
		ok, err = graph.Exists(ngac.FormatUsername(agency.Users.AcquisitionSpecialist, agency.MSPID))
		require.NoError(t, err)
		require.True(t, ok)
		ok, err = graph.Exists(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.True(t, ok)
		children, err := graph.GetChildren(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, children, ngac.FormatUsername(agency.Users.SystemOwner, agency.MSPID))
		require.Contains(t, children, ngac.FormatUsername(agency.Users.SystemAdministrator, agency.MSPID))
		require.Contains(t, children, ngac.FormatUsername(agency.Users.AcquisitionSpecialist, agency.MSPID))
		parents, err := graph.GetParents(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, parents, dac.UserAttributeName)
		ok, err = graph.Exists(agencypap.ObjectAttributeName(agency.Name))
		require.NoError(t, err)
		require.True(t, ok)
		parents, err = graph.GetParents(agencypap.ObjectAttributeName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, parents, dac.ObjectAttributeName)
		ok, err = graph.Exists(agencypap.InfoObjectName(agency.Name))
		require.NoError(t, err)
		require.True(t, ok)
		parents, err = graph.GetParents(agencypap.InfoObjectName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, parents, agencypap.ObjectAttributeName(agency.Name))
		assocs, err := graph.GetAssociationsForSubject(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, assocs, agencypap.ObjectAttributeName(agency.Name))
		require.Contains(t, assocs[agencypap.ObjectAttributeName(agency.Name)], pip.AllOps)

	})

	t.Run("test RBAC", func(t *testing.T) {
		parents, err := graph.GetParents(agencypap.InfoObjectName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, parents, rbac.AgenciesOA)
		parents, err = graph.GetParents(ngac.FormatUsername(agency.Users.SystemOwner, agency.MSPID))
		require.NoError(t, err)
		require.Contains(t, parents, rbac.SystemOwnerUA)
		parents, err = graph.GetParents(ngac.FormatUsername(agency.Users.AcquisitionSpecialist, agency.MSPID))
		require.NoError(t, err)
		require.Contains(t, parents, rbac.AcquisitionSpecialistUA)
		parents, err = graph.GetParents(ngac.FormatUsername(agency.Users.SystemAdministrator, agency.MSPID))
		require.NoError(t, err)
		require.Contains(t, parents, rbac.SystemAdministratorUA)

	})

	t.Run("test Status", func(t *testing.T) {
		parents, err := graph.GetParents(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, parents, status.PendingUA)
		parents, err = graph.GetParents(agencypap.InfoObjectName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, parents, status.AgenciesOA)
	})
}

func TestUpdateAgencyStatus(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)
	graphBytes, err := json.Marshal(graph)
	require.NoError(t, err)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.GetStateReturns(graphBytes, nil)

	require.NoError(t, err)

	agencyAdmin := NewAgencyAdmin()
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

	// update mock graph
	graph = agencyAdmin.graph
	graphBytes, err = json.Marshal(graph)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(graphBytes, nil)

	t.Run("test approved", func(t *testing.T) {
		err = agencyAdmin.UpdateAgencyStatus(transactionContext, agency.Name, model.Approved)
		require.NoError(t, err)
		parents, err := agencyAdmin.graph.GetParents(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, parents, status.ActiveUA)
		require.NotContains(t, parents, status.InactiveUA)
		require.NotContains(t, parents, status.PendingUA)
	})
	t.Run("test pending", func(t *testing.T) {
		err = agencyAdmin.UpdateAgencyStatus(transactionContext, agency.Name, model.PendingATO)
		require.NoError(t, err)
		parents, err := agencyAdmin.graph.GetParents(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.NotContains(t, parents, status.ActiveUA)
		require.NotContains(t, parents, status.InactiveUA)
		require.Contains(t, parents, status.PendingUA)
	})
	t.Run("test inactive", func(t *testing.T) {
		err = agencyAdmin.UpdateAgencyStatus(transactionContext, agency.Name, model.InactiveATO)
		require.NoError(t, err)
		parents, err := agencyAdmin.graph.GetParents(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.NotContains(t, parents, status.ActiveUA)
		require.Contains(t, parents, status.InactiveUA)
		require.NotContains(t, parents, status.PendingUA)
	})
}
