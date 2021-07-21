package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
	"time"
)

func TestGetSwIDsAssociatedWithLicense(t *testing.T) {
	mock := mocks.New()

	swidBytes := make([][]byte, 0)
	swid := model.SwID{
		PrimaryTag:      "pt1",
		XML:             "xml",
		Asset:           "test-asset",
		License:         "test-asset:1",
		LeaseExpiration: time.Time{},
	}
	b, err := json.Marshal(swid)
	require.NoError(t, err)
	swidBytes = append(swidBytes, b)

	swid = model.SwID{
		PrimaryTag:      "pt2",
		XML:             "xml",
		Asset:           "test-asset",
		License:         "test-asset:2",
		LeaseExpiration: time.Time{},
	}
	b, err = json.Marshal(swid)
	require.NoError(t, err)
	swidBytes = append(swidBytes, b)

	swid = model.SwID{
		PrimaryTag:      "pt3",
		XML:             "xml",
		Asset:           "other-asset",
		License:         "other-asset:1",
		LeaseExpiration: time.Time{},
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
	swids, err := cc.getSwIDsAssociatedWithAsset(mock.Stub, "test-asset")
	require.NoError(t, err)
	require.Equal(t, 2, len(swids))
}
