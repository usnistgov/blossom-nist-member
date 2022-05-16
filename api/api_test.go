package api

import (
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/collections"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
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
		mock := mocks.NewCtx()
		mock.CreateCollection(collections.Catalog(), []string{"BlossomMSP"}, []string{"BlossomMSP"})
		require.NoError(t, mock.SetClientIdentity(mocks.A1SystemOwner))
		_, err := bcc.GetAssets(mock)
		require.Error(t, err)

		require.NoError(t, mock.SetClientIdentity(mocks.Super))
		err = mock.SetTransient("asset", onboardAssetTransientInput{Licenses: []model.License{
			{"1", "exp1"}, {"2", "exp2"},
		}})
		require.NoError(t, err)
		err = bcc.OnboardAsset(mock, "123", "asset1", "onboard-date", "expiration-date")
		require.Error(t, err)
	})

	t.Run("test after initngac", func(t *testing.T) {
		bcc := new(BlossomSmartContract)
		testCtx := mocks.NewCtx()
		testCtx.CreateCollection(collections.Catalog(), []string{"BlossomMSP", "A1MSP", "A2MSP"}, []string{"BlossomMSP"})

		require.NoError(t, testCtx.SetClientIdentity(mocks.Super))
		err := bcc.InitNGAC(testCtx)
		require.NoError(t, err)

		_, err = bcc.GetAssets(testCtx)
		require.NoError(t, err)
	})

	t.Run("test initngac unauthorized", func(t *testing.T) {
		bcc := new(BlossomSmartContract)
		testCtx := mocks.NewCtx()
		testCtx.CreateCollection(collections.Catalog(), []string{"BlossomMSP", "A1MSP", "A2MSP"}, []string{"BlossomMSP"})

		require.NoError(t, testCtx.SetClientIdentity(mocks.A1SystemAdmin))
		err := bcc.InitNGAC(testCtx)
		require.Error(t, err)
	})
}
