package api

import (
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/adminmsp"
	"github.com/usnistgov/blossom/chaincode/collections"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
)

const Org2MSP = "Org2MSP"
const Org3MSP = "Org3MSP"

var Org2Collection = collections.Account(Org2MSP)
var Org3Collection = collections.Account(Org3MSP)

func newTestStub(t *testing.T) *mocks.Ctx {
	ctx := mocks.NewCtx()
	ctx.CreateCollection(collections.Catalog(),
		[]string{Org2MSP, Org3MSP, adminmsp.AdminMSP},
		[]string{adminmsp.AdminMSP})
	ctx.CreateCollection(collections.Account(Org2MSP),
		[]string{Org2MSP, adminmsp.AdminMSP},
		[]string{Org2MSP, adminmsp.AdminMSP})
	ctx.CreateCollection(collections.Account(Org3MSP),
		[]string{Org3MSP, adminmsp.AdminMSP},
		[]string{Org3MSP, adminmsp.AdminMSP})
	ctx.CreateCollection(collections.Licenses(),
		[]string{adminmsp.AdminMSP},
		[]string{adminmsp.AdminMSP})

	bcc := BlossomSmartContract{}
	err := ctx.SetClientIdentity(mocks.Super)
	require.NoError(t, err)
	err = bcc.InitNGAC(ctx)
	require.NoError(t, err)

	return ctx
}

func requestTestAccount(t *testing.T, ctx *mocks.Ctx, account string) {
	bcc := BlossomSmartContract{}
	if account == Org2MSP {
		err := ctx.SetClientIdentity(mocks.Org2SystemOwner)
		require.NoError(t, err)
	} else {
		err := ctx.SetClientIdentity(mocks.Org3SystemOwner)
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

	err = ctx.SetClientIdentity(mocks.Org2SystemOwner)
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
