package api

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
)

type BlossomSmartContract struct {
	contractapi.Contract
}

func (b *BlossomSmartContract) InitBlossom(ctx contractapi.TransactionContextInterface) error {
	adminPDP := pdp.NewAdminDecider()
	return adminPDP.InitGraph(ctx)
}

func (b *BlossomSmartContract) Test(ctx contractapi.TransactionContextInterface) ([]string, error) {
	s := make([]string, 0)

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		var queryResponse *queryresult.KV
		if queryResponse, err = resultsIterator.Next(); err != nil {
			return nil, err
		}

		s = append(s, string(queryResponse.Value))
	}

	return s, nil
}
