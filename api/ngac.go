package api

/*import (
	"encoding/json"
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	adminpdp "github.com/usnistgov/blossom/chaincode/ngac/pdp"
)

// UpdateGraph updates the ledger graph with the graph provided.  The requesting user needs to have permission to
// make each change.
func (b *BlossomSmartContract) UpdateGraph(ctx contractapi.TransactionContextInterface, jsonStr string) error {
	pdp := new(adminpdp.AdministrativePDP)

	ledgerGraph, err := pdp.GetGraph(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting graph to update")
	}

	// unmarshal the graph json
	jsonGraph := memory.NewGraph()
	if err := jsonGraph.UnmarshalJSON([]byte(jsonStr)); err != nil {
		return fmt.Errorf("error unmarshaling provided graph json: %v", err)
	}

	// update the graph using the pdp to check for permissions
	if err := pdp.UpdateGraph(ctx, ledgerGraph, jsonGraph); err != nil {
		return fmt.Errorf("error updating graph: %v", err)
	}

	// marshal the graph to json
	var bytes []byte
	if bytes, err = ledgerGraph.MarshalJSON(); err != nil {
		return fmt.Errorf("error marshaling ledger graph after update")
	}

	// store the updated graph
	return ctx.GetStub().PutState("graph", bytes)
}

// GetGraph retrieves an instance of ngac from the world state
func (b *BlossomSmartContract) GetGraph(ctx contractapi.TransactionContextInterface) (string, error) {
	pdp := new(adminpdp.AdministrativePDP)

	graph, err := pdp.GetGraph(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting NGAC graph")
	}

	bytes, err := json.Marshal(graph)
	if err != nil {
		return "", errors.Wrap(err, "error marshaling graph")
	}

	return string(bytes), nil
}
*/
