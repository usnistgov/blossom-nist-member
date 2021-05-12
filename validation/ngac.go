package validation

import (
	"fmt"

	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric-protos-go/peer"
	. "github.com/hyperledger/fabric/core/handlers/validation/api/identities"
	validation "github.com/hyperledger/fabric/core/handlers/validation/api/state"
)

type (
	// canInvokeFunc defines a function that check if a given user can invoke the function in the invokeSpec parameter.
	// The NGAC decider is used to make an access decision
	canInvokeFunc func(user string, decider pdp.Decider, invokeSpec *peer.ChaincodeInvocationSpec) (bool, error)

	// NGACVSCC validates transactions using NGAC.
	NGACVSCC struct {
		sf validation.StateFetcher
	}
)

const GraphKey = "graph"

func New(sf validation.StateFetcher) *NGACVSCC {
	return &NGACVSCC{sf: sf}
}

// CanInvoke checks if the given user has the permission to invoke the chaincode function that is called in the invokeSpec.
// This function retrieves the NGAC graph from the world state and uses it to decide if the user can invoke the chaincode.
func (n NGACVSCC) CanInvoke(identity *IdentityIdentifier, namespace string, invokeSpec *peer.ChaincodeInvocationSpec) (bool, error) {
	var (
		err   error
		state validation.State
	)

	if state, err = n.sf.FetchState(); err != nil {
		return false, nil
	}

	defer state.Done()

	// construct the ngac user name which is userid:mspid
	user := fmt.Sprintf(identity.Id, ":", identity.Mspid)

	// get graph from world state
	var graphBytes []byte
	if graphBytes, err = getState(state, namespace, GraphKey); err != nil {
		return false, err
	}

	g := memory.NewGraph()
	if err := g.UnmarshalJSON(graphBytes); err != nil {
		return false, err
	}

	// get the canInvoke function for the invoked chaincode function
	var canInvoke canInvokeFunc
	canInvoke, err = getCanInvokeFunc(string(invokeSpec.ChaincodeSpec.Input.Args[1]))
	if err != nil {
		return false, err
	}

	// check if the user can perform the operation
	// return false is an error occurs or the user is denied permission
	var ok bool
	if ok, err = canInvoke(user, pdp.NewDecider(g), invokeSpec); err != nil {
		return false, err
	} else if !ok {
		return ok, nil
	}

	// commit any changes to the graph
	return true, nil
}

// getCanInvokeFunc returns the canInvokeFunc that corresponds to the provided function name
func getCanInvokeFunc(funcName string) (canInvokeFunc, error) {
	switch funcName {
	case "CreateAsset":
		return canCreateAsset, nil
	case "ReadAsset":
		return canReadAsset, nil
	case "UpdateAsset":
		return canUpdateAsset, nil
	case "DeleteAsset":
		return canDeleteAsset, nil
	}

	return nil, nil
}

func canCreateAsset(user string, decider pdp.Decider, invokeSpec *peer.ChaincodeInvocationSpec) (bool, error) {
	if ok, err := decider.HasPermissions(user, "assets", "CreateAsset"); err != nil {
		return false, fmt.Errorf("error deciding if user can create an asset: %v", err)
	} else if !ok {
		return ok, fmt.Errorf("user %q does not have permssion to create assets", user)
	}

	return true, nil
}

func canReadAsset(user string, decider pdp.Decider, invokeSpec *peer.ChaincodeInvocationSpec) (bool, error) {
	return true, nil
}

func canUpdateAsset(user string, decider pdp.Decider, invokeSpec *peer.ChaincodeInvocationSpec) (bool, error) {
	return true, nil
}

func canDeleteAsset(user string, decider pdp.Decider, invokeSpec *peer.ChaincodeInvocationSpec) (bool, error) {
	return true, nil
}

// getState retrieves the value for the given key in the given namespace
func getState(s validation.State, namespace string, key string) ([]byte, error) {
	values, err := s.GetStateMultipleKeys(namespace, []string{key})
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}
	return values[0], nil
}
