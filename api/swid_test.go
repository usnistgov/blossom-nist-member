package api

import (
	"encoding/json"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/api/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
	"time"
)

func TestGetSwIDsAssociatedWithLicense(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	swidBytes := make([][]byte, 0)
	swid := model.SwID{
		PrimaryTag:      "pt1",
		XML:             "xml",
		Asset:           "test-license",
		License:         "test-license:1",
		LeaseExpiration: time.Time{},
	}
	b, err := json.Marshal(swid)
	require.NoError(t, err)
	swidBytes = append(swidBytes, b)

	swid = model.SwID{
		PrimaryTag:      "pt2",
		XML:             "xml",
		Asset:           "test-license",
		License:         "test-license:2",
		LeaseExpiration: time.Time{},
	}
	b, err = json.Marshal(swid)
	require.NoError(t, err)
	swidBytes = append(swidBytes, b)

	swid = model.SwID{
		PrimaryTag:      "pt3",
		XML:             "xml",
		Asset:           "other-license",
		License:         "other-license:1",
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

	chaincodeStub.GetStateByRangeReturns(iterator, nil)

	cc := BlossomSmartContract{}
	swids, err := cc.getSwIDsAssociatedWithAsset(transactionContext, "test-license")
	require.NoError(t, err)
	require.Equal(t, 2, len(swids))
}
