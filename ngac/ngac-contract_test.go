/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/mock"
	"testing"
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
	testGraph := memory.NewGraph()
	testGraph.CreateNode("pc1", pip.PolicyClass, nil)
	testGraphBytes, _ := testGraph.MarshalJSON()

	ms := new(MockStub)
	ms.On("GetState", "graph").Return(testGraphBytes, nil)
	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

	mc := new(MockContext)
	mc.On("GetStub").Return(ms)

	return mc, ms
}

func TestUpdateGraph(t *testing.T) {
	pm := new(NGACContract)
	ctx, _ := configureStub()
	pm.InitNGAC(ctx)

	expected := memory.NewGraph()
	expected.CreateNode("pc2", pip.PolicyClass, map[string]string{"k": "v"})
	expectedBytes, _ := expected.MarshalJSON()
	pm.UpdateGraph(ctx, string(expectedBytes))

	actualStr, _ := pm.GetGraph(ctx)
	actual := memory.NewGraph()
	actual.UnmarshalJSON([]byte(actualStr))
	if _, err := actual.GetNode("pc2"); err != nil {
		t.Fatal("graph was not updated")
	}
}

func TestPolicy(t *testing.T) {
	graph := memory.NewGraph()

	// create the initial configuration
	if err := graph.CreateNode("assets_pc", pip.PolicyClass, nil); err != nil {
		t.Fatal(err)
	}
	if err := graph.CreateNode("assets", pip.ObjectAttribute, nil); err != nil {
		t.Fatal(err)
	}
	if err := graph.CreateNode("reader", pip.UserAttribute, nil); err != nil {
		t.Fatal(err)
	}
	if err := graph.CreateNode("writer", pip.UserAttribute, nil); err != nil {
		t.Fatal(err)
	}
	if err := graph.CreateNode("creator", pip.UserAttribute, nil); err != nil {
		t.Fatal(err)
	}

	// testUser1
	testUser1 := "eDUwOTo6Q049dGVzdFVzZXIxLE9VPWNsaWVudDo6Q049T3JnMSBDQQ==:Org1MSP"
	if err := graph.CreateNode(testUser1, pip.User, nil); err != nil {
		t.Fatal(err)
	}
	if err := graph.Assign(testUser1, "reader"); err != nil {
		t.Fatal(err)
	}

	// org1 admin
	org1Admin := "eDUwOTo6Q049T3JnMSBBZG1pbixPVT1hZG1pbjo6Q049T3JnMSBDQQ==:Org1MSP"
	if err := graph.CreateNode(org1Admin, pip.User, nil); err != nil {
		t.Fatal(err)
	}
	if err := graph.Assign(org1Admin, "creator"); err != nil {
		t.Fatal(err)
	}

	if err := graph.Assign("creator", "writer"); err != nil {
		t.Fatal(err)
	}

	if err := graph.Assign("writer", "reader"); err != nil {
		t.Fatal(err)
	}

	if err := graph.Assign("assets", "assets_pc"); err != nil {
		t.Fatal(err)
	}

	if err := graph.Associate("creator", "assets", pip.ToOps("CreateAsset")); err != nil {
		t.Fatal(err)
	}
	if err := graph.Associate("writer", "assets", pip.ToOps("UpdateAsset")); err != nil {
		t.Fatal(err)
	}
	if err := graph.Associate("reader", "assets", pip.ToOps("ReadAsset")); err != nil {
		t.Fatal(err)
	}

	if err := graph.CreateNode("o1", pip.Object, nil); err != nil {
		t.Fatal(err)
	}
	if err := graph.Assign("o1", "assets"); err != nil {
		t.Fatal(err)
	}

	decider := pdp.NewDecider(graph)
	fmt.Println(decider.Decide(org1Admin, "o1", "ReadAsset"))
}
