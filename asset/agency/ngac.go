package agency

import (
	"asset/ngac"
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func createAgency(ctx contractapi.TransactionContextInterface, agency Agency) error {
	// add agency to agencies attribute in NGAC
	graph, err := ngac.GetGraph(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving ngac graph from ")
	}

	// create an object to represent the agency
	if _, err = graph.CreateNode(agency.Name, pip.Object, nil); err != nil {
		return fmt.Errorf("error creating agency in NGAC: %w", err)
	}

	// assign the agency object to the agencies attribute
	if err = graph.Assign(agency.Name, "agencies"); err != nil {
		return fmt.Errorf("error assigning agency %q to agencies attribute: %w", agency.Name, err)
	}

	return ngac.UpdateGraph(ctx, graph)
}
