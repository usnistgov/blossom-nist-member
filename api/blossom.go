package api

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
)

type BlossomSmartContract struct {
	contractapi.Contract
}

func (b *BlossomSmartContract) InitBlossom(ctx contractapi.TransactionContextInterface) error {
	adminPDP, err := pdp.NewAdminDecider(ctx)
	if err != nil {
		return errors.Errorf("error initializing administrative decider")
	}

	return adminPDP.InitGraph(ctx)
}
