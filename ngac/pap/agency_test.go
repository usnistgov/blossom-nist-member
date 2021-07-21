package pap

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac"
	agencypap "github.com/usnistgov/blossom/chaincode/ngac/pap/agency"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy"
	dacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/policy/dac"
	rbacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/policy/rbac"
	statuspolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/policy/status"
	"testing"
)

func TestRequestAccount(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	mock := mocks.New()
	mock.SetGraphState(graph)

	require.NoError(t, err)

	agencyAdmin, err := NewAgencyAdmin(mock.Stub)
	require.NoError(t, err)
	agency := model.Agency{
		Name:  "Org2",
		ATO:   "",
		MSPID: "",
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
		require.Contains(t, parents, dacpolicy.UserAttributeName)
		ok, err = graph.Exists(agencypap.ObjectAttributeName(agency.Name))
		require.NoError(t, err)
		require.True(t, ok)
		parents, err = graph.GetParents(agencypap.ObjectAttributeName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, parents, dacpolicy.ObjectAttributeName)
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
		require.Contains(t, parents, rbacpolicy.AgenciesOA)
		parents, err = graph.GetParents(ngac.FormatUsername(agency.Users.SystemOwner, agency.MSPID))
		require.NoError(t, err)
		require.Contains(t, parents, rbacpolicy.SystemOwnerUA)
		parents, err = graph.GetParents(ngac.FormatUsername(agency.Users.AcquisitionSpecialist, agency.MSPID))
		require.NoError(t, err)
		require.Contains(t, parents, rbacpolicy.AcquisitionSpecialistUA)
		parents, err = graph.GetParents(ngac.FormatUsername(agency.Users.SystemAdministrator, agency.MSPID))
		require.NoError(t, err)
		require.Contains(t, parents, rbacpolicy.SystemAdministratorUA)

	})

	t.Run("test Status", func(t *testing.T) {
		parents, err := graph.GetParents(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, parents, statuspolicy.PendingUA)
		parents, err = graph.GetParents(agencypap.InfoObjectName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, parents, statuspolicy.AgenciesOA)
	})
}

func TestUpdateAgencyStatus(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	mock := mocks.New()
	mock.SetGraphState(graph)

	require.NoError(t, err)

	agencyAdmin, err := NewAgencyAdmin(mock.Stub)
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
	err = agencyAdmin.RequestAccount(mock.Stub, agency)
	require.NoError(t, err)

	// update mock graph
	graph = agencyAdmin.graph
	mock.SetGraphState(graph)

	t.Run("test approved", func(t *testing.T) {
		err = agencyAdmin.UpdateAgencyStatus(mock.Stub, agency.Name, model.Approved)
		require.NoError(t, err)
		parents, err := agencyAdmin.graph.GetParents(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.Contains(t, parents, statuspolicy.ActiveUA)
		require.NotContains(t, parents, statuspolicy.InactiveUA)
		require.NotContains(t, parents, statuspolicy.PendingUA)
	})
	t.Run("test pending", func(t *testing.T) {
		err = agencyAdmin.UpdateAgencyStatus(mock.Stub, agency.Name, model.PendingATO)
		require.NoError(t, err)
		parents, err := agencyAdmin.graph.GetParents(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.NotContains(t, parents, statuspolicy.ActiveUA)
		require.NotContains(t, parents, statuspolicy.InactiveUA)
		require.Contains(t, parents, statuspolicy.PendingUA)
	})
	t.Run("test inactive", func(t *testing.T) {
		err = agencyAdmin.UpdateAgencyStatus(mock.Stub, agency.Name, model.InactiveATO)
		require.NoError(t, err)
		parents, err := agencyAdmin.graph.GetParents(agencypap.UserAttributeName(agency.Name))
		require.NoError(t, err)
		require.NotContains(t, parents, statuspolicy.ActiveUA)
		require.Contains(t, parents, statuspolicy.InactiveUA)
		require.NotContains(t, parents, statuspolicy.PendingUA)
	})
}
