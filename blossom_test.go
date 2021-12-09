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
		require.NoError(t, mock.SetUser(mocks.A1SystemOwner))
		mock.SetFunctionAndArgs("test", "hello world")
		result := bcc.Invoke(mock)
		require.Equal(t, int32(500), result.Status)
		require.Equal(t, "ngac not initialized", result.Message)
	})

	t.Run("test after initngac", func(t *testing.T) {
		bcc := new(BlossomSmartContract)
		mock := mocks.NewMemCCStub()
		mock.CreateCollection(CatalogCollectionName(), []string{"BlossomMSP", "A1MSP", "A2MSP"}, []string{"BlossomMSP"})

		require.NoError(t, mock.SetUser(mocks.Super))
		mock.SetFunctionAndArgs("InitNGAC")
		result := bcc.Invoke(mock)
		require.Equal(t, int32(200), result.Status)
		require.Equal(t, "", result.Message)
		require.NoError(t, mock.SetUser(mocks.A1SystemOwner))

		mock.SetFunctionAndArgs("test", "awesome blossom")
		result = bcc.Invoke(mock)
		require.Equal(t, int32(200), result.Status)
		require.Equal(t, "", result.Message)
	})

	t.Run("test initngac unauthorized", func(t *testing.T) {
		bcc := new(BlossomSmartContract)
		mock := mocks.NewMemCCStub()
		mock.CreateCollection(CatalogCollectionName(), []string{"BlossomMSP", "A1MSP", "A2MSP"}, []string{"BlossomMSP"})

		require.NoError(t, mock.SetUser(mocks.A1SystemAdmin))
		mock.SetFunctionAndArgs("InitNGAC")
		result := bcc.Invoke(mock)
		require.Equal(t, int32(500), result.Status)
		require.Equal(t, "user a1_system_admin:A1MSP does not have permission init_blossom on blossom_object", result.Message)
	})

}
