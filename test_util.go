package main

import (
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
	"time"
)

const A1MSP = "A1MSP"
const A2MSP = "A2MSP"
const BlossomMSP = "BlossomMSP"

var A1Collection = AccountCollection(A1MSP)
var A2Collection = AccountCollection(A1MSP)

func newTestStub(t *testing.T) *mocks.MemChaincodeStub {
	stub := mocks.NewMemCCStub()
	stub.CreateCollection(CatalogCollection(),
		[]string{A1MSP, A2MSP, BlossomMSP},
		[]string{BlossomMSP})
	stub.CreateCollection(AccountCollection(A1MSP),
		[]string{A1MSP, BlossomMSP},
		[]string{A1MSP, BlossomMSP})
	stub.CreateCollection(AccountCollection(A2MSP),
		[]string{A2MSP, BlossomMSP},
		[]string{A2MSP, BlossomMSP})
	stub.CreateCollection(LicensesCollection(),
		[]string{BlossomMSP},
		[]string{BlossomMSP})

	bcc := BlossomSmartContract{}
	stub.SetFunctionAndArgs("", "BlossomMSP")
	bcc.Init(stub)
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
		err := stub.SetUser(mocks.A1SystemOwner)
		require.NoError(t, err)
		err = stub.SetTransient("account", accountTransientInput{"a1_system_owner", "a1_system_admin", "a1_acq_spec"})
		require.NoError(t, err)
	} else {
		err := stub.SetUser(mocks.A2SystemOwner)
		require.NoError(t, err)
		err = stub.SetTransient("account", accountTransientInput{"a2_system_owner", "a2_system_admin", "a2_acq_spec"})
		require.NoError(t, err)
	}
	result := bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status)

	err := stub.SetUser(mocks.Super)
	require.NoError(t, err)

	stub.SetFunctionAndArgs("ApproveAccount", account)
	result = bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status)

	acct, err := bcc.Account(stub, account)
	require.NoError(t, err)
	require.Equal(t, model.PendingATO, acct.Status)

	err = stub.SetUser(mocks.A1SystemOwner)
	require.NoError(t, err)

	err = stub.SetTransient("ato", uploadATOTransientInput{ATO: "test ato"})
	require.NoError(t, err)
	err = bcc.UploadATO(stub)
	require.NoError(t, err)
	require.Equal(t, model.PendingATO, acct.Status)

	// udpate account status to authorized as super user
	err = stub.SetUser(mocks.Super)
	require.NoError(t, err)

	stub.SetFunctionAndArgs("UpdateAccountStatus", account, "AUTHORIZED")
	result = bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status, result.Message)
}

func onboardTestAsset(t *testing.T, stub *mocks.MemChaincodeStub, id, name string, licenses []string) {
	bcc := BlossomSmartContract{}
	stub.SetFunctionAndArgs("OnboardAsset", id, name, time.Now().AddDate(5, 0, 0).String())
	err := stub.SetTransient("asset", onboardAssetTransientInput{Licenses: licenses})
	require.NoError(t, err)
	result := bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status)
}
