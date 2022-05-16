package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
)

func TestRequestAccount(t *testing.T) {
	t.Run("test with ato", func(t *testing.T) {
		ctx := newTestStub(t)

		requestTestAccount(t, ctx, A1MSP)

		bytes, err := ctx.GetStub().GetState(model.AccountKey("A1MSP"))
		require.NoError(t, err)

		acctPub := &model.AccountPublic{}
		err = json.Unmarshal(bytes, acctPub)
		require.NoError(t, err)
		require.Equal(t, A1MSP, acctPub.Name)
		require.Equal(t, A1MSP, acctPub.MSPID)
		require.Equal(t, model.Authorized, acctPub.Status)

		bytes, err = ctx.GetStub().GetPrivateData(A1Collection, model.AccountKey("A1MSP"))
		require.NoError(t, err)

		acctPvt := &model.AccountPrivate{}
		err = json.Unmarshal(bytes, acctPvt)
		require.NoError(t, err)
	})

	t.Run("test without ato", func(t *testing.T) {
		ctx := newTestStub(t)

		err := ctx.SetClientIdentity(mocks.A1SystemOwner)
		require.NoError(t, err)

		bcc := BlossomSmartContract{}
		err = ctx.SetTransient("account", accountTransientInput{"a1_system_owner", "a1_system_admin", "a1_acq_spec"})
		require.NoError(t, err)
		err = bcc.RequestAccount(ctx)
		require.NoError(t, err)

		bytes, err := ctx.GetStub().GetState(model.AccountKey(A1MSP))
		require.NoError(t, err)

		acctPub := &model.AccountPublic{}
		err = json.Unmarshal(bytes, acctPub)
		require.NoError(t, err)
		require.Equal(t, A1MSP, acctPub.Name)
		require.Equal(t, A1MSP, acctPub.MSPID)
		require.Equal(t, model.PendingApproval, acctPub.Status)

		bytes, err = ctx.GetStub().GetPrivateData(A1Collection, model.AccountKey(A1MSP))
		require.NoError(t, err)

		acctPvt := &model.AccountPrivate{}
		err = json.Unmarshal(bytes, acctPvt)
		require.NoError(t, err)
	})

}

func TestUploadATO(t *testing.T) {
	ctx := newTestStub(t)

	requestTestAccount(t, ctx, A1MSP)

	err := ctx.SetClientIdentity(mocks.A1SystemOwner)
	require.NoError(t, err)

	bcc := BlossomSmartContract{}
	err = ctx.SetTransient("ato", uploadATOTransientInput{"my ato"})
	require.NoError(t, err)
	err = bcc.UploadATO(ctx)
	require.NoError(t, err)

	account, err := bcc.GetAccount(ctx, A1MSP)
	require.NoError(t, err)
	require.Equal(t, "my ato", account.ATO)

	err = ctx.SetClientIdentity(mocks.A1AcqSpec)
	require.NoError(t, err)

	err = ctx.SetTransient("ato", uploadATOTransientInput{"my ato"})
	require.NoError(t, err)
	err = bcc.UploadATO(ctx)
	require.Error(t, err)
	require.Errorf(t, err, "error uploading ATO for account A1MSP: user a1_acq_spec does not have permission upload_ato on A1MSP_object", err.Error())
}

func TestUpdateAccountStatus(t *testing.T) {
	ctx := newTestStub(t)

	requestTestAccount(t, ctx, A1MSP)

	err := ctx.SetClientIdentity(mocks.A1SystemOwner)
	require.NoError(t, err)

	bcc := BlossomSmartContract{}
	err = bcc.UpdateAccountStatus(ctx, A1MSP, "AUTHORIZED")
	require.Error(t, err)

	err = ctx.SetClientIdentity(mocks.Super)
	require.NoError(t, err)

	err = bcc.UpdateAccountStatus(ctx, A1MSP, "AUTHORIZED")
	require.NoError(t, err)
}

func TestAccounts(t *testing.T) {
	ctx := newTestStub(t)

	requestTestAccount(t, ctx, A1MSP)
	requestTestAccount(t, ctx, A2MSP)

	err := ctx.SetClientIdentity(mocks.A2SystemOwner)
	require.NoError(t, err)

	bcc := BlossomSmartContract{}
	accounts, err := bcc.GetAccounts(ctx)
	require.NoError(t, err)

	require.Equal(t, 2, len(accounts))
}

func TestAccount(t *testing.T) {
	ctx := newTestStub(t)

	requestTestAccount(t, ctx, A1MSP)
	requestTestAccount(t, ctx, A2MSP)

	err := ctx.SetClientIdentity(mocks.A2SystemOwner)
	require.NoError(t, err)

	bcc := BlossomSmartContract{}
	acct, err := bcc.GetAccount(ctx, A1MSP)
	require.NoError(t, err)
	require.Equal(t, A1MSP, acct.Name)
	require.Equal(t, A1MSP, acct.MSPID)
	require.Equal(t, model.Authorized, acct.Status)
	require.Equal(t, "", acct.ATO)
	require.Empty(t, acct.Assets)

	acct, err = bcc.GetAccount(ctx, A2MSP)
	require.NoError(t, err)
	require.Equal(t, A2MSP, acct.Name)
	require.Equal(t, A2MSP, acct.MSPID)
	require.Equal(t, model.Authorized, acct.Status)
	require.Empty(t, acct.Assets)
}
