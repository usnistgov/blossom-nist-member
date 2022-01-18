package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
)

func TestOnboardAsset(t *testing.T) {
	stub := newTestStub(t)

	err := stub.SetUser(mocks.Super)
	require.NoError(t, err)

	onboardTestAsset(t, stub, "123", "myasset", []string{"1", "2"})

	data, err := stub.GetPrivateData(CatalogCollection(), model.AssetKey("123"))
	require.NoError(t, err)

	assetPub := model.AssetPublic{}
	err = json.Unmarshal(data, &assetPub)
	require.NoError(t, err)
	require.Equal(t, "123", assetPub.ID)
	require.Equal(t, "myasset", assetPub.Name)
	require.Equal(t, 2, assetPub.Available)

	data, err = stub.GetPrivateData(LicensesCollection(), model.AssetKey("123"))
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
	stub := newTestStub(t)
	bcc := BlossomSmartContract{}

	onboardTestAsset(t, stub, "123", "myasset", []string{"1", "2"})

	stub.SetFunctionAndArgs("OffboardAsset", "123")
	result := bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status, result.Message)

	data, err := stub.GetPrivateData(CatalogCollection(), model.AssetKey("123"))
	require.NoError(t, err)
	require.Nil(t, data)

	data, err = stub.GetPrivateData(LicensesCollection(), model.AssetKey("123"))
	require.NoError(t, err)
	require.Nil(t, data)
}

func TestGetAssets(t *testing.T) {
	stub := newTestStub(t)
	bcc := BlossomSmartContract{}

	onboardTestAsset(t, stub, "123", "myasset1", []string{"1", "2"})
	onboardTestAsset(t, stub, "321", "myasset2", []string{"1", "2"})

	assets, err := bcc.GetAssets(stub)
	require.NoError(t, err)
	require.Equal(t, 2, len(assets))
}

func TestGetAsset(t *testing.T) {
	stub := newTestStub(t)
	bcc := BlossomSmartContract{}

	err := stub.SetUser(mocks.Super)
	require.NoError(t, err)

	onboardTestAsset(t, stub, "123", "myasset", []string{"1", "2"})

	asset, err := bcc.GetAsset(stub, "123")
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
	stub := newTestStub(t)

	bcc := BlossomSmartContract{}
	onboardTestAsset(t, stub, "123", "myasset", []string{"1", "2"})

	requestTestAccount(t, stub, A1MSP)

	err := stub.SetUser(mocks.Super)
	require.NoError(t, err)

	err = bcc.UpdateAccountStatus(stub, A1MSP, "AUTHORIZED")
	require.NoError(t, err)

	t.Run("error unauthorized to request checkout", func(t *testing.T) {
		err = stub.SetUser(mocks.A1SystemOwner)
		require.NoError(t, err)

		stub.SetFunctionAndArgs("RequestCheckout")
		err = stub.SetTransient("checkout", requestCheckoutTransientInput{"123", 1})
		require.NoError(t, err)
		result := bcc.Invoke(stub)
		require.Equal(t, int32(500), result.Status)
	})

	t.Run("authorized request checkout", func(t *testing.T) {
		err = stub.SetUser(mocks.A1SystemAdmin)
		require.NoError(t, err)
		stub.SetFunctionAndArgs("RequestCheckout")
		err = stub.SetTransient("checkout", requestCheckoutTransientInput{"123", 1})
		require.NoError(t, err)
		result := bcc.Invoke(stub)
		require.Equal(t, int32(200), result.Status)

		err = stub.SetUser(mocks.Super)
		require.NoError(t, err)

		stub.SetFunctionAndArgs("ApproveCheckout")
		err = stub.SetTransient("checkout", approveCheckoutTransientInput{A1MSP, "123"})
		require.NoError(t, err)
		result = bcc.Invoke(stub)
		require.Equal(t, int32(200), result.Status)

		err = stub.SetUser(mocks.A1SystemAdmin)
		require.NoError(t, err)

		stub.SetFunctionAndArgs("GetLicenses", A1MSP, "123")
		result = bcc.Invoke(stub)
		require.Equal(t, int32(200), result.Status)

		licenses := make(map[string]string, 0)
		err = json.Unmarshal(result.Payload, &licenses)
		require.NoError(t, err)

		err = stub.SetUser(mocks.Super)
		require.NoError(t, err)

		info, err := bcc.GetAsset(stub, "123")
		require.NoError(t, err)
		require.Equal(t, 2, info.TotalAmount)
		require.Equal(t, "123", info.ID)
		require.Equal(t, "myasset", info.Name)
		require.Equal(t, 1, len(info.AvailableLicenses))
		require.Equal(t, 1, info.Available)
		require.Equal(t, map[string]map[string]string{A1MSP: {"1": licenses["1"]}}, info.CheckedOut)

		// check in
		require.NoError(t, stub.SetUser(mocks.A1SystemAdmin))
		stub.SetFunctionAndArgs("InitiateCheckin")
		require.NoError(t, stub.SetTransient("checkin", initiateCheckinTransientInput{
			AssetID:  "123",
			Licenses: []string{"1"},
		}))
		result = bcc.Invoke(stub)
		require.Equal(t, int32(200), result.Status, result.Message)

		require.NoError(t, stub.SetUser(mocks.Super))
		stub.SetFunctionAndArgs("ProcessCheckin")
		require.NoError(t, stub.SetTransient("checkin", processCheckinTransientInput{
			Account: A1MSP,
			AssetID: "123",
		}))
		result = bcc.Invoke(stub)
		require.Equal(t, int32(200), result.Status, result.Message)

		stub.SetFunctionAndArgs("GetLicenses", A1MSP, "123")
		result = bcc.Invoke(stub)
		require.Equal(t, int32(200), result.Status)

		licenses = make(map[string]string, 0)
		err = json.Unmarshal(result.Payload, &licenses)
		require.NoError(t, err)
		require.Equal(t, 0, len(licenses))

		// update account to pending
		err = bcc.UpdateAccountStatus(stub, A1MSP, "PENDING_ATO")
		require.NoError(t, err)

		// checkout should fail
		err = stub.SetUser(mocks.A1SystemAdmin)
		require.NoError(t, err)

		stub.SetFunctionAndArgs("RequestCheckout")
		err = stub.SetTransient("checkout", requestCheckoutTransientInput{"123", 1})
		require.NoError(t, err)
		result = bcc.Invoke(stub)
		require.Equal(t, int32(500), result.Status)
	})
}

func TestCheckoutRequests(t *testing.T) {
	stub := newTestStub(t)

	onboardTestAsset(t, stub, "123", "myasset1", []string{"1", "2"})
	onboardTestAsset(t, stub, "456", "myasset2", []string{"1", "2"})

	requestTestAccount(t, stub, A1MSP)

	bcc := BlossomSmartContract{}
	err := stub.SetUser(mocks.A1SystemAdmin)
	require.NoError(t, err)
	err = stub.SetTransient("checkout", requestCheckoutTransientInput{"123", 1})
	require.NoError(t, err)
	err = bcc.RequestCheckout(stub)
	require.NoError(t, err)

	err = stub.SetTransient("checkout", requestCheckoutTransientInput{"456", 1})
	require.NoError(t, err)
	err = bcc.RequestCheckout(stub)
	require.NoError(t, err)

	result, err := bcc.GetCheckoutRequests(stub, A1MSP)
	require.NoError(t, err)
	require.Equal(t, 2, len(result))
}
