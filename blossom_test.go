package main

import (
	"encoding/json"
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy"
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
		mock := mocks.New()
		require.NoError(t, mock.SetUser(mocks.A1SystemOwner))
		mock.Stub.GetFunctionAndParametersReturns("test", []string{})
		result := bcc.Invoke(mock.Stub)
		require.Equal(t, int32(500), result.Status)
		require.Equal(t, "ngac not initialized", result.Message)
	})

	t.Run("test after initngac", func(t *testing.T) {
		bcc := new(BlossomSmartContract)
		mock := mocks.New()
		require.NoError(t, mock.SetUser(mocks.Super))
		mock.Stub.GetFunctionAndParametersReturns("InitNGAC", []string{})
		result := bcc.Invoke(mock.Stub)
		require.Equal(t, int32(200), result.Status)
		require.Equal(t, "", result.Message)

		g := memory.NewGraph()
		require.NoError(t, policy.Configure(g))
		mock.SetGraphState(g)

		require.NoError(t, mock.SetUser(mocks.A1SystemOwner))
		mock.Stub.GetFunctionAndParametersReturns("test", []string{"awesome blossom"})
		result = bcc.Invoke(mock.Stub)
		fmt.Println(result.Message)
		require.Equal(t, int32(200), result.Status)
		require.Equal(t, "", result.Message)
	})

	t.Run("test initngac unauthorized", func(t *testing.T) {
		bcc := new(BlossomSmartContract)
		mock := mocks.New()
		require.NoError(t, mock.SetUser(mocks.A1SystemAdmin))
		mock.Stub.GetFunctionAndParametersReturns("InitNGAC", []string{})
		result := bcc.Invoke(mock.Stub)
		fmt.Println(result.Message)
		require.Equal(t, int32(500), result.Status)
		require.Equal(t, "user a1_system_admin:A1MSP does not have permission to initialize blossom", result.Message)
	})

}

func TestUseCase(t *testing.T) {
	// initialize ngac with super user
	bcc := new(BlossomSmartContract)
	stub := mocks.NewMemCCStub()
	stub.SetUser(mocks.Super)
	stub.SetFunctionAndArgs("InitNGAC")
	result := bcc.Invoke(&stub)
	require.Equal(t, int32(200), result.Status)

	// try to initialize ngac with unauthorized user
	stub.SetUser(mocks.A1SystemOwner)
	stub.SetFunctionAndArgs("InitNGAC")
	result = bcc.Invoke(&stub)
	fmt.Println(result.Message)
	require.Equal(t, int32(500), result.Status)
	require.Equal(t, "user a1_system_owner:A1MSP does not have permission to initialize blossom", result.Message)

	// onboard asset as super user
	stub.SetUser(mocks.Super)
	asset := []byte(`
	{
	  "id": "test-asset-id",
	  "name": "test-asset",
	  "total_amount": 10,
	  "available": 10,
	  "cost": 100.00,
	  "onboarding_date": "",
	  "expiration": "2025-01-01",
	  "licenses": [
	    "test-asset-1",
	    "test-asset-2",
	    "test-asset-3",
	    "test-asset-4",
	    "test-asset-5",
	    "test-asset-6",
	    "test-asset-7",
	    "test-asset-8",
	    "test-asset-9",
	    "test-asset-10"
	  ],
	  "available_licenses": [],
	  "checked_out": {}
	}`)
	stub.SetFunctionAndArgs("OnboardAsset", asset)
	result = bcc.Invoke(&stub)
	require.Equal(t, int32(200), result.Status)

	// request account
	stub.SetUser(mocks.A1SystemOwner)
	account := []byte(`
{
  "name": "Agency1",
  "ato": "this is a test ato",
  "mspid": "A1MSP",
  "users": {
    "system_owner": "a1_system_owner",
    "acquisition_specialist": "a1_acq_spec",
    "system_administrator": "a1_system_admin"
  },
  "status": "",
  "assets": {}
}`)
	stub.SetFunctionAndArgs("RequestAccount", account)
	result = bcc.Invoke(&stub)
	require.Equal(t, int32(200), result.Status)

	// set account to active
	stub.SetUser(mocks.Super)
	stub.SetFunctionAndArgs("UpdateAccountStatus", []byte("Agency1"), []byte(model.Approved))
	result = bcc.Invoke(&stub)
	require.Equal(t, int32(200), result.Status)

	// get assets
	stub.SetUser(mocks.A1SystemAdmin)
	stub.SetFunctionAndArgs("Assets")
	result = bcc.Invoke(&stub)
	require.Equal(t, int32(200), result.Status)
	assets := make([]*model.Asset, 0)
	require.NoError(t, json.Unmarshal(result.Payload, &assets))
	require.Equal(t, 1, len(assets))
	require.Equal(t, "test-asset", assets[0].Name)
	require.Equal(t, 100.00, assets[0].Cost)
	require.Equal(t, 0, len(assets[0].Licenses))
	require.Equal(t, 0, len(assets[0].AvailableLicenses))

	// checkout asset
	stub.SetFunctionAndArgs("Checkout", []byte("test-asset-id"), []byte("Agency1"), []byte("2"))
	result = bcc.Invoke(&stub)
	require.Equal(t, int32(200), result.Status)
	checkedOut := make(map[string]model.DateTime)
	require.NoError(t, json.Unmarshal(result.Payload, &checkedOut))
	require.Equal(t, 2, len(checkedOut))

	// checkout asset with unauthorized user should fail
	stub.SetUser(mocks.A1AcqSpec)
	stub.SetFunctionAndArgs("Checkout", []byte("test-asset-id"), []byte("Agency1"), []byte("2"))
	result = bcc.Invoke(&stub)
	require.Equal(t, int32(500), result.Status)

	// report swid with valid license
	stub.SetUser(mocks.A1SystemAdmin)
	swid := []byte(fmt.Sprintf(`
{
	"primary_tag": "swid-1",
	"xml": "<swid>test</swid>",
	"asset": "test-asset-id",
	"license": "test-asset-1",
	"lease_expiration": "%s"
}
`, checkedOut["test-asset-1"]))
	stub.SetFunctionAndArgs("ReportSwID", swid, []byte("Agency1"))
	result = bcc.Invoke(&stub)
	fmt.Println(result.Message)
	require.Equal(t, int32(200), result.Status)

	// report swid with invalid license
	swid = []byte(fmt.Sprintf(`
{
	"primary_tag": "swid-1",
	"xml": "<swid>test</swid>",
	"asset": "test-asset-id",
	"license": "test-asset-3",
	"lease_expiration": "%s"
}
`, checkedOut["test-asset-1"]))
	stub.SetFunctionAndArgs("ReportSwID", swid, []byte("Agency1"))
	result = bcc.Invoke(&stub)
	require.Equal(t, int32(500), result.Status)

	stub.SetFunctionAndArgs("GetSwIDsAssociatedWithAsset", []byte("test-asset-id"))
	result = bcc.Invoke(&stub)
	swids := make([]*model.SwID, 0)
	require.NoError(t, json.Unmarshal(result.Payload, &swids))
	require.Equal(t, 1, len(swids))
	require.Equal(t, "swid-1", swids[0].PrimaryTag)
}
