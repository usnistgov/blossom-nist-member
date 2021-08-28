package pap

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac"
	accountpap "github.com/usnistgov/blossom/chaincode/ngac/pap/account"
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

	accountAdmin, err := NewAccountAdmin(mock.Stub)
	require.NoError(t, err)
	account := &model.Account{
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
	err = accountAdmin.RequestAccount(mock.Stub, account)
	require.NoError(t, err)

	graph = accountAdmin.graph

	t.Run("test DAC", func(t *testing.T) {
		ok, err := graph.Exists(ngac.FormatUsername(account.Users.SystemOwner, account.MSPID))
		require.NoError(t, err)
		require.True(t, ok)
		ok, err = graph.Exists(ngac.FormatUsername(account.Users.SystemAdministrator, account.MSPID))
		require.NoError(t, err)
		require.True(t, ok)
		ok, err = graph.Exists(ngac.FormatUsername(account.Users.AcquisitionSpecialist, account.MSPID))
		require.NoError(t, err)
		require.True(t, ok)
		ok, err = graph.Exists(accountpap.UserAttributeName(account.Name))
		require.NoError(t, err)
		require.True(t, ok)
		children, err := graph.GetChildren(accountpap.UserAttributeName(account.Name))
		require.NoError(t, err)
		require.Contains(t, children, ngac.FormatUsername(account.Users.SystemOwner, account.MSPID))
		require.Contains(t, children, ngac.FormatUsername(account.Users.SystemAdministrator, account.MSPID))
		require.Contains(t, children, ngac.FormatUsername(account.Users.AcquisitionSpecialist, account.MSPID))
		parents, err := graph.GetParents(accountpap.UserAttributeName(account.Name))
		require.NoError(t, err)
		require.Contains(t, parents, dacpolicy.UserAttributeName)
		ok, err = graph.Exists(accountpap.ObjectAttributeName(account.Name))
		require.NoError(t, err)
		require.True(t, ok)
		parents, err = graph.GetParents(accountpap.ObjectAttributeName(account.Name))
		require.NoError(t, err)
		require.Contains(t, parents, dacpolicy.ObjectAttributeName)
		ok, err = graph.Exists(accountpap.InfoObjectName(account.Name))
		require.NoError(t, err)
		require.True(t, ok)
		parents, err = graph.GetParents(accountpap.InfoObjectName(account.Name))
		require.NoError(t, err)
		require.Contains(t, parents, accountpap.ObjectAttributeName(account.Name))
		assocs, err := graph.GetAssociationsForSubject(accountpap.UserAttributeName(account.Name))
		require.NoError(t, err)
		require.Contains(t, assocs, accountpap.ObjectAttributeName(account.Name))
		require.Contains(t, assocs[accountpap.ObjectAttributeName(account.Name)], pip.AllOps)

	})

	t.Run("test RBAC", func(t *testing.T) {
		parents, err := graph.GetParents(accountpap.InfoObjectName(account.Name))
		require.NoError(t, err)
		require.Contains(t, parents, rbacpolicy.AccountsOA)
		parents, err = graph.GetParents(ngac.FormatUsername(account.Users.SystemOwner, account.MSPID))
		require.NoError(t, err)
		require.Contains(t, parents, rbacpolicy.SystemOwnerUA)
		parents, err = graph.GetParents(ngac.FormatUsername(account.Users.AcquisitionSpecialist, account.MSPID))
		require.NoError(t, err)
		require.Contains(t, parents, rbacpolicy.AcquisitionSpecialistUA)
		parents, err = graph.GetParents(ngac.FormatUsername(account.Users.SystemAdministrator, account.MSPID))
		require.NoError(t, err)
		require.Contains(t, parents, rbacpolicy.SystemAdministratorUA)

	})

	t.Run("test Status", func(t *testing.T) {
		parents, err := graph.GetParents(accountpap.UserAttributeName(account.Name))
		require.NoError(t, err)
		require.Contains(t, parents, statuspolicy.PendingUA)
		parents, err = graph.GetParents(accountpap.InfoObjectName(account.Name))
		require.NoError(t, err)
		require.Contains(t, parents, statuspolicy.AccountsOA)
	})
}

func TestUpdateAccountStatus(t *testing.T) {
	graph := memory.NewGraph()
	err := policy.Configure(graph)
	require.NoError(t, err)

	mock := mocks.New()
	mock.SetGraphState(graph)

	require.NoError(t, err)

	accountAdmin, err := NewAccountAdmin(mock.Stub)
	require.NoError(t, err)
	account := &model.Account{
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
	err = accountAdmin.RequestAccount(mock.Stub, account)
	require.NoError(t, err)

	// update mock graph
	graph = accountAdmin.graph
	mock.SetGraphState(graph)

	t.Run("test approved", func(t *testing.T) {
		err = accountAdmin.UpdateAccountStatus(mock.Stub, account.Name, model.Approved)
		require.NoError(t, err)
		parents, err := accountAdmin.graph.GetParents(accountpap.UserAttributeName(account.Name))
		require.NoError(t, err)
		require.Contains(t, parents, statuspolicy.ActiveUA)
		require.NotContains(t, parents, statuspolicy.InactiveUA)
		require.NotContains(t, parents, statuspolicy.PendingUA)
	})
	t.Run("test pending", func(t *testing.T) {
		err = accountAdmin.UpdateAccountStatus(mock.Stub, account.Name, model.PendingATO)
		require.NoError(t, err)
		parents, err := accountAdmin.graph.GetParents(accountpap.UserAttributeName(account.Name))
		require.NoError(t, err)
		require.NotContains(t, parents, statuspolicy.ActiveUA)
		require.NotContains(t, parents, statuspolicy.InactiveUA)
		require.Contains(t, parents, statuspolicy.PendingUA)
	})
	t.Run("test inactive", func(t *testing.T) {
		err = accountAdmin.UpdateAccountStatus(mock.Stub, account.Name, model.InactiveATO)
		require.NoError(t, err)
		parents, err := accountAdmin.graph.GetParents(accountpap.UserAttributeName(account.Name))
		require.NoError(t, err)
		require.NotContains(t, parents, statuspolicy.ActiveUA)
		require.Contains(t, parents, statuspolicy.InactiveUA)
		require.NotContains(t, parents, statuspolicy.PendingUA)
	})
}
