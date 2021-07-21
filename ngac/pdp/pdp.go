package pdp

import (
	"encoding/json"
	"fmt"
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/ledger"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy"
)

var ErrAccessDenied = errors.New("access denied")

func FormatUsername(user string, mspid string) string {
	return fmt.Sprintf("%s:%s", user, mspid)
}

func GetUser(stub shim.ChaincodeStubInterface) (string, error) {
	cert, err := cid.GetX509Certificate(stub)
	if err != nil {
		return "", err
	}

	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return "", err
	}

	return FormatUsername(cert.Subject.CommonName, mspid), nil
}

// AdminDecider is the administrative Policy Decision Point (PDP) for the Blossom NGAC smart contract.
// An administrative PDP performs permission checks for all actions on an NGAC graph.  The permissions it is looking for
// are NGAC specific such as "create node" or "assign to"
type AdminDecider struct {
	user string
}

func NewAdminDecider() *AdminDecider {
	return &AdminDecider{}
}

func (a *AdminDecider) setup(stub shim.ChaincodeStubInterface) error {
	user, err := GetUser(stub)
	if err != nil {
		return errors.Wrapf(err, "error getting user from request")
	}

	a.user = user

	return nil
}

func (a *AdminDecider) InitGraph(stub shim.ChaincodeStubInterface) error {
	if err := a.setup(stub); err != nil {
		return errors.Wrap(err, "error initializing admin decider")
	}

	graph := memory.NewGraph()

	if err := policy.Configure(graph); err != nil {
		return errors.Wrap(err, "error configuring access control policy")
	}

	var (
		bytes []byte
		err   error
	)

	decider := pdp.NewDecider(graph)
	if ok, err := decider.HasPermissions(a.user, policy.BlossomObject, operations.InitBlossom); err != nil {
		return errors.Wrapf(err, "error checking if user can initialize blossom")
	} else if !ok {
		return errors.Errorf("user %s does not have permission to initialize blossom", a.user)
	}

	if bytes, err = graph.MarshalJSON(); err != nil {
		return errors.Wrap(err, "error serializing graph")
	}

	if err = stub.PutState("graph", bytes); err != nil {
		return errors.Wrap(err, "error updating graph on ledger")
	}

	return nil
}

func (a *AdminDecider) GetGraph(stub shim.ChaincodeStubInterface) (pip.Graph, error) {
	bytes, err := ledger.GetGraphBytes(stub)
	if err != nil {
		return nil, errors.Wrap(err, "error getting graph")
	}

	graph := memory.NewGraph()
	if err = json.Unmarshal(bytes, graph); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling graph")
	}

	return graph, nil
}

// UpdateGraph updates the NGAC graph with the given graph json. It first identifies the differences between the ledger
// graph and the provided graph and checks the requesting user has permission to carry out all the actions.
// If the user can carry out all the actions, the ledger graph is replaced with the graph provided.
func (a *AdminDecider) UpdateGraph(stub shim.ChaincodeStubInterface, ledgerGraph pip.Graph, jsonGraph pip.Graph) error {
	// check the client can execute request
	var (
		cmds []GraphCmd
		err  error
		user string
	)

	if user, err = GetUser(stub); err != nil {
		return fmt.Errorf("error getting user from request: %w", err)
	}

	// check the user can update graph
	if cmds, err = checkPermissions(user, ledgerGraph, jsonGraph); err != nil {
		return fmt.Errorf("client %v does not have permissions to update graph: %w", user, err)
	}

	// execute graph update
	if err := updateGraph(ledgerGraph, cmds...); err != nil {
		return fmt.Errorf("error updating NGAC graph: %w", err)
	}

	return nil
}

func checkPermissions(user string, ledgerGraph pip.Graph, jsonGraph pip.Graph) ([]GraphCmd, error) {
	// get differences between ledger graph and json graph
	graphCmds := differGraphs(ledgerGraph, jsonGraph)

	// check if user can execute all cmds
	for _, graphCmd := range graphCmds {
		if ok, err := graphCmd.CanExecute(user, ledgerGraph); err != nil {
			return nil, fmt.Errorf("error checking if user can execute command %v: %w", graphCmd, err)
		} else if !ok {
			return nil, fmt.Errorf("could not execute %s", graphCmd.String())
		}
	}

	return graphCmds, nil
}

// differGraphs finds the difference between the ledger graph and the json graph and formulates the differences into a
// series of commands
func differGraphs(ledgerGraph pip.Graph, jsonGraph pip.Graph) []GraphCmd {
	// get new nodes
	createdNodes := differCreatedNodes(ledgerGraph, jsonGraph)
	// get deleted nodes
	deletedNodes := differDeletedNodes(ledgerGraph, jsonGraph)
	// get new assignments
	assignments := differAssignments(ledgerGraph, jsonGraph, createdNodes)
	// get deleted assignments
	deassignments := differDeassignments(ledgerGraph, jsonGraph, deletedNodes)
	// get new associations
	associations := differAssociations(ledgerGraph, jsonGraph)
	// get deleted associations
	dissoctiations := differDissociations(ledgerGraph, jsonGraph)

	cmds := make([]GraphCmd, 0)

	// add created node cmds
	for _, cmd := range createdNodes {
		cmds = append(cmds, cmd)
	}

	// add the deleted node cmds
	for _, cmd := range deletedNodes {
		cmds = append(cmds, cmd)
	}

	// add the rest of the commands
	cmds = append(cmds, assignments...)
	cmds = append(cmds, deassignments...)
	cmds = append(cmds, associations...)
	cmds = append(cmds, dissoctiations...)

	return cmds
}

func differCreatedNodes(ledgerGraph pip.Graph, jsonGraph pip.Graph) map[string]CreateNodeCmd {
	createdNodes := make(map[string]CreateNodeCmd)

	// get created nodes
	nodes, _ := jsonGraph.GetNodes()
	for _, node := range nodes {
		// if the node already exists in the ledger graph - skip
		if ok, _ := ledgerGraph.Exists(node.Name); ok {
			continue
		}

		createdNodes[node.Name] = CreateNodeCmd{
			node:    node,
			parents: make(map[string]bool),
		}
	}

	return createdNodes
}

func differDeletedNodes(ledgerGraph pip.Graph, jsonGraph pip.Graph) map[string]DeleteNodeCmd {
	deletedNodes := make(map[string]DeleteNodeCmd)
	nodes, _ := ledgerGraph.GetNodes()

	for _, node := range nodes {
		// if the node does not exist in the json graph - add DeleteNodeCmd
		if ok, _ := jsonGraph.Exists(node.Name); ok {
			continue
		}

		deletedNodes[node.Name] = DeleteNodeCmd{
			name: node.Name,
		}
	}

	return deletedNodes
}

func differAssignments(ledgerGraph pip.Graph, jsonGraph pip.Graph, createdNodes map[string]CreateNodeCmd) []GraphCmd {
	cmds := make([]GraphCmd, 0)
	jsonAssignments, _ := jsonGraph.GetAssignments()
	ledgerAssignments, _ := ledgerGraph.GetAssignments()

	for child, parents := range jsonAssignments {
		if cmd, ok := createdNodes[child]; ok {
			// this node was created, add to the create node cmd
			cmd.parents = parents
			createdNodes[child] = cmd
		} else {
			// create assign commands for each assignment
			for parent := range parents {
				if ledgerAssignments[child][parent] {
					continue
				}

				cmds = append(cmds, AssignCmd{
					child:  child,
					parent: parent,
				})
			}
		}
	}

	return cmds
}

func differDeassignments(ledgerGraph pip.Graph, jsonGraph pip.Graph, deletedNodes map[string]DeleteNodeCmd) []GraphCmd {
	cmds := make([]GraphCmd, 0)
	jsonAssignments, _ := jsonGraph.GetAssignments()
	ledgerAssignments, _ := ledgerGraph.GetAssignments()

	for child, parents := range ledgerAssignments {
		// skip if the child has been deleted
		if _, ok := deletedNodes[child]; ok {
			continue
		}

		for parent := range parents {
			// if the assignment exists in the json graph it has not been deleted
			if jsonAssignments[child][parent] {
				continue
			}

			cmds = append(cmds, DeassignCmd{
				child:  child,
				parent: parent,
			})
		}
	}

	return cmds
}

func differAssociations(ledgerGraph pip.Graph, jsonGraph pip.Graph) []GraphCmd {
	cmds := make([]GraphCmd, 0)
	jsonAssociations, _ := jsonGraph.GetAssociations()
	ledgerAssociations, _ := ledgerGraph.GetAssociations()

	for subject, assocs := range jsonAssociations {
		for target, ops := range assocs {
			if _, ok := ledgerAssociations[subject][target]; ok {
				continue
			}

			cmds = append(cmds, AssociateCmd{
				subject:    subject,
				target:     target,
				operations: ops,
			})
		}
	}

	return cmds
}

func differDissociations(ledgerGraph pip.Graph, jsonGraph pip.Graph) []GraphCmd {
	cmds := make([]GraphCmd, 0)
	jsonAssociations, _ := jsonGraph.GetAssociations()
	ledgerAssociations, _ := ledgerGraph.GetAssociations()

	for subject, assocs := range ledgerAssociations {
		for target := range assocs {
			// if the association still exists in the json graph it was not deleted
			if _, ok := jsonAssociations[subject][target]; ok {
				continue
			}

			cmds = append(cmds, DissociateCmd{
				subject: subject,
				target:  target,
			})
		}
	}

	return cmds
}

func updateGraph(graph pip.Graph, cmds ...GraphCmd) error {
	for _, cmd := range cmds {
		if err := cmd.Execute(graph); err != nil {
			return err
		}
	}

	return nil
}
