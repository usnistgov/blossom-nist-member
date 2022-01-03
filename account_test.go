package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
)

func TestRequestAccount(t *testing.T) {
	t.Run("test with ato", func(t *testing.T) {
		stub := newTestStub(t)

		requestTestAccount(t, stub, A1MSP)

		bytes, err := stub.GetState(model.AccountKey("A1MSP"))
		require.NoError(t, err)

		acctPub := &model.AccountPublic{}
		err = json.Unmarshal(bytes, acctPub)
		require.NoError(t, err)
		require.Equal(t, A1MSP, acctPub.Name)
		require.Equal(t, A1MSP, acctPub.MSPID)
		require.Equal(t, model.Active, acctPub.Status)

		bytes, err = stub.GetPrivateData(A1Collection, model.AccountKey("A1MSP"))
		require.NoError(t, err)

		acctPvt := &model.AccountPrivate{}
		err = json.Unmarshal(bytes, acctPvt)
		require.NoError(t, err)
		require.Equal(t, model.Users{
			SystemOwner:           "a1_system_owner",
			AcquisitionSpecialist: "a1_acq_spec",
			SystemAdministrator:   "a1_system_admin",
		}, acctPvt.Users)
	})

	t.Run("test without ato", func(t *testing.T) {
		stub := newTestStub(t)

		err := stub.SetUser(mocks.A1SystemOwner)
		require.NoError(t, err)

		bcc := BlossomSmartContract{}
		stub.SetFunctionAndArgs("RequestAccount")
		err = stub.SetTransient("account", accountTransientInput{"a1_system_owner", "a1_system_admin", "a1_acq_spec", ""})
		require.NoError(t, err)
		result := bcc.Invoke(stub)
		require.Equal(t, int32(200), result.Status)

		bytes, err := stub.GetState(model.AccountKey(A1MSP))
		require.NoError(t, err)

		acctPub := &model.AccountPublic{}
		err = json.Unmarshal(bytes, acctPub)
		require.NoError(t, err)
		require.Equal(t, A1MSP, acctPub.Name)
		require.Equal(t, A1MSP, acctPub.MSPID)
		require.Equal(t, model.PendingApproval, acctPub.Status)

		bytes, err = stub.GetPrivateData(A1Collection, model.AccountKey(A1MSP))
		require.NoError(t, err)

		acctPvt := &model.AccountPrivate{}
		err = json.Unmarshal(bytes, acctPvt)
		require.NoError(t, err)
		require.Equal(t, model.Users{
			SystemOwner:           "a1_system_owner",
			AcquisitionSpecialist: "a1_acq_spec",
			SystemAdministrator:   "a1_system_admin",
		}, acctPvt.Users)
	})

}

func TestUploadATO(t *testing.T) {
	stub := newTestStub(t)

	requestTestAccount(t, stub, A1MSP)

	err := stub.SetUser(mocks.A1SystemOwner)
	require.NoError(t, err)

	bcc := BlossomSmartContract{}
	stub.SetFunctionAndArgs("UploadATO")
	err = stub.SetTransient("ato", uploadATOTransientInput{"my ato"})
	require.NoError(t, err)
	result := bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status, result.Message)

	account, err := bcc.Account(stub, A1MSP)
	require.NoError(t, err)
	require.Equal(t, "my ato", account.ATO)

	err = stub.SetUser(mocks.A1AcqSpec)
	require.NoError(t, err)

	err = stub.SetTransient("ato", uploadATOTransientInput{"my ato"})
	require.NoError(t, err)
	result = bcc.Invoke(stub)
	require.Equal(t, int32(500), result.Status)
	require.Equal(t, "error uploading ATO for account A1MSP: user \"a1_acq_spec:A1MSP\" does not have permission \"upload_ato\" on \"A1MSP_object\"", result.Message)
}

func TestUpdateAccountStatus(t *testing.T) {
	stub := newTestStub(t)

	requestTestAccount(t, stub, A1MSP)

	err := stub.SetUser(mocks.A1SystemOwner)
	require.NoError(t, err)

	bcc := BlossomSmartContract{}
	err = bcc.UpdateAccountStatus(stub, A1MSP, "ACTIVE")
	require.Error(t, err)

	err = stub.SetUser(mocks.Super)
	require.NoError(t, err)

	err = bcc.UpdateAccountStatus(stub, A1MSP, "ACTIVE")
	require.NoError(t, err)
}

func TestAccounts(t *testing.T) {
	stub := newTestStub(t)

	requestTestAccount(t, stub, A1MSP)
	requestTestAccount(t, stub, A2MSP)

	err := stub.SetUser(mocks.A2SystemOwner)
	require.NoError(t, err)

	bcc := BlossomSmartContract{}
	accounts, err := bcc.Accounts(stub)
	require.NoError(t, err)

	require.Equal(t, 2, len(accounts))
}

func TestAccount(t *testing.T) {
	stub := newTestStub(t)

	requestTestAccount(t, stub, A1MSP)
	requestTestAccount(t, stub, A2MSP)

	err := stub.SetUser(mocks.A2SystemOwner)
	require.NoError(t, err)

	bcc := BlossomSmartContract{}
	acct, err := bcc.Account(stub, A1MSP)
	require.NoError(t, err)
	require.Equal(t, A1MSP, acct.Name)
	require.Equal(t, A1MSP, acct.MSPID)
	require.Equal(t, model.Active, acct.Status)
	require.Equal(t, "", acct.ATO)
	require.Equal(t, model.Users{}, acct.Users)
	require.Empty(t, acct.Assets)

	acct, err = bcc.Account(stub, A2MSP)
	require.NoError(t, err)
	require.Equal(t, A2MSP, acct.Name)
	require.Equal(t, A2MSP, acct.MSPID)
	require.Equal(t, model.Active, acct.Status)
	require.Equal(t, "my ato", acct.ATO)
	require.Equal(t, model.Users{
		SystemOwner:           "a2_system_owner",
		AcquisitionSpecialist: "a2_acq_spec",
		SystemAdministrator:   "a2_system_admin",
	}, acct.Users)
	require.Empty(t, acct.Assets)
}
