package api

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
)

type BlossomSmartContract struct {
	contractapi.Contract
}

func (b *BlossomSmartContract) InitBlossom(ctx contractapi.TransactionContextInterface) error {
	adminPDP := pdp.NewAdminDecider()
	return adminPDP.InitGraph(ctx)
}
