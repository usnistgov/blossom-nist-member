package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
)

type BlossomSmartContract struct {
}

func (t *BlossomSmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	adminPDP := pdp.NewAdminDecider()
	if err := adminPDP.InitGraph(stub); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *BlossomSmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}
