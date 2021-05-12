package ledger

import (
	"encoding/json"

	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
)

func GetGraph(ctx contractapi.TransactionContextInterface) (pip.Graph, error) {
	bytes, err := ctx.GetStub().GetState("graph")
	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return nil, errors.Errorf("NGAC graph has not been initialized")
	}

	graph := memory.NewGraph()
	if err = json.Unmarshal(bytes, graph); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling graph")
	}

	return graph, nil
}

func GetGraphBytes(ctx contractapi.TransactionContextInterface) ([]byte, error) {
	bytes, err := ctx.GetStub().GetState("graph")
	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return nil, errors.Errorf("NGAC graph has not been initialized")
	}

	return bytes, nil
}

func UpdateGraphState(ctx contractapi.TransactionContextInterface, graph pip.Graph) error {
	bytes, err := json.Marshal(graph)
	if err != nil {
		return errors.Wrapf(err, "error serializing graph")
	}

	if err = ctx.GetStub().PutState("graph", bytes); err != nil {
		return errors.Wrapf(err, "error updating graph state")
	}

	return nil
}
