package rbac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
)

const (
	PolicyClassName         = "RBAC"
	ObjectAttributeName     = "RBAC_OA"
	UserAttributeName       = "RBAC_UA"
	SystemOwnerUA           = "SystemOwner"
	AcquisitionSpecialistUA = "AcquisitionSpecialist"
	SystemAdministratorUA   = "SystemAdministrator"
	AgenciesOA              = "Agencies"
	AgenciesUA              = "Agencies_UA"
	LicensesOA              = "Licenses"
)

func Configure(graph pip.Graph, adminUA string) error {
	// create RBAC policy class node
	rbacPC, err := graph.CreateNode(PolicyClassName, pip.PolicyClass, nil)
	if err != nil {
		return errors.Wrap(err, "error creating RBAC policy class")
	}

	// create default attributes
	// these are used when a user wants to create a new attribute in the policy class
	// we can't check if the user has permissions to create a new node in a policy class
	// we can check if they can create a new node in an already existing node
	rbacUA, err := graph.CreateNode(UserAttributeName, pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrap(err, "error creating RBAC user attribute")
	}

	if err = graph.Assign(rbacUA.Name, rbacPC.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", rbacUA.Name, rbacPC.Name)
	}

	rbacOA, err := graph.CreateNode(ObjectAttributeName, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating RBAC object attribute")
	}

	if err = graph.Assign(rbacOA.Name, rbacPC.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", rbacOA.Name, rbacPC.Name)
	}

	// create a UA to hold each agency UA
	agenciesUA, err := graph.CreateNode(AgenciesUA, pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating agencies base user attribute")
	}

	if err = graph.Assign(agenciesUA.Name, rbacUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", agenciesUA.Name, rbacUA.Name)
	}

	// associate the admin UA with the default attributes, giving them * permissions on all nodes in the policy class
	if err = graph.Associate(adminUA, rbacUA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", adminUA, rbacUA.Name)
	}
	if err = graph.Associate(adminUA, rbacOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", adminUA, rbacOA.Name)
	}

	agenciesOA, err := graph.CreateNode(AgenciesOA, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating agencies base object attribute")
	}

	if err = graph.Assign(agenciesOA.Name, rbacOA.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", agenciesOA.Name, rbacOA.Name)
	}

	licensesOA, err := graph.CreateNode(LicensesOA, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating licenses base object attribute")
	}

	if err = graph.Assign(licensesOA.Name, rbacOA.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", licensesOA.Name, rbacOA.Name)
	}

	systemOwnersUA, err := graph.CreateNode(SystemOwnerUA, pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating SystemOwners user attribute")
	}

	systemAdminsUA, err := graph.CreateNode(SystemAdministratorUA, pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating SystemAdmins user attribute")
	}

	acqSpecUA, err := graph.CreateNode(AcquisitionSpecialistUA, pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating AcquisitionSpecialists user attribute")
	}

	if err = graph.Assign(systemOwnersUA.Name, rbacUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", systemOwnersUA.Name, rbacUA.Name)
	}

	if err = graph.Assign(systemAdminsUA.Name, rbacUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", systemAdminsUA.Name, rbacUA.Name)
	}

	if err = graph.Assign(acqSpecUA.Name, rbacUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", acqSpecUA.Name, rbacUA.Name)
	}

	// system owners can view "agencies"
	if err = graph.Associate(systemOwnersUA.Name, agenciesOA.Name, SystemOwnerPermissions); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", systemOwnersUA.Name, agenciesOA.Name)
	}
	// system admins can read, assign, deassign (assign and deassign for the license keys) "licenses"
	if err = graph.Associate(systemAdminsUA.Name, licensesOA.Name, SystemAdminPermissions); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", systemAdminsUA.Name, licensesOA.Name)
	}
	// acquisition specialists can audit agency licenses
	if err = graph.Associate(acqSpecUA.Name, licensesOA.Name, AcqSpecPermissions); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", acqSpecUA.Name, licensesOA.Name)
	}

	return nil
}
