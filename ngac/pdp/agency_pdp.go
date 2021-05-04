package pdp

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/usnistgov/blossom/chaincode/asset/chaincode"
	"github.com/usnistgov/blossom/chaincode/asset/operations"
)

type AgencyPDP struct {
}

func (a AgencyPDP) RequestAccount(ctx contractapi.TransactionContextInterface, agencyName string) error {
	// add agency to agencies attribute in NGAC
	graph, err := GetGraph(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving ngac graph from ledger: %w", err)
	}

	// simulate an obligation in which the admin creates the user DAC configuration
	// get the user from the request -- this will be the system owner for the agency
	var user string
	if user, err = GetUser(ctx); err != nil {
		return errors.Wrap(err, "error getting user from request")
	}

	// create the user

	// create user who submitted the request as the agency system owner
	// this will allow them to view and create the agency in the graph

	// create an object to represent the agency
	var agencyObj pip.Node
	if agencyObj, err = graph.CreateNode(agency.Name+"_object", pip.Object, map[string]string{"agency": agency.Name, "type": "agency"}); err != nil {
		return fmt.Errorf("error creating agency in NGAC: %w", err)
	}

	// assign the agency object to the agencies attribute
	if err = graph.Assign(agencyObj.Name, "agencies"); err != nil {
		return fmt.Errorf("error assigning agency %q to agencies attribute: %w", agency.Name, err)
	}

	return UpdateGraph(ctx, graph)
}

func (a AgencyPDP) UploadATO(ctx contractapi.TransactionContextInterface, agencyName string) error {
	panic("implement me")
}

func (a AgencyPDP) UpdateAgencyStatus(ctx contractapi.TransactionContextInterface, agencyName string) error {
	panic("implement me")
}

func (a AgencyPDP) ApproveAccountRequest(ctx contractapi.TransactionContextInterface, agencyName string) error {
	panic("implement me")
}

func (a AgencyPDP) DenyAccountRequest(ctx contractapi.TransactionContextInterface, agencyName string) error {
	panic("implement me")
}

func (a AgencyPDP) Agencies(ctx contractapi.TransactionContextInterface) ([]string, error) {
	panic("implement me")
}

func (a AgencyPDP) Agency(ctx contractapi.TransactionContextInterface, agencyName string) error {
	panic("implement me")
}

func ngacApproveAgency(ctx contractapi.TransactionContextInterface, agency *chaincode.Agency) error {
	// create DAC attributes for agency
	graph, err := GetGraph(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving ngac graph from ledger: %w", err)
	}

	// create a user attribute
	agencyUA, err := graph.CreateNode(agency.Name+"_ua", pip.UserAttribute, nil)
	if err != nil {
		return fmt.Errorf("error creating agency user attribute: %w", err)
	}

	if err = graph.Assign(agencyUA.Name, "agencies_ua"); err != nil {
		return fmt.Errorf("error assigning agency UA to agencies_ua: %w", err)
	}

	// create users
	userNameFunc := func(name string, mspid string) string {
		return fmt.Sprintf("%s:%s", name, mspid)
	}

	soUser, err := graph.CreateNode(userNameFunc(agency.Users.SystemOwner, agency.MSPID), pip.User, nil)
	if err != nil {
		return fmt.Errorf("error creating agency system owner: %w", err)
	}

	asUser, err := graph.CreateNode(userNameFunc(agency.Users.AcquisitionSpecialist, agency.MSPID), pip.User, nil)
	if err != nil {
		return fmt.Errorf("error creating agency acquisition specialist: %w", err)
	}

	saUser, err := graph.CreateNode(userNameFunc(agency.Users.SystemAdministrator, agency.MSPID), pip.User, nil)
	if err != nil {
		return fmt.Errorf("error creating agency system administrator: %w", err)
	}

	if err = graph.Assign(soUser.Name, agencyUA.Name); err != nil {
		return fmt.Errorf("error assigning agency system owner to agency UA: %w", err)
	}

	if err = graph.Assign(asUser.Name, agencyUA.Name); err != nil {
		return fmt.Errorf("error assigning agency acquisition specialist to agency UA: %w", err)
	}

	if err = graph.Assign(saUser.Name, agencyUA.Name); err != nil {
		return fmt.Errorf("error assigning agency system administrator to agency UA: %w", err)
	}

	// create OA
	agencyOA, err := graph.CreateNode(agency.Name+"_licenses", pip.ObjectAttribute, map[string]string{"licenses": agency.Name})
	if err != nil {
		return fmt.Errorf("error creating agency licenses OA: %w", err)
	}

	if err = graph.Assign(agencyOA.Name, "DAC_OA"); err != nil {
		return fmt.Errorf("error assiging agency licenses OA to DAC_OA: %w", err)
	}

	if err = graph.Assign(agency.Name+"_object", agencyOA.Name); err != nil {
		return fmt.Errorf("error assigning agency_object to agency licenses OA")
	}

	if err = graph.Associate(agencyUA.Name, agencyOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return fmt.Errorf("error associating agency UA with agency licenses OA")
	}

	return nil
}

func ngacAgencies(ctx contractapi.TransactionContextInterface) ([]string, error) {
	graph, err := GetGraph(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ngac graph from ledger: %w", err)
	}

	user, err := GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting user from context: %w", err)
	}

	objects, err := graph.Find(pip.Object, map[string]string{"type": "agency"})
	if err != nil {
		return nil, fmt.Errorf("error getting ngac objects: %w", err)
	}

	agencies := make([]string, 0)
	for _, node := range objects {
		// do not add any agencies the user cannot view
		decider := pdp.NewDecider(graph)
		if ok, err := decider.Decide(user, node.Name, operations.ViewAgency); err != nil {
			return nil, fmt.Errorf("error checking permissions on %s: %w", node.Name, err)
		} else if !ok {
			continue
		}

		agencies = append(agencies, node.Properties["agency"])
	}

	return agencies, nil
}

func ngacAgency(ctx contractapi.TransactionContextInterface, agency string) (bool, error) {
	graph, err := GetGraph(ctx)
	if err != nil {
		return false, fmt.Errorf("error retrieving ngac graph from ledger: %w", err)
	}

	user, err := GetUser(ctx)
	if err != nil {
		return false, fmt.Errorf("error getting user from context: %w", err)
	}

	decider := pdp.NewDecider(graph)
	return decider.Decide(user, agency+"_object", operations.ViewAgency)
}
