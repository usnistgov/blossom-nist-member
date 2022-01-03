package pdp

import (
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
	"testing"
)

func TestInitCatalogNGAC(t *testing.T) {
	stub := mocks.NewMemCCStub()
	err := stub.PutState(pap.AdminMSPKey, []byte("BlossomMSP"))
	require.NoError(t, err)

	stub.CreateCollection("catalog_coll", []string{"BlossomMSP"}, []string{"BlossomMSP"})
	err = stub.SetUser(mocks.Super)
	require.NoError(t, err)

	err = InitCatalogNGAC(stub, "catalog_coll")
	require.NoError(t, err)

	err = stub.SetUser(mocks.A1SystemOwner)
	require.NoError(t, err)

	err = InitCatalogNGAC(stub, "catalog_coll")
	require.Error(t, err)
}
