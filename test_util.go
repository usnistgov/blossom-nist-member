package main

import (
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/collections"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
)

const A1MSP = "A1MSP"
const A2MSP = "A2MSP"

var A1Collection = collections.Account(A1MSP)
var A2Collection = collections.Account(A1MSP)

func newTestStub(t *testing.T) *mocks.MemChaincodeStub {
	stub := mocks.NewMemCCStub()
	stub.CreateCollection(collections.Catalog(),
		[]string{A1MSP, A2MSP, "BlossomMSP"},
		[]string{"BlossomMSP"})
	stub.CreateCollection(collections.Account(A1MSP),
		[]string{A1MSP, "BlossomMSP"},
		[]string{A1MSP, "BlossomMSP"})
	stub.CreateCollection(collections.Account(A2MSP),
		[]string{A2MSP, "BlossomMSP"},
		[]string{A2MSP, "BlossomMSP"})
	stub.CreateCollection(collections.Licenses(),
		[]string{"BlossomMSP"},
		[]string{"BlossomMSP"})

	bcc := BlossomSmartContract{}
	err := stub.SetUser(mocks.Super)
	require.NoError(t, err)
	stub.SetFunctionAndArgs("InitNGAC")
	result := bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status, result.Message)

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
	require.Equal(t, int32(200), result.Status, result.Message)

	acct, err := bcc.GetAccount(stub, account)
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
	licensesMap := make([]model.License, 0)
	for _, l := range licenses {
		licensesMap = append(licensesMap, model.License{
			LicenseID:  l,
			Expiration: "exp",
		})
	}

	bcc := BlossomSmartContract{}
	stub.SetFunctionAndArgs("OnboardAsset", id, name, "onboard-date", "expiration-date")
	err := stub.SetTransient("asset", onboardAssetTransientInput{Licenses: licensesMap})
	require.NoError(t, err)
	result := bcc.Invoke(stub)
	require.Equal(t, int32(200), result.Status, result.Message)
}
