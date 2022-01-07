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

	onboardTestAsset(t, stub, "123", "myasset", []string{"1", "2"})
	require.NoError(t, err)

	requestTestAccount(t, stub, A1MSP)

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

	// report swid on license they did not checkout
	stub.SetFunctionAndArgs("ReportSwID")
	err = stub.SetTransient("swid", reportSwIDTransientInput{
		PrimaryTag: "primary_tag_1",
		Asset:      "123",
		License:    "2",
		Xml:        "swid_xml",
	})
	require.NoError(t, err)
	result = bcc.Invoke(stub)
	require.Equal(t, int32(500), result.Status)

	stub.SetFunctionAndArgs("ReportSwID")
	err = stub.SetTransient("swid", reportSwIDTransientInput{
		PrimaryTag: "primary_tag_1",
		Asset:      "123",
		License:    "1",
		Xml:        "swid_xml",
	})
	require.NoError(t, err)
	result = bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status)

	// check swid in collection
	stub.SetFunctionAndArgs("GetSwID")
	err = stub.SetTransient("swid", swidTransientInput{
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
		Asset:      "123",
		License:    "1",
	}, swid)

	stub.SetFunctionAndArgs("GetSwIDsAssociatedWithAsset", A1MSP, "123")
	result = bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status)
	swids := make([]*model.SwID, 0)
	err = json.Unmarshal(result.Payload, &swids)
	require.NoError(t, err)
	require.Equal(t, 1, len(swids))
	require.Equal(t, &model.SwID{
		PrimaryTag: "primary_tag_1",
		XML:        "swid_xml",
		Asset:      "123",
		License:    "1",
	}, swids[0])

	err = stub.SetUser(mocks.A1SystemOwner)
	require.NoError(t, err)

	// try deleting as unauthorized user
	stub.SetFunctionAndArgs("DeleteSwID")
	err = stub.SetTransient("swid", swidTransientInput{
		Account:    A1MSP,
		PrimaryTag: "primary_tag_1",
	})
	require.NoError(t, err)
	result = bcc.Invoke(stub)
	require.Equal(t, int32(500), result.Status)

	err = stub.SetUser(mocks.A1SystemAdmin)
	require.NoError(t, err)

	stub.SetFunctionAndArgs("DeleteSwID")
	err = stub.SetTransient("swid", swidTransientInput{
		Account:    A1MSP,
		PrimaryTag: "primary_tag_1",
	})
	require.NoError(t, err)
	result = bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status)

	stub.SetFunctionAndArgs("GetSwID")
	err = stub.SetTransient("swid", swidTransientInput{
		Account:    A1MSP,
		PrimaryTag: "primary_tag_1",
	})
	require.NoError(t, err)
	result = bcc.Invoke(stub)
	require.Equal(t, int32(500), result.Status)
}
