package main

import (
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"testing"
	"time"
)

const A1MSP = "A1MSP"
const A2MSP = "A2MSP"
const BlossomMSP = "BlossomMSP"

var A1Collection = AccountCollectionName(A1MSP)
var A2Collection = AccountCollectionName(A1MSP)

func newTestStub(t *testing.T) *mocks.MemChaincodeStub {
	stub := mocks.NewMemCCStub()
	stub.CreateCollection(CatalogCollectionName(),
		[]string{A1MSP, A2MSP, BlossomMSP},
		[]string{BlossomMSP})
	stub.CreateCollection(AccountCollectionName(A1MSP),
		[]string{A1MSP, BlossomMSP},
		[]string{A1MSP, BlossomMSP})
	stub.CreateCollection(AccountCollectionName(A2MSP),
		[]string{A2MSP, BlossomMSP},
		[]string{A2MSP, BlossomMSP})
	stub.CreateCollection(LicensesCollectionName(),
		[]string{BlossomMSP},
		[]string{BlossomMSP})

	bcc := BlossomSmartContract{}
	err := stub.SetUser(mocks.Super)
	require.NoError(t, err)
	err = bcc.handleInitNGAC(stub)
	require.NoError(t, err)

	return stub
}

func requestTestAccount(t *testing.T, stub *mocks.MemChaincodeStub, account string) {
	bcc := BlossomSmartContract{}
	stub.SetFunctionAndArgs("RequestAccount")
	if account == A1MSP {
		err := stub.SetTransient("account", accountTransientInput{"a1_system_owner", "a1_system_admin", "a1_acq_spec", "my ato"})
		require.NoError(t, err)
	} else {
		err := stub.SetTransient("account", accountTransientInput{"a2_system_owner", "a2_system_admin", "a2_acq_spec", "my ato"})
		require.NoError(t, err)
	}
	result := bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status)
}

func onboardTestAsset(t *testing.T, stub *mocks.MemChaincodeStub, id, name string, licenses []string) {
	bcc := BlossomSmartContract{}
	stub.SetFunctionAndArgs("OnboardAsset", id, name, time.Now().AddDate(5, 0, 0).String())
	err := stub.SetTransient("asset", onboardAssetTransientInput{Licenses: licenses})
	require.NoError(t, err)
	result := bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status)
}
