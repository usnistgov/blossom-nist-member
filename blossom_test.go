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

func TestInit(t *testing.T) {
	bcc := BlossomSmartContract{}
	stub := mocks.NewMemCCStub()
	stub.CreateCollection(CatalogCollection(), []string{BlossomMSP}, []string{BlossomMSP})

	err := stub.SetUser(mocks.Super)
	require.NoError(t, err)

	t.Run("error - init without admin msp arg", func(t *testing.T) {
		stub.SetFunctionAndArgs("init")
		result := bcc.Init(stub)
		require.Equal(t, int32(500), result.Status)
	})

	t.Run("init with admin msp arg", func(t *testing.T) {
		stub.SetFunctionAndArgs("init", "BlossomMSP")
		result := bcc.Init(stub)
		require.Equal(t, int32(200), result.Status, result.Message)
	})
}
