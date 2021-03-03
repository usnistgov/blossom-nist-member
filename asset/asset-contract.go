/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"asset/ngac"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// AssetContract contract for managing CRUD for Asset
type AssetContract struct {
	contractapi.Contract
}

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
	exists, err := c.assetExists(ctx, assetID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("The asset %s already exists", assetID)
	}

	// check that the requesting user has permissions to create an asset
	// the function returned will be used to commit the new asset to the
	// ngac graph once it has been successfully added to the ledger
	ngacCommit, err := ngac.CreateAsset(ctx, assetID)
	if err != nil {
		return err
	}

	asset := new(Asset)
	asset.Value = value

	bytes, _ := json.Marshal(asset)

	if err = ctx.GetStub().PutState(assetID, bytes); err != nil {
		return err
	}

	// commit the new asset to the ngac graph
	return ngacCommit()
}

// ReadAsset retrieves an instance of Asset from the world state
func (c *AssetContract) ReadAsset(ctx contractapi.TransactionContextInterface, assetID string) (*Asset, error) {
	// check that the asset exists
	exists, err := c.assetExists(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("1asset %s not found", assetID)
	}

	// check the user can read the asset
	if err = ngac.ReadAsset(ctx, assetID); err != nil {
		return nil, err
	}

	bytes, _ := ctx.GetStub().GetState(assetID)

	asset := new(Asset)

	err = json.Unmarshal(bytes, asset)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Asset")
	}

	return asset, nil
}

// UpdateAsset retrieves an instance of Asset from the world state and updates its value
func (c *AssetContract) UpdateAsset(ctx contractapi.TransactionContextInterface, assetID string, newValue string) error {
	exists, err := c.assetExists(ctx, assetID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", assetID)
	}

	// check the user can update the asset
	if err = ngac.UpdateAsset(ctx, assetID); err != nil {
		return err
	}

	asset := new(Asset)
	asset.Value = newValue

	bytes, _ := json.Marshal(asset)

	return ctx.GetStub().PutState(assetID, bytes)
}

// DeleteAsset deletes an instance of Asset from the world state
func (c *AssetContract) DeleteAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	if exists, err := c.assetExists(ctx, assetID); err != nil {
		return fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("the asset %s does not exist", assetID)
	}

	// check the user can read the asset
	ngacCommit, err := ngac.DeleteAsset(ctx, assetID)
	if err != nil {
		return err
	}

	if err := ctx.GetStub().DelState(assetID); err != nil {
		return err
	}

	// commit the changes to the ngac graph
	return ngacCommit()
}
