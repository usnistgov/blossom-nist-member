package main

import (
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

	err = stub.SetUser(mocks.A1SystemOwner)
	require.NoError(t, err)

	requestTestAccount(t, stub, A1MSP)

	err = stub.SetUser(mocks.Super)
	require.NoError(t, err)

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

	err = bcc.ReportSwID(stub, A1MSP, "primary_tag_1", "myasset1", "1", "swid_xml")
	require.NoError(t, err)

	// check swid in collection
	swid, err := bcc.GetSwID(stub, A1MSP, "primary_tag_1")
	require.NoError(t, err)
	require.Equal(t, &model.SwID{
		PrimaryTag: "primary_tag_1",
		XML:        "swid_xml",
		Asset:      "myasset1",
		License:    "1",
	}, swid)

	swids, err := bcc.GetSwIDsAssociatedWithAsset(stub, A1MSP, "myasset1")
	require.NoError(t, err)
	require.Equal(t, 1, len(swids))
	require.Equal(t, &model.SwID{
		PrimaryTag: "primary_tag_1",
		XML:        "swid_xml",
		Asset:      "myasset1",
		License:    "1",
	}, swids[0])
}

/*func TestGetSwIDsAssociatedWithLicense(t *testing.T) {
	mock := mocks.New()

	swidBytes := make([][]byte, 0)
	swid := model.SwID{
		PrimaryTag:      "pt1",
		XML:             "xml",
		Asset:           "test-asset",
		License:         "test-asset:1",
		LeaseExpiration: "",
	}
	b, err := json.Marshal(swid)
	require.NoError(t, err)
	swidBytes = append(swidBytes, b)

	swid = model.SwID{
		PrimaryTag:      "pt2",
		XML:             "xml",
		Asset:           "test-asset",
		License:         "test-asset:2",
		LeaseExpiration: "",
	}
	b, err = json.Marshal(swid)
	require.NoError(t, err)
	swidBytes = append(swidBytes, b)

	swid = model.SwID{
		PrimaryTag:      "pt3",
		XML:             "xml",
		Asset:           "other-asset",
		License:         "other-asset:1",
		LeaseExpiration: "",
	}
	b, err = json.Marshal(swid)
	require.NoError(t, err)
	swidBytes = append(swidBytes, b)

	iterator := &mocks.StateQueryIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, true)
	iterator.HasNextReturnsOnCall(2, true)
	iterator.HasNextReturnsOnCall(3, false)
	iterator.NextReturnsOnCall(0, &queryresult.KV{Key: model.SwIDKey("pt1"), Value: swidBytes[0]}, nil)
	iterator.NextReturnsOnCall(1, &queryresult.KV{Key: model.SwIDKey("pt2"), Value: swidBytes[1]}, nil)
	iterator.NextReturnsOnCall(2, &queryresult.KV{Key: model.SwIDKey("pt3"), Value: swidBytes[2]}, nil)

	mock.Stub.GetStateByRangeReturns(iterator, nil)

	cc := BlossomSmartContract{}
	swids, err := cc.GetSwIDsAssociatedWithAsset(mock.Stub, "test-asset")
	require.NoError(t, err)
	require.Equal(t, 2, len(swids))
}
*/
