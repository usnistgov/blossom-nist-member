/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"ngac/pdp"

	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// NGACContract contract for managing ngac
type NGACContract struct {
	contractapi.Contract
}

// InitNGAC initializes the ngac components on the ledger
func (c *NGACContract) InitNGAC(ctx contractapi.TransactionContextInterface) error {
	graph := memory.NewGraph()

	// create the initial configuration
	if err := graph.CreateNode("assets_pc", pip.PolicyClass, nil); err != nil {
		return err
	}
	if err := graph.CreateNode("assets", pip.ObjectAttribute, nil); err != nil {
		return err
	}
	if err := graph.CreateNode("reader", pip.UserAttribute, nil); err != nil {
		return err
	}
	if err := graph.CreateNode("writer", pip.UserAttribute, nil); err != nil {
		return err
	}
	if err := graph.CreateNode("creator", pip.UserAttribute, nil); err != nil {
		return err
	}

	// testUser1
	testUser1 := "eDUwOTo6Q049dGVzdFVzZXIxLE9VPWNsaWVudDo6Q049T3JnMSBDQQ==:Org1MSP"
	if err := graph.CreateNode(testUser1, pip.User, nil); err != nil {
		return err
	}
	if err := graph.Assign(testUser1, "reader"); err != nil {
		return err
	}

	// org1 admin
	org1Admin := "eDUwOTo6Q049T3JnMSBBZG1pbixPVT1hZG1pbjo6Q049T3JnMSBDQQ==:Org1MSP"
	if err := graph.CreateNode(org1Admin, pip.User, nil); err != nil {
		return err
	}
	if err := graph.Assign(org1Admin, "creator"); err != nil {
		return err
	}

	if err := graph.Assign("creator", "writer"); err != nil {
		return err
	}

	if err := graph.Assign("writer", "reader"); err != nil {
		return err
	}

	if err := graph.Assign("assets", "assets_pc"); err != nil {
		return err
	}

	if err := graph.Associate("creator", "assets", pip.ToOps("CreateAsset", pdp.CreateNodePermission)); err != nil {
		return err
	}
	if err := graph.Associate("writer", "assets", pip.ToOps("UpdateAsset")); err != nil {
		return err
	}
	if err := graph.Associate("reader", "assets", pip.ToOps("ReadAsset")); err != nil {
		return err
	}

	bytes, _ := graph.MarshalJSON()
	if err := ctx.GetStub().PutState("graph", bytes); err != nil {
		return err
	}

	return nil
}

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
