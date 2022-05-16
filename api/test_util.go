package api

import (
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/collections"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
)

const A1MSP = "A1MSP"
const A2MSP = "A2MSP"

var A1Collection = collections.Account(A1MSP)
var A2Collection = collections.Account(A2MSP)

func newTestStub(t *testing.T) *mocks.Ctx {
	ctx := mocks.NewCtx()
	ctx.CreateCollection(collections.Catalog(),
		[]string{A1MSP, A2MSP, "BlossomMSP"},
		[]string{"BlossomMSP"})
	ctx.CreateCollection(collections.Account(A1MSP),
		[]string{A1MSP, "BlossomMSP"},
		[]string{A1MSP, "BlossomMSP"})
	ctx.CreateCollection(collections.Account(A2MSP),
		[]string{A2MSP, "BlossomMSP"},
		[]string{A2MSP, "BlossomMSP"})
	ctx.CreateCollection(collections.Licenses(),
		[]string{"BlossomMSP"},
		[]string{"BlossomMSP"})

	bcc := BlossomSmartContract{}
	err := ctx.SetClientIdentity(mocks.Super)
	require.NoError(t, err)
	err = bcc.InitNGAC(ctx)
	require.NoError(t, err)

	return ctx
}

func requestTestAccount(t *testing.T, ctx *mocks.Ctx, account string) {
	bcc := BlossomSmartContract{}
	if account == A1MSP {
		err := ctx.SetClientIdentity(mocks.A1SystemOwner)
		require.NoError(t, err)
		err = ctx.SetTransient("account", accountTransientInput{"a1_system_owner", "a1_system_admin", "a1_acq_spec"})
		require.NoError(t, err)
	} else {
		err := ctx.SetClientIdentity(mocks.A2SystemOwner)
		require.NoError(t, err)
		err = ctx.SetTransient("account", accountTransientInput{"a2_system_owner", "a2_system_admin", "a2_acq_spec"})
		require.NoError(t, err)
	}
	err := bcc.RequestAccount(ctx)
	require.NoError(t, err)

	err = ctx.SetClientIdentity(mocks.Super)
	require.NoError(t, err)

	err = bcc.ApproveAccount(ctx, account)
	require.NoError(t, err)

	acct, err := bcc.GetAccount(ctx, account)
	require.NoError(t, err)
	require.Equal(t, model.PendingATO, acct.Status)

	err = ctx.SetClientIdentity(mocks.A1SystemOwner)
	require.NoError(t, err)

	err = ctx.SetTransient("ato", uploadATOTransientInput{ATO: "test ato"})
	require.NoError(t, err)
	err = bcc.UploadATO(ctx)
	require.NoError(t, err)
	require.Equal(t, model.PendingATO, acct.Status)

	// udpate account status to authorized as super user
	err = ctx.SetClientIdentity(mocks.Super)
	require.NoError(t, err)

	err = bcc.UpdateAccountStatus(ctx, account, "AUTHORIZED")
	require.NoError(t, err)
}

func onboardTestAsset(t *testing.T, ctx *mocks.Ctx, id, name string, licenses []string) {
	licensesMap := make([]model.License, 0)
	for _, l := range licenses {
		licensesMap = append(licensesMap, model.License{
			LicenseID:  l,
			Expiration: "exp",
		})
	}

	bcc := BlossomSmartContract{}
	err := ctx.SetTransient("asset", onboardAssetTransientInput{Licenses: licensesMap})
	require.NoError(t, err)
	err = bcc.OnboardAsset(ctx, id, name, "onboard-date", "expiration-date")
	require.NoError(t, err)
}
