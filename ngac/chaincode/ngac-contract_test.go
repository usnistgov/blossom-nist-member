package chaincode

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/pdp"
	pdp2 "github.com/usnistgov/blossom/chaincode/ngac/pdp"
	"testing"

	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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

func TestPolicy(t *testing.T) {
	graph := memory.NewGraph()

	// create the admin user
	a0Admin, err := graph.CreateNode("Org1 Admin:Org1MSP", pip.User, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating admin user node: %w", err)
	}

	// create a UA for the admin user
	// this is the node that will be used to set the admin user's policies
	a0AdminUA, err := graph.CreateNode(a0Admin.Name+"_UA", pip.UserAttribute, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating admin user attribute: %w", err)
	}

	if err = graph.Assign(a0Admin.Name, a0AdminUA.Name); err != nil {
		// return nil, fmt.Errorf("error assigning %q to %q: %w", a0Admin.Name, a0AdminUA.Name, err)
	}

	// create RBAC policy class node
	rbacPC, err := graph.CreateNode("RBAC", pip.PolicyClass, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating RBAC policy class: %w", err)
	}

	// create default attributes
	// these are used when a user wants to create a new attribute in the policy class
	// we can't check if the user has permissions to create a new node in a policy class
	// we can check if they can create a new node in an already existing node
	rbacUA, err := graph.CreateNode("RBAC_UA", pip.UserAttribute, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating RBAC user attribute: %w", err)
	}

	if err = graph.Assign(rbacUA.Name, rbacPC.Name); err != nil {
		// return nil, fmt.Errorf("error assigning %q to %q: %w", rbacUA.Name, rbacPC.Name, err)
	}

	rbacOA, err := graph.CreateNode("RBAC_OA", pip.ObjectAttribute, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating RBAC object attribute: %w", err)
	}

	if err = graph.Assign(rbacOA.Name, rbacPC.Name); err != nil {
		// return nil, fmt.Errorf("error assigning %q to %q: %w", rbacOA.Name, rbacPC.Name, err)
	}

	// create a UA to hold each agency UA
	agenciesUA, err := graph.CreateNode("agencies_ua", pip.UserAttribute, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating agencies base user attribute: %w", err)
	}

	if err = graph.Assign(agenciesUA.Name, rbacUA.Name); err != nil {
		// return nil, fmt.Errorf("error assigning %q to %q: %w", agenciesUA.Name, rbacUA.Name, err)
	}

	// associate the admin UA with the default attributes, giving them * permissions on all nodes in the policy class
	if err = graph.Associate(a0AdminUA.Name, rbacUA.Name, pip.ToOps(pip.AllOps)); err != nil {
		// return nil, fmt.Errorf("error associating %q with %q: %w", a0AdminUA.Name, rbacUA.Name, err)
	}
	if err = graph.Associate(a0AdminUA.Name, rbacOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		// return nil, fmt.Errorf("error associating %q with %q: %w", a0AdminUA.Name, rbacOA.Name, err)
	}

	agenciesOA, err := graph.CreateNode("agencies", pip.ObjectAttribute, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating agencies base object attribute: %w", err)
	}

	if err = graph.Assign(agenciesOA.Name, rbacOA.Name); err != nil {
		// return nil, fmt.Errorf("error assigning %q to %q: %w", agenciesOA.Name, rbacOA.Name, err)
	}

	licensesOA, err := graph.CreateNode("licenses", pip.ObjectAttribute, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating licenses base object attribute: %w", err)
	}

	if err = graph.Assign(licensesOA.Name, rbacOA.Name); err != nil {
		// return nil, fmt.Errorf("error assigning %q to %q: %w", licensesOA.Name, rbacOA.Name, err)
	}

	systemOwnersUA, err := graph.CreateNode("SystemOwners", pip.UserAttribute, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating SystemOwners user attribute: %w", err)
	}

	systemAdminsUA, err := graph.CreateNode("SystemAdmins", pip.UserAttribute, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating SystemAdmins user attribute: %w", err)
	}

	acqSpecUA, err := graph.CreateNode("AcquisitionSpecialists", pip.UserAttribute, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating AcquisitionSpecialists user attribute: %w", err)
	}

	if err = graph.Assign(systemOwnersUA.Name, rbacUA.Name); err != nil {
		// return nil, fmt.Errorf("error assigning %q to %q: %w", systemOwnersUA.Name, rbacUA.Name, err)
	}

	if err = graph.Assign(systemAdminsUA.Name, rbacUA.Name); err != nil {
		// return nil, fmt.Errorf("error assigning %q to %q: %w", systemAdminsUA.Name, rbacUA.Name, err)
	}

	if err = graph.Assign(acqSpecUA.Name, rbacUA.Name); err != nil {
		// return nil, fmt.Errorf("error assigning %q to %q: %w", acqSpecUA.Name, rbacUA.Name, err)
	}

	// system owners can view "agencies"
	if err = graph.Associate(systemOwnersUA.Name, agenciesOA.Name, pip.ToOps(ViewAgency)); err != nil {
		// return nil, fmt.Errorf("error associating %q with %q: %w", systemOwnersUA.Name, agenciesOA.Name, err)
	}
	// system admins can read, assign, deassign (assign and deassign for the license keys) "licenses"
	if err = graph.Associate(systemAdminsUA.Name, licensesOA.Name,
		pip.ToOps(ViewLicense, CheckOutLicense, CheckInLicense)); err != nil {
		// return nil, fmt.Errorf("error associating %q with %q: %w", systemAdminsUA.Name, licensesOA.Name, err)
	}
	// acquisition specialists can audit agency licenses
	if err = graph.Associate(acqSpecUA.Name, licensesOA.Name, pip.ToOps(ViewLicense, ViewAgency)); err != nil {
		// return nil, fmt.Errorf("error associating %q with %q: %w", acqSpecUA.Name, licensesOA.Name, err)
	}

	// create DAC policy class node
	dacPC, err := graph.CreateNode("DAC", pip.PolicyClass, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating DAC policy class node: %w", err)
	}

	// same default nodes as RBAC
	dacUA, err := graph.CreateNode("DAC_UA", pip.UserAttribute, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating DAC user attribute node: %w", err)
	}

	if err = graph.Assign(dacUA.Name, dacPC.Name); err != nil {
		// return nil, fmt.Errorf("error assigning %q to %q: %w", dacUA.Name, dacPC.Name, err)
	}

	dacOA, err := graph.CreateNode("DAC_OA", pip.ObjectAttribute, nil)
	if err != nil {
		// return nil, fmt.Errorf("error creating DAC object attribute node: %w", err)
	}

	if err = graph.Assign(dacOA.Name, dacPC.Name); err != nil {
		// return nil, fmt.Errorf("error assigning %q to %q: %w", dacUA.Name, dacPC.Name, err)
	}

	if err = graph.Associate(a0AdminUA.Name, dacUA.Name, pip.ToOps(pip.AllOps)); err != nil {
		// return nil, fmt.Errorf("error associating %q with %q: %w", a0AdminUA.Name, dacUA.Name, err)
	}

	if err = graph.Associate(a0AdminUA.Name, dacOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		// return nil, fmt.Errorf("error associating %q with %q: %w", a0AdminUA.Name, dacOA.Name, err)
	}

	d := pdp.NewDecider(graph)
	fmt.Println(d.Decide(a0Admin.Name, "licenses", pdp2.CreateNodePermission))
}
