package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"testing"
)

//go:generate counterfeiter -o ../mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o ../mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

//go:generate counterfeiter -o ../mocks/clientIdentity.go -fake-name ClientIdentity . clientIdentity
type clientIdentity interface {
	cid.ClientIdentity
}

func TestInitNGAC(t *testing.T) {
	t.Run("test without initngac", func(t *testing.T) {
		bcc := new(BlossomSmartContract)
		mock := mocks.NewMemCCStub()
		mock.CreateCollection(CatalogCollection(), []string{"BlossomMSP"}, []string{"BlossomMSP"})
		require.NoError(t, mock.SetUser(mocks.A1SystemOwner))
		mock.SetFunctionAndArgs("test", "hello world")
		result := bcc.Invoke(mock)
		require.Equal(t, int32(200), result.Status)

		require.NoError(t, mock.SetUser(mocks.Super))
		mock.SetFunctionAndArgs("OnboardAsset", "123", "asset1", "onboard-date", "expiration-date")
		err := mock.SetTransient("asset", onboardAssetTransientInput{Licenses: map[string]string{"1": "exp1", "2": "exp2"}})
		require.NoError(t, err)
		result = bcc.Invoke(mock)
		require.Equal(t, int32(500), result.Status, result.Message)
	})

	t.Run("test after initngac", func(t *testing.T) {
		bcc := new(BlossomSmartContract)
		stub := mocks.NewMemCCStub()
		stub.CreateCollection(CatalogCollection(), []string{"BlossomMSP", "A1MSP", "A2MSP"}, []string{"BlossomMSP"})

		require.NoError(t, stub.SetUser(mocks.Super))
		stub.SetFunctionAndArgs("InitNGAC")
		result := bcc.Invoke(stub)
		require.Equal(t, int32(200), result.Status)
		require.Equal(t, "", result.Message)
		require.NoError(t, stub.SetUser(mocks.A1SystemOwner))

		stub.SetFunctionAndArgs("test", "awesome blossom")
		result = bcc.Invoke(stub)
		require.Equal(t, int32(200), result.Status)
		require.Equal(t, "", result.Message)
	})

	t.Run("test initngac unauthorized", func(t *testing.T) {
		bcc := new(BlossomSmartContract)
		stub := mocks.NewMemCCStub()
		stub.CreateCollection(CatalogCollection(), []string{"BlossomMSP", "A1MSP", "A2MSP"}, []string{"BlossomMSP"})

		require.NoError(t, stub.SetUser(mocks.A1SystemAdmin))
		stub.SetFunctionAndArgs("InitNGAC")
		result := bcc.Invoke(stub)
		require.Equal(t, int32(500), result.Status)
	})

}
