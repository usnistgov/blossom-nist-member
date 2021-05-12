package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/usnistgov/blossom/chaincode/api"
	"log"
)

func main() {
	assetChaincode, err := contractapi.NewChaincode(&api.BlossomSmartContract{})
	if err != nil {
		log.Panicf("error creating blossom chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("error starting blossom chaincode: %v", err)
	}
}
