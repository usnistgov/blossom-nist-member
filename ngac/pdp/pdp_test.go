package pdp

import (
	"github.com/stretchr/testify/require"
	mocks "github.com/usnistgov/blossom/chaincode/mocks"
	"testing"
)

func TestInitGraph(t *testing.T) {
	// create the mock chaincode stub and context
	mock := mocks.New()

	// create the mock client identity
	mock.SetUser(mocks.Super)

	adminDecider := NewAdminDecider()
	err := adminDecider.InitGraph(mock.Stub)
	// no error should occur since Org1 Admin has permission to init blossom
	require.NoError(t, err)

	// do the same with an unauthorized user - a1_system_owner
	mock.SetUser(mocks.A1SystemOwner)

	// make a call to init blossom
	err = adminDecider.InitGraph(mock.Stub)
	// an error should occur because a1_system_owner cannot init blossom
	require.Error(t, err)
}
