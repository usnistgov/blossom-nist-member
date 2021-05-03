package main

import (
	"fmt"

	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
)

// NGACContract contract for managing ngac
type NGACContract struct {
	contractapi.Contract
}

func main() {}

func (c *NGACContract) graphExists(ctx contractapi.TransactionContextInterface) (bool, error) {
	data, err := ctx.GetStub().GetState("graph")

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// UpdateGraph updates the ledger graph with the graph provided.  The requesting user needs to have permission to
// make each change.
func (c *NGACContract) UpdateGraph(ctx contractapi.TransactionContextInterface, jsonStr string) error {
	bytes, err := ctx.GetStub().GetState("graph")
	if err != nil {
		return fmt.Errorf("error getting graph state: %v", err)
	} else if bytes == nil {
		return fmt.Errorf("policy machine has not been initialized with InitPolicyMachine")
	}
	fmt.Println(string(bytes))

	// unmarshal the ledger pm
	ledgerGraph := memory.NewGraph()
	if err = ledgerGraph.UnmarshalJSON(bytes); err != nil {
		return fmt.Errorf("could not unmarshal world state data to type PolicyMachine: %v", err)
	}

	// unmarshal the graph json
	jsonGraph := memory.NewGraph()
	if err = jsonGraph.UnmarshalJSON([]byte(jsonStr)); err != nil {
		return fmt.Errorf("error unmarshaling provided graph json: %v", err)
	}

	// update the graph
	pdp := new(pdp.PDP)
	if err = pdp.UpdateGraph(ctx, ledgerGraph, jsonGraph); err != nil {
		return fmt.Errorf("error updating graph: %v", err)
	}

	// marshal the graph to json
	if bytes, err = ledgerGraph.MarshalJSON(); err != nil {
		return fmt.Errorf("error marshaling ledger graph after update")
	}

	// store the updated graph
	return ctx.GetStub().PutState("graph", bytes)
}

// GetGraph retrieves an instance of ngac from the world state
func (c *NGACContract) GetGraph(ctx contractapi.TransactionContextInterface) (string, error) {
	exists, err := c.graphExists(ctx)
	if err != nil {
		return "", fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return "", fmt.Errorf("graph does not exist")
	}

	bytes, _ := ctx.GetStub().GetState("graph")

	return string(bytes), nil
}
