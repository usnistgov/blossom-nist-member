package pdp

import (
	"fmt"

	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/PM-Master/policy-machine-go/pip"
)

type (
	// GraphCmd provides methods used to execute graph commands on an NGAC graph
	GraphCmd interface {
		// CanExecute checks that the user can execute the GraphCmd on the graph
		CanExecute(user string, graph pip.Graph) (bool, error)
		// execute the GraphCmd on the graph
		Execute(graph pip.Graph) error

		fmt.Stringer
	}

	// CreateNodeCmd stores the node being created and the parents of that node
	CreateNodeCmd struct {
		node    pip.Node
		parents map[string]bool
	}

	// DeleteNodeCmd stores the name of the node being deleted
	DeleteNodeCmd struct {
		name string
	}

	// AssignCmd stores the child and parent of an assignment (child -> parent)
	AssignCmd struct {
		child  string
		parent string
	}

	// DeassignCmd stores the child and parents of an assignment to delete
	DeassignCmd struct {
		child  string
		parent string
	}

	// AssociateCmd stores the subject, target, and operations of a new association
	AssociateCmd struct {
		subject    string
		target     string
		operations pip.Operations
	}

	// DissociateCmd stores the subject and target of an association to delete
	DissociateCmd struct {
		subject string
		target  string
	}
)

const (
	CreateNodePermission   = "create node"
	DeleteNodePermission   = "delete node"
	AssignPermission       = "assign"
	AssignToPermission     = "assign to"
	DeassignPermission     = "deassign"
	DeassignFromPermission = "deassign from"
	AssociatePermission    = "associate"
	DissociatePermission   = "dissociate"
)

func (c CreateNodeCmd) CanExecute(user string, graph pip.Graph) (bool, error) {
	decider := pdp.NewDecider(graph)
	for parent := range c.parents {
		if ok, err := decider.Decide(user, parent, CreateNodePermission); err != nil {
			return false, err
		} else if !ok {
			return false, nil
		}
	}

	return true, nil
}

func (c CreateNodeCmd) Execute(graph pip.Graph) error {
	// create the node
	if _, err := graph.CreateNode(c.node.Name, c.node.Kind, c.node.Properties); err != nil {
		return fmt.Errorf("could not execute command %v: %v", c, err)
	}

	// assign the node to the parents
	for parent := range c.parents {
		if err := graph.Assign(c.node.Name, parent); err != nil {
			return err
		}
	}

	return nil
}

func (c CreateNodeCmd) String() string {
	return fmt.Sprintf("create node %v in %v", c.node, c.parents)
}

func (c DeleteNodeCmd) CanExecute(user string, graph pip.Graph) (bool, error) {
	decider := pdp.NewDecider(graph)

	parents, err := graph.GetParents(c.name)
	if err != nil {
		return false, fmt.Errorf("could not get parents of %q", c.name)
	}

	if ok, err := decider.Decide(user, c.name, DeleteNodePermission); err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}

	for parent := range parents {
		if ok, err := decider.Decide(user, parent, DeassignFromPermission); err != nil {
			return false, err
		} else if !ok {
			return false, nil
		}
	}

	return true, nil
}

func (c DeleteNodeCmd) Execute(graph pip.Graph) error {
	return graph.DeleteNode(c.name)
}

func (c DeleteNodeCmd) String() string {
	return fmt.Sprintf("delete node %v", c.name)
}

func (c AssignCmd) CanExecute(user string, graph pip.Graph) (bool, error) {
	decider := pdp.NewDecider(graph)

	// check user can assign child
	if ok, err := decider.Decide(user, c.child, AssignPermission); err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}

	// check user can assign to parent
	return decider.Decide(user, c.parent, AssignToPermission)
}

func (c AssignCmd) Execute(graph pip.Graph) error {
	return graph.Assign(c.child, c.parent)
}

func (c AssignCmd) String() string {
	return fmt.Sprintf("assign %v to %v", c.child, c.parent)
}

func (c DeassignCmd) CanExecute(user string, graph pip.Graph) (bool, error) {
	decider := pdp.NewDecider(graph)

	// check user can assign child
	if ok, err := decider.Decide(user, c.child, DeassignPermission); err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}

	// check user can assign to parent
	return decider.Decide(user, c.parent, DeassignFromPermission)
}

func (c DeassignCmd) Execute(graph pip.Graph) error {
	return graph.Deassign(c.child, c.parent)
}

func (c DeassignCmd) String() string {
	return fmt.Sprintf("deassign %v from %v", c.child, c.parent)
}

func (c AssociateCmd) CanExecute(user string, graph pip.Graph) (bool, error) {
	decider := pdp.NewDecider(graph)

	// check user can associate the subject
	if ok, err := decider.Decide(user, c.subject, AssociatePermission); err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}

	// check user can associate the target
	return decider.Decide(user, c.target, AssociatePermission)
}

func (c AssociateCmd) Execute(graph pip.Graph) error {
	return graph.Associate(c.subject, c.target, c.operations)
}

func (c AssociateCmd) String() string {
	return fmt.Sprintf("associate %v with %v with ops %v", c.subject, c.target, c.operations)
}

func (c DissociateCmd) CanExecute(user string, graph pip.Graph) (bool, error) {
	decider := pdp.NewDecider(graph)

	// check user can assign child
	if ok, err := decider.Decide(user, c.subject, DissociatePermission); err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}

	// check user can assign to parent
	return decider.Decide(user, c.target, DissociatePermission)
}

func (c DissociateCmd) Execute(graph pip.Graph) error {
	return graph.Dissociate(c.subject, c.target)
}

func (c DissociateCmd) String() string {
	return fmt.Sprintf("dissociate %v from %v", c.subject, c.target)
}
