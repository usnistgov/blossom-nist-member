package main

import (
	"encoding/json"
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// AssetContract contract for managing CRUD for Asset
type AssetContract struct {
	contractapi.Contract
}

const GraphKey = "graph"

// AssetExists returns true when asset with given ID exists in world state
func (c *AssetContract) assetExists(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {
	data, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateAsset creates a new instance of Asset
func (c *AssetContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetID string, value string) error {
	// first check that the asset does not exist yet
	if exists, err := c.assetExists(ctx, assetID); err != nil {
		return fmt.Errorf("could not read from world state. %w", err)
	} else if exists {
		return fmt.Errorf("asset %s already exists", assetID)
	}

	asset := new(Asset)
	asset.Value = value

	var (
		bytes []byte
		err   error
	)

	// serialize byte array
	if bytes, err = json.Marshal(asset); err != nil {
		return fmt.Errorf("error marshaling asset %v: %w", assetID, err)
	}

	// create asset on ledger
	if err = ctx.GetStub().PutState(assetID, bytes); err != nil {
		return fmt.Errorf("error updating world state: %w", err)
	}

	// get the ngac graph
	if bytes, err = ctx.GetStub().GetState(GraphKey); err != nil {
		return fmt.Errorf("error getting graph state: %v", err)
	}

	// deserialize the ledger graph
	graph := memory.NewGraph()
	if err = graph.UnmarshalJSON(bytes); err != nil {
		return fmt.Errorf("could not unmarshal world state data to type PolicyMachine: %v", err)
	}

	// create asset in ngac
	// create the object node representing the asset
	if _, err := graph.CreateNode(assetID, pip.Object, nil); err != nil {
		return err
	}

	// assign the new asset to the assets container
	if err := graph.Assign(assetID, "assets"); err != nil {
		return err
	}

	// serialize the updated graph
	if bytes, err = graph.MarshalJSON(); err != nil {
		return err
	}

	// update the graph on the ledger
	if err = ctx.GetStub().PutState(GraphKey, bytes); err != nil {
		return fmt.Errorf("error updating graph world state: %w", err)
	}

	return nil
}

// ReadAsset retrieves an instance of Asset from the world state
func (c *AssetContract) ReadAsset(ctx contractapi.TransactionContextInterface, assetID string) (*Asset, error) {
	// check that the asset exists
	if exists, err := c.assetExists(ctx, assetID); err != nil {
		return nil, fmt.Errorf("could not read from world state: %w", err)
	} else if !exists {
		return nil, fmt.Errorf("1asset %s not found", assetID)
	}

	bytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return nil, fmt.Errorf("error getting asset ID from world state: %w", err)
	}

	asset := new(Asset)

	if err = json.Unmarshal(bytes, asset); err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Asset")
	}

	return asset, nil
}

// UpdateAsset retrieves an instance of Asset from the world state and updates its value
func (c *AssetContract) UpdateAsset(ctx contractapi.TransactionContextInterface, assetID string, newValue string) error {
	exists, err := c.assetExists(ctx, assetID)
	if err != nil {
		return fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("asset %s does not exist", assetID)
	}

	asset := new(Asset)
	asset.Value = newValue

	bytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("error marshaling asset %v: %w", assetID, err)
	}

	return ctx.GetStub().PutState(assetID, bytes)
}

// DeleteAsset deletes an instance of Asset from the world state
func (c *AssetContract) DeleteAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	if exists, err := c.assetExists(ctx, assetID); err != nil {
		return fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("asset %s does not exist", assetID)
	}

	if err := ctx.GetStub().DelState(assetID); err != nil {
		return fmt.Errorf("error deleting asset from wolrd state %v: %w", assetID, err)
	}

	var (
		bytes []byte
		err   error
	)

	// commit the changes to the ngac graph
	// get the ngac graph
	if bytes, err = ctx.GetStub().GetState(GraphKey); err != nil {
		return fmt.Errorf("error getting graph state: %v", err)
	}

	// deserialize the ledger graph
	graph := memory.NewGraph()
	if err = graph.UnmarshalJSON(bytes); err != nil {
		return fmt.Errorf("could not unmarshal world state data to type PolicyMachine: %v", err)
	}

	// create asset in ngac
	// create the object node representing the asset
	if err := graph.DeleteNode(assetID); err != nil {
		return err
	}

	// serialize the updated graph
	if bytes, err = graph.MarshalJSON(); err != nil {
		return err
	}

	// update the graph on the ledger
	if err = ctx.GetStub().PutState(GraphKey, bytes); err != nil {
		return fmt.Errorf("error updating graph world state: %w", err)
	}

	return nil
}
