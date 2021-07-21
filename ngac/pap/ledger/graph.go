package ledger

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/pkg/errors"
)

func GetGraph(stub shim.ChaincodeStubInterface) (pip.Graph, error) {
	bytes, err := stub.GetState("graph")
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

func GetGraphBytes(stub shim.ChaincodeStubInterface) ([]byte, error) {
	bytes, err := stub.GetState("graph")
	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return nil, errors.Errorf("NGAC graph has not been initialized")
	}

	return bytes, nil
}

func UpdateGraphState(stub shim.ChaincodeStubInterface, graph pip.Graph) error {
	bytes, err := json.Marshal(graph)
	if err != nil {
		return errors.Wrapf(err, "error serializing graph")
	}

	if err = stub.PutState("graph", bytes); err != nil {
		return errors.Wrapf(err, "error updating graph state")
	}

	return nil
}
