package rbac

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
)

const (
	PolicyClassName         = "RBAC"
	ObjectAttributeName     = "RBAC_OA"
	UserAttributeName       = "RBAC_UA"
	SystemOwnerUA           = "SystemOwner"
	AcquisitionSpecialistUA = "AcquisitionSpecialist"
	SystemAdministratorUA   = "SystemAdministrator"
	AccountsOA              = "Accounts"
	AccountsUA              = "Accounts_UA"
	AssetsOA                = "Assets"
	SwIDsOA                 = "SwIDs"
)

var SystemOwnerPermissions = pip.ToOps(
	operations.ViewAccount,
	operations.ViewAccountLicenses,
	operations.UploadATO,
	operations.ViewATO,
	operations.ViewMSPID,
	operations.ViewUsers,
	operations.ViewStatus)

var SystemAdminLicensesPermissions = pip.ToOps(
	operations.ViewAsset,
	operations.CheckOut,
	operations.CheckIn,
	operations.ReportSwid)

var SystemAdminAccountsPermissions = pip.ToOps(
	operations.ViewAccount,
	operations.ViewAccountLicenses,
)

var AcqSpecLicensesPermissions = pip.ToOps(
	operations.ViewAsset)

var AcqSpecAccountsPermissions = pip.ToOps(
	operations.ViewAccountLicenses,
	operations.ViewAccount,
	operations.ViewStatus)

var SystemAdminSwidPermissions = pip.ToOps(
	operations.ViewSwID,
	operations.ReportSwid)

var AcqSpecswidPermissions = pip.ToOps(
	operations.ViewSwID)

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

	// create a UA to hold each account UA
	accountsUA, err := graph.CreateNode(AccountsUA, pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating accounts base user attribute")
	}

	if err = graph.Assign(accountsUA.Name, rbacUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", accountsUA.Name, rbacUA.Name)
	}

	// associate the admin UA with the default attributes, giving them * permissions on all nodes in the policy class
	if err = graph.Associate(adminUA, rbacUA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", adminUA, rbacUA.Name)
	}
	if err = graph.Associate(adminUA, rbacOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", adminUA, rbacOA.Name)
	}

	accountsOA, err := graph.CreateNode(AccountsOA, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating accounts base object attribute")
	}

	if err = graph.Assign(accountsOA.Name, rbacOA.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", accountsOA.Name, rbacOA.Name)
	}

	licensesOA, err := graph.CreateNode(AssetsOA, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating licenses base object attribute")
	}

	if err = graph.Assign(licensesOA.Name, rbacOA.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", licensesOA.Name, rbacOA.Name)
	}

	swidsOA, err := graph.CreateNode(SwIDsOA, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating swids base object attribute")
	}

	if err = graph.Assign(swidsOA.Name, rbacOA.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", swidsOA.Name, rbacOA.Name)
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

	// system owners are only associated with accounts
	if err = graph.Associate(systemOwnersUA.Name, accountsOA.Name, SystemOwnerPermissions); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", systemOwnersUA.Name, accountsOA.Name)
	}

	// system admins are associated with licenses and accounts
	if err = graph.Associate(systemAdminsUA.Name, licensesOA.Name, SystemAdminLicensesPermissions); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", systemAdminsUA.Name, licensesOA.Name)
	}

	if err = graph.Associate(systemAdminsUA.Name, accountsOA.Name, SystemAdminAccountsPermissions); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", systemAdminsUA.Name, accountsOA.Name)
	}

	if err = graph.Associate(systemAdminsUA.Name, swidsOA.Name, SystemAdminSwidPermissions); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", systemAdminsUA.Name, swidsOA.Name)
	}

	// acquisition specialists are associated with licenses and accounts
	if err = graph.Associate(acqSpecUA.Name, licensesOA.Name, AcqSpecLicensesPermissions); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", acqSpecUA.Name, licensesOA.Name)
	}

	if err = graph.Associate(acqSpecUA.Name, accountsOA.Name, AcqSpecAccountsPermissions); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", acqSpecUA.Name, accountsOA.Name)
	}

	if err = graph.Associate(acqSpecUA.Name, swidsOA.Name, AcqSpecswidPermissions); err != nil {
		return errors.Wrapf(err, "error associating %q with %q", acqSpecUA.Name, swidsOA.Name)
	}

	return nil
}
