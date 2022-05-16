package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/collections"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
)

func TestOnboardAsset(t *testing.T) {
	ctx := newTestStub(t)

	err := ctx.SetClientIdentity(mocks.Super)
	require.NoError(t, err)

	onboardTestAsset(t, ctx, "123", "myasset", []string{"1", "2"})

	data, err := ctx.GetStub().GetPrivateData(collections.Catalog(), model.AssetKey("123"))
	require.NoError(t, err)

	assetPub := model.AssetPublic{}
	err = json.Unmarshal(data, &assetPub)
	require.NoError(t, err)
	require.Equal(t, "123", assetPub.ID)
	require.Equal(t, "myasset", assetPub.Name)
	require.Equal(t, 2, assetPub.Available)

	data, err = ctx.GetStub().GetPrivateData(collections.Licenses(), model.AssetKey("123"))
	require.NoError(t, err)

	assetPvt := model.AssetPrivate{}
	err = json.Unmarshal(data, &assetPvt)
	require.NoError(t, err)
	require.Equal(t, 2, assetPvt.TotalAmount)
	require.Equal(t, map[string]string{"1": "exp", "2": "exp"}, assetPvt.Licenses)
	require.Equal(t, 2, len(assetPvt.AvailableLicenses))
	require.Empty(t, assetPvt.CheckedOut)
}

func TestOffboardAsset(t *testing.T) {
	ctx := newTestStub(t)
	bcc := BlossomSmartContract{}

	onboardTestAsset(t, ctx, "123", "myasset", []string{"1", "2"})

	err := bcc.OffboardAsset(ctx, "123")
	require.NoError(t, err)

	data, err := ctx.GetStub().GetPrivateData(collections.Catalog(), model.AssetKey("123"))
	require.NoError(t, err)
	require.Nil(t, data)

	data, err = ctx.GetStub().GetPrivateData(collections.Licenses(), model.AssetKey("123"))
	require.NoError(t, err)
	require.Nil(t, data)
}

func TestGetAssets(t *testing.T) {
	ctx := newTestStub(t)
	bcc := BlossomSmartContract{}

	onboardTestAsset(t, ctx, "123", "myasset1", []string{"1", "2"})
	onboardTestAsset(t, ctx, "321", "myasset2", []string{"1", "2"})

	assets, err := bcc.GetAssets(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(assets))
}

func TestGetAsset(t *testing.T) {
	ctx := newTestStub(t)
	bcc := BlossomSmartContract{}

	err := ctx.SetClientIdentity(mocks.Super)
	require.NoError(t, err)

	onboardTestAsset(t, ctx, "123", "myasset", []string{"1", "2"})

	asset, err := bcc.GetAsset(ctx, "123")
	require.NoError(t, err)
	require.Equal(t, "123", asset.ID)
	require.Equal(t, "myasset", asset.Name)
	require.Equal(t, 2, asset.TotalAmount)
	require.Equal(t, 2, asset.Available)
	require.Equal(t, 2, len(asset.AvailableLicenses))
	require.Equal(t, map[string]string{"1": "exp", "2": "exp"}, asset.Licenses)
	require.Empty(t, asset.CheckedOut)
}

func TestCheckout(t *testing.T) {
	ctx := newTestStub(t)

	bcc := BlossomSmartContract{}
	onboardTestAsset(t, ctx, "123", "myasset", []string{"1", "2"})

	requestTestAccount(t, ctx, A1MSP)

	err := ctx.SetClientIdentity(mocks.Super)
	require.NoError(t, err)

	err = bcc.UpdateAccountStatus(ctx, A1MSP, "AUTHORIZED")
	require.NoError(t, err)

	t.Run("error unauthorized to request checkout", func(t *testing.T) {
		err = ctx.SetClientIdentity(mocks.A1SystemOwner)
		require.NoError(t, err)

		err = ctx.SetTransient("checkout", requestCheckoutTransientInput{"123", 1})
		require.NoError(t, err)
		err = bcc.RequestCheckout(ctx)
		require.Error(t, err)
	})

	t.Run("authorized request checkout", func(t *testing.T) {
		err = ctx.SetClientIdentity(mocks.A1SystemAdmin)
		require.NoError(t, err)
		err = ctx.SetTransient("checkout", requestCheckoutTransientInput{"123", 1})
		require.NoError(t, err)
		err = bcc.RequestCheckout(ctx)
		require.NoError(t, err)

		err = ctx.SetClientIdentity(mocks.Super)
		require.NoError(t, err)

		err = ctx.SetTransient("checkout", approveCheckoutTransientInput{A1MSP, "123"})
		require.NoError(t, err)
		err = bcc.ApproveCheckout(ctx)
		require.NoError(t, err)

		err = ctx.SetClientIdentity(mocks.A1SystemAdmin)
		require.NoError(t, err)

		licenses := make(map[string]string, 0)
		licenses, err = bcc.GetLicenses(ctx, A1MSP, "123")
		require.NoError(t, err)
		require.Equal(t, 1, len(licenses))

		err = ctx.SetClientIdentity(mocks.Super)
		require.NoError(t, err)

		info := &model.Asset{}
		info, err = bcc.GetAsset(ctx, "123")
		require.NoError(t, err)
		require.Equal(t, 2, info.TotalAmount)
		require.Equal(t, "123", info.ID)
		require.Equal(t, "myasset", info.Name)
		require.Equal(t, 1, len(info.AvailableLicenses))
		require.Equal(t, 1, info.Available)
		require.Equal(t, map[string]map[string]string{A1MSP: {"1": licenses["1"]}}, info.CheckedOut)

		// check in
		require.NoError(t, ctx.SetClientIdentity(mocks.A1SystemAdmin))
		require.NoError(t, ctx.SetTransient("checkin", initiateCheckinTransientInput{
			AssetID:  "123",
			Licenses: []string{"1"},
		}))
		err = bcc.InitiateCheckin(ctx)
		require.NoError(t, err)

		require.NoError(t, ctx.SetClientIdentity(mocks.Super))
		require.NoError(t, ctx.SetTransient("checkin", processCheckinTransientInput{
			Account: A1MSP,
			AssetID: "123",
		}))
		err = bcc.ProcessCheckin(ctx)
		require.NoError(t, err)

		licenses, err = bcc.GetLicenses(ctx, A1MSP, "123")
		require.NoError(t, err)

		require.Equal(t, 0, len(licenses))

		// update account to pending
		err = bcc.UpdateAccountStatus(ctx, A1MSP, "PENDING_ATO")
		require.NoError(t, err)

		// checkout should fail
		err = ctx.SetClientIdentity(mocks.A1SystemAdmin)
		require.NoError(t, err)

		err = ctx.SetTransient("checkout", requestCheckoutTransientInput{"123", 1})
		require.NoError(t, err)
		err = bcc.RequestCheckout(ctx)
		require.Error(t, err)
	})
}

func TestCheckoutRequests(t *testing.T) {
	ctx := newTestStub(t)

	onboardTestAsset(t, ctx, "123", "myasset1", []string{"1", "2"})
	onboardTestAsset(t, ctx, "456", "myasset2", []string{"1", "2"})

	requestTestAccount(t, ctx, A1MSP)

	bcc := BlossomSmartContract{}
	err := ctx.SetClientIdentity(mocks.A1SystemAdmin)
	require.NoError(t, err)
	err = ctx.SetTransient("checkout", requestCheckoutTransientInput{"123", 1})
	require.NoError(t, err)
	err = bcc.RequestCheckout(ctx)
	require.NoError(t, err)

	err = ctx.SetTransient("checkout", requestCheckoutTransientInput{"456", 1})
	require.NoError(t, err)
	err = bcc.RequestCheckout(ctx)
	require.NoError(t, err)

	result, err := bcc.GetCheckoutRequests(ctx, A1MSP)
	require.NoError(t, err)
	require.Equal(t, 2, len(result))
}

func TestViewAssetPermissions(t *testing.T) {
	ctx := newTestStub(t)
	requestTestAccount(t, ctx, A1MSP)
	requestTestAccount(t, ctx, A2MSP)
	require.NoError(t, ctx.SetClientIdentity(mocks.Super))
	onboardTestAsset(t, ctx, "123", "myasset1", []string{"1", "2"})
	onboardTestAsset(t, ctx, "456", "myasset2", []string{"1", "2"})

	bcc := BlossomSmartContract{}
	assets, err := bcc.GetAssets(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(assets))

	require.NoError(t, ctx.SetClientIdentity(mocks.A1SystemAdmin))
	assets, err = bcc.GetAssets(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(assets))

	// update account status
	require.NoError(t, ctx.SetClientIdentity(mocks.Super))
	err = bcc.UpdateAccountStatus(ctx, A1MSP, "UNAUTHORIZED_DENIED")
	require.NoError(t, err)

	require.NoError(t, ctx.SetClientIdentity(mocks.A1SystemAdmin))
	assets, err = bcc.GetAssets(ctx)
	require.Error(t, err)

	require.NoError(t, ctx.SetClientIdentity(mocks.A1SystemAdmin))
	_, err = bcc.GetAsset(ctx, "123")
	require.Error(t, err)
}
