package ngac

import (
	"fmt"

	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func getUser(ctx contractapi.TransactionContextInterface) (string, error) {
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

func getGraph(ctx contractapi.TransactionContextInterface) (pip.Graph, error) {
	// invoke the ngac chaincode to get the ngac graph
	// leaving the channel empty assumes the same channel
	response := ctx.GetStub().InvokeChaincode("n", [][]byte{[]byte("GetGraph")}, "")
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

func updateGraph(ctx contractapi.TransactionContextInterface, graph []byte) error {
	response := ctx.GetStub().InvokeChaincode("n", [][]byte{[]byte("UpdateGraph"), graph}, "mychannel")
	if response.Status != 200 {
		return fmt.Errorf("error invoking ngac chaincode: %v", response.Message)
	}

	return nil
}

func CreateAsset(ctx contractapi.TransactionContextInterface, assetID string) (func() error, error) {
	var (
		user  string
		graph pip.Graph
		err   error
	)

	if user, err = getUser(ctx); err != nil {
		return nil, err
	}

	if graph, err = getGraph(ctx); err != nil {
		return nil, err
	}

	// decide if the requesting user is allowed to create an asset
	decider := pdp.NewDecider(graph)
	if ok, err := decider.Decide(user, "assets", "CreateAsset"); err != nil {
		return nil, fmt.Errorf("error deciding if user can create an asset: %v", err)
	} else if !ok {
		return nil, fmt.Errorf("user %q does not have permssion to create assets", user)
	}

	// the commit func will be returned and allows the asset to be created in the ngac graph after it has been
	// added to the ledger
	commit := func() error {
		// create the object node representing the asset
		if err := graph.CreateNode(assetID, pip.Object, nil); err != nil {
			return err
		}

		// assign the new asset to the assets container
		if err := graph.Assign(assetID, "assets"); err != nil {
			return err
		}

		bytes, err := graph.MarshalJSON()
		if err != nil {
			return err
		}

		return updateGraph(ctx, bytes)
	}

	return commit, err
}

func ReadAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	var (
		user  string
		graph pip.Graph
		err   error
	)

	if user, err = getUser(ctx); err != nil {
		return err
	}

	if graph, err = getGraph(ctx); err != nil {
		return err
	}

	decider := pdp.NewDecider(graph)
	if ok, err := decider.Decide(user, assetID, "ReadAsset"); err != nil {
		return fmt.Errorf("error deciding if user can read asset: %v", err)
	} else if !ok {
		return fmt.Errorf("asset %q not found", assetID)
	}

	return nil
}

func UpdateAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	var (
		user  string
		graph pip.Graph
		err   error
	)

	if user, err = getUser(ctx); err != nil {
		return err
	}

	if graph, err = getGraph(ctx); err != nil {
		return err
	}

	decider := pdp.NewDecider(graph)
	if ok, err := decider.Decide(user, assetID, "UpdateAsset"); err != nil {
		return fmt.Errorf("error deciding if user can update asset: %v", err)
	} else if !ok {
		return fmt.Errorf("asset %q could not be updated", assetID)
	}

	return nil
}

func DeleteAsset(ctx contractapi.TransactionContextInterface, assetID string) (func() error, error) {
	var (
		user  string
		graph pip.Graph
		err   error
	)

	if user, err = getUser(ctx); err != nil {
		return nil, err
	}

	if graph, err = getGraph(ctx); err != nil {
		return nil, err
	}

	decider := pdp.NewDecider(graph)
	if ok, err := decider.Decide(user, assetID, "DeleteAsset"); err != nil {
		return nil, fmt.Errorf("error deciding if user can delete asset: %v", err)
	} else if !ok {
		return nil, fmt.Errorf("asset %q could not be deleted", assetID)
	}

	// the commit func will be returned and allows the asset to be deleted from the ngac graph after it has been
	// deleted from the ledger
	commit := func() error {
		// create the object node representing the asset
		if err := graph.CreateNode(assetID, pip.Object, nil); err != nil {
			return err
		}

		// assign the new asset to the assets container
		if err := graph.Assign(assetID, "assets"); err != nil {
			return err
		}

		bytes, err := graph.MarshalJSON()
		if err != nil {
			return err
		}

		return updateGraph(ctx, bytes)
	}

	return commit, nil
}
