package api

import (
	"fmt"
	"github.com/hyperledger/fabric/core/shim"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/api/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
)

func TestRequestAccount(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	clientIdentity := &mocks.ClientIdentity{}
	clientIdentity.GetMSPIDReturns("Org1MSP", nil)
	clientIdentity.GetX509CertificateReturns(Org1AdminCert(), nil)
	transactionContext.GetClientIdentityReturns(clientIdentity)

	blossomCC := BlossomSmartContract{}
	err := blossomCC.InitBlossom(transactionContext)
	require.NoError(t, err)

	clientIdentity.GetMSPIDReturns("Org2MSP", nil)
	clientIdentity.GetX509CertificateReturns(A1SystemOwnerCert(), nil)
	err = blossomCC.RequestAccount(transactionContext, model.Agency{
		Name:  "test-agency",
		ATO:   "test-ato",
		MSPID: "TestMSP",
		Users: model.Users{
			SystemOwner:           "system_owner",
			AcquisitionSpecialist: "acq_spec",
			SystemAdministrator:   "system_admin",
		},
		Status:   "test-status",
		Licenses: nil,
	})
	require.NoError(t, err)

	chaincodeStub.GetStateCalls(func(s string) ([]byte, error) {
		switch s {
		case "graph":
			return []byte{1, 2, 3}, nil
		case "test":
			return []byte{9, 9, 9}, nil
		}
		return nil, nil
	})
	b, err := chaincodeStub.GetState("graph")
	require.NoError(t, err)
	fmt.Println(b)
}

func TestName(t *testing.T) {
	mockStub := shim.N
}
