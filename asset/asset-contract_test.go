package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const getStateError = "world state get error"

type MockStub struct {
	shim.ChaincodeStubInterface
	mock.Mock
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)

	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)

	return args.Error(0)
}

func (ms *MockStub) DelState(key string) error {
	args := ms.Called(key)

	return args.Error(0)
}

type MockContext struct {
	contractapi.TransactionContextInterface
	mock.Mock
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	args := mc.Called()

	return args.Get(0).(*MockStub)
}

func configureStub() (*MockContext, *MockStub) {
	var nilBytes []byte

	testAsset := new(Asset)
	testAsset.Value = "set value"
	assetBytes, _ := json.Marshal(testAsset)

	ms := new(MockStub)
	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
	ms.On("GetState", "missingkey").Return(nilBytes, nil)
	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
	ms.On("GetState", "assetkey").Return(assetBytes, nil)
	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

	mc := new(MockContext)
	mc.On("GetStub").Return(ms)

	return mc, ms
}

func TestAssetExists(t *testing.T) {
	var exists bool
	var err error

	ctx, _ := configureStub()
	c := new(AssetContract)

	exists, err = c.assetExists(ctx, "statebad")
	assert.EqualError(t, err, getStateError)
	assert.False(t, exists, "should return false on error")

	exists, err = c.assetExists(ctx, "missingkey")
	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
	assert.False(t, exists, "should return false when no value for key in world state")

	exists, err = c.assetExists(ctx, "existingkey")
	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
	assert.True(t, exists, "should return true when value for key in world state")
}

func TestCreateAsset(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(AssetContract)

	err = c.CreateAsset(ctx, "statebad", "some value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.CreateAsset(ctx, "existingkey", "some value")
	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

	err = c.CreateAsset(ctx, "missingkey", "some value")
	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
}

func TestReadAsset(t *testing.T) {
	var asset *Asset
	var err error

	ctx, _ := configureStub()
	c := new(AssetContract)

	asset, err = c.ReadAsset(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
	assert.Nil(t, asset, "should not return Asset when exists errors when reading")

	asset, err = c.ReadAsset(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
	assert.Nil(t, asset, "should not return Asset when key does not exist in world state when reading")

	asset, err = c.ReadAsset(ctx, "existingkey")
	assert.EqualError(t, err, "Could not unmarshal world state data to type Asset", "should error when data in key is not Asset")
	assert.Nil(t, asset, "should not return Asset when data in key is not of type Asset")

	asset, err = c.ReadAsset(ctx, "assetkey")
	expectedAsset := new(Asset)
	expectedAsset.Value = "set value"
	assert.Nil(t, err, "should not return error when Asset exists in world state when reading")
	assert.Equal(t, expectedAsset, asset, "should return deserialized Asset from world state")
}

func TestUpdateAsset(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(AssetContract)

	err = c.UpdateAsset(ctx, "statebad", "new value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

	err = c.UpdateAsset(ctx, "missingkey", "new value")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when updating")

	err = c.UpdateAsset(ctx, "assetkey", "new value")
	expectedAsset := new(Asset)
	expectedAsset.Value = "new value"
	expectedAssetBytes, _ := json.Marshal(expectedAsset)
	assert.Nil(t, err, "should not return error when Asset exists in world state when updating")
	stub.AssertCalled(t, "PutState", "assetkey", expectedAssetBytes)
}

func TestDeleteAsset(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(AssetContract)

	err = c.DeleteAsset(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.DeleteAsset(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when deleting")

	err = c.DeleteAsset(ctx, "assetkey")
	assert.Nil(t, err, "should not return error when Asset exists in world state when deleting")
	stub.AssertCalled(t, "DelState", "assetkey")
}
