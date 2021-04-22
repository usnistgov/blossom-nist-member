package ledger

import (
	"asset/agency"
	"asset/license"
	"asset/ngac"
	"asset/operations"
	"asset/swid"
	"encoding/json"
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract for managing the ledger
type Contract struct {
	contractapi.Contract
}

const (
	AgenciesKey = "agencies"
	LicensesKey = "licenses"
	SwidsKey    = "swids"
)

// InitLedger initializes the ledger components including: Agencies, Licenses, and SwID tags. This method also
// invokes NGAC chaincode to initialize the NGAC components.
func (c Contract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// if the ledger is already initialized do nothing
	if data, err := ctx.GetStub().GetState(AgenciesKey); err != nil {
		return err
	} else if data != nil {
		return nil
	}

	// init ngac graph
	if err := initNGAC(ctx); err != nil {
		return fmt.Errorf("error initializing NGAC: %w", err)
	}

	// init agency, license, and swid collections
	agencies := make([]agency.Agency, 0)
	agenciesBytes, err := json.Marshal(agencies)
	if err != nil {
		return fmt.Errorf("error marshaling agency array: %w", err)
	}
	if err = ctx.GetStub().PutState(AgenciesKey, agenciesBytes); err != nil {
		return fmt.Errorf("error initializing agency collection on ledger")
	}

	licenses := make([]license.License, 0)
	licensesBytes, err := json.Marshal(licenses)
	if err != nil {
		return fmt.Errorf("error marshaling license array: %w", err)
	}
	if err = ctx.GetStub().PutState(LicensesKey, licensesBytes); err != nil {
		return fmt.Errorf("error initializing license collection on ledger")
	}

	swids := make([]swid.SwID, 0)
	swidsBytes, err := json.Marshal(swids)
	if err != nil {
		return fmt.Errorf("error marshaling swid array: %w", err)
	}
	if err = ctx.GetStub().PutState(SwidsKey, swidsBytes); err != nil {
		return fmt.Errorf("error initializing swid collection on ledger")
	}

	return nil
}

func initNGAC(ctx contractapi.TransactionContextInterface) error {
	graph, err := initGraph()
	if err != nil {
		return fmt.Errorf("error initializing NGAC graph on ledger: %w", err)
	}

	if err = ngac.UpdateGraph(ctx, graph); err != nil {
		return fmt.Errorf("error updating NGAC graph: %w", err)
	}

	return nil
}

func initGraph() (pip.Graph, error) {
	graph := memory.NewGraph()

	// create the admin user
	a0Admin, err := graph.CreateNode("A0admin", pip.User, nil)
	if err != nil {
		return nil, err
	}

	// create a UA for the admin user
	// this is the node that will be used to set the admin user's policies
	a0AdminUA, err := graph.CreateNode(a0Admin.Name+"_UA", pip.UserAttribute, nil)
	graph.Assign(a0Admin.Name, a0AdminUA.Name)

	// create RBAC policy class node
	rbacPC, err := graph.CreateNode("RBAC", pip.PolicyClass, nil)

	// create default attributes
	// these are used when a user wants to create a new attribute in the policy class
	// we can't check if the user has permissions to create a new node in a policy class
	// we can check if they can create a new node in an already existing node
	rbacUA, err := graph.CreateNode("RBAC_UA", pip.UserAttribute, nil)
	graph.Assign(rbacUA.Name, rbacPC.Name)
	rbacOA, err := graph.CreateNode("RBAC_OA", pip.ObjectAttribute, nil)
	graph.Assign(rbacOA.Name, rbacPC.Name)

	// associate the admin UA with the default attributes, giving them * permissions on all nodes in the policy class
	graph.Associate(a0AdminUA.Name, rbacUA.Name, pip.ToOps(pip.AllOps))
	graph.Associate(a0AdminUA.Name, rbacOA.Name, pip.ToOps(pip.AllOps))

	agenciesOA, err := graph.CreateNode("agencies", pip.ObjectAttribute, nil)
	graph.Assign(agenciesOA.Name, rbacOA.Name)
	licensesOA, err := graph.CreateNode("licenses", pip.ObjectAttribute, nil)
	graph.Assign(licensesOA.Name, rbacOA.Name)

	systemOwnersUA, err := graph.CreateNode("SystemOwners", pip.UserAttribute, nil)
	systemAdminsUA, err := graph.CreateNode("SystemAdmins", pip.UserAttribute, nil)
	acqSpecUA, err := graph.CreateNode("AcquisitionSpecialists", pip.UserAttribute, nil)
	graph.Assign(systemOwnersUA.Name, rbacUA.Name)
	graph.Assign(systemAdminsUA.Name, rbacUA.Name)
	graph.Assign(acqSpecUA.Name, rbacUA.Name)

	// system owners can view "agencies"
	graph.Associate(systemOwnersUA.Name, agenciesOA.Name, pip.ToOps(operations.ViewAgency))
	// system admins can read, assign, deassign (assign and deassign for the license keys) "licenses"
	graph.Associate(systemAdminsUA.Name, licensesOA.Name,
		pip.ToOps(operations.ViewLicense, operations.CheckOutLicense, operations.CheckInLicense))

	// create DAC policy class node
	graph.CreateNode("DAC", pip.PolicyClass, nil)

	// same default nodes as RBAC
	graph.CreateNode("DAC_UA", pip.UserAttribute, nil)
	graph.Assign("DAC_UA", "DAC")
	graph.CreateNode("DAC_OA", pip.ObjectAttribute, nil)
	graph.Assign("DAC_OA", "DAC")

	graph.Associate(a0AdminUA.Name, "DAC_UA", pip.ToOps(pip.AllOps))
	graph.Associate(a0AdminUA.Name, "DAC_OA", pip.ToOps(pip.AllOps))

	return graph, nil
}
