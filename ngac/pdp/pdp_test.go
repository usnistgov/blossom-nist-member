package pdp

import (
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"testing"
)

func TestInitCatalogNGAC(t *testing.T) {
	stub := mocks.NewMemCCStub()

	stub.CreateCollection("catalog_coll", []string{"BlossomMSP"}, []string{"BlossomMSP"})
	err := stub.SetUser(mocks.Super)
	require.NoError(t, err)

	err = InitCatalogNGAC(stub)
	require.NoError(t, err)

	err = stub.SetUser(mocks.A1SystemOwner)
	require.NoError(t, err)

	err = InitCatalogNGAC(stub)
	require.Error(t, err)
}
