package ngac

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func GetUser(ctx contractapi.TransactionContextInterface) (string, error) {
	var (
		cID   string
		mspID string
		err   error
	)

	// get the client and msp ids from the request to formulate user id
	if cID, err = ctx.GetClientIdentity().GetID(); err != nil {
		return "", fmt.Errorf("error retrieving client ID from request: %v", err)
	}

	if mspID, err = ctx.GetClientIdentity().GetMSPID(); err != nil {
		return "", fmt.Errorf("error retrieving MSP ID from request: %v", err)
	}

	return fmt.Sprintf("%s:%s", cID, mspID), nil
}

func GetGraph(ctx contractapi.TransactionContextInterface) (pip.Graph, error) {
	// invoke the ngac chaincode to get the ngac graph
	response := ctx.GetStub().InvokeChaincode("ngac", [][]byte{[]byte("GetGraph")}, "mychannel")
	if response.Status != 200 {
		return nil, fmt.Errorf("error invoking ngac chaincode: %v", response.Message)
	}

	// unmarshal the graph returned from the ngac chaincode
	g := memory.NewGraph()
	if err := g.UnmarshalJSON(response.GetPayload()); err != nil {
		return nil, fmt.Errorf("error unmarshaling graph json: %v", err)
	}

	return g, nil
}

func UpdateGraph(ctx contractapi.TransactionContextInterface, graph pip.Graph) error {
	// convert the graph to a byte array
	bytes, err := graph.MarshalJSON()
	if err != nil {
		return err
	}

	response := ctx.GetStub().InvokeChaincode("ngac", [][]byte{[]byte("UpdateGraph"), bytes}, "mychannel")
	if response.Status != 200 {
		return fmt.Errorf("error invoking ngac chaincode: %v", response.Message)
	}

	return nil
}
