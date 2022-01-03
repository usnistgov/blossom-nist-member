package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
)

func TestSwID(t *testing.T) {
	stub := newTestStub(t)
	err := stub.SetUser(mocks.Super)
	require.NoError(t, err)

	bcc := BlossomSmartContract{}

	err = bcc.handleInitNGAC(stub)
	require.NoError(t, err)

	onboardTestAsset(t, stub, "123", "myasset", []string{"1", "2"})
	require.NoError(t, err)

	requestTestAccount(t, stub, A1MSP)

	err = bcc.UpdateAccountStatus(stub, A1MSP, "ACTIVE")
	require.NoError(t, err)

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

	stub.SetFunctionAndArgs("ReportSwID")
	err = stub.SetTransient("swid", reportSwIDTransientInput{
		Account:    A1MSP,
		PrimaryTag: "primary_tag_1",
		Asset:      "myasset1",
		License:    "1",
		Xml:        "swid_xml",
	})
	require.NoError(t, err)
	result = bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status)

	// check swid in collection
	stub.SetFunctionAndArgs("GetSwID")
	err = stub.SetTransient("swid", getSwIDTransientInput{
		Account:    A1MSP,
		PrimaryTag: "primary_tag_1",
	})
	require.NoError(t, err)
	result = bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status)
	swid := &model.SwID{}
	err = json.Unmarshal(result.Payload, swid)
	require.NoError(t, err)
	require.Equal(t, &model.SwID{
		PrimaryTag: "primary_tag_1",
		XML:        "swid_xml",
		Asset:      "myasset1",
		License:    "1",
	}, swid)

	stub.SetFunctionAndArgs("GetSwIDsAssociatedWithAsset")
	err = stub.SetTransient("swid", getSwIDsAssociatedWithAssetTransientInput{
		Account: A1MSP,
		AssetID: "myasset1",
	})
	require.NoError(t, err)
	result = bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status)
	swids := make([]*model.SwID, 0)
	err = json.Unmarshal(result.Payload, &swids)
	require.NoError(t, err)
	require.Equal(t, 1, len(swids))
	require.Equal(t, &model.SwID{
		PrimaryTag: "primary_tag_1",
		XML:        "swid_xml",
		Asset:      "myasset1",
		License:    "1",
	}, swids[0])
}
