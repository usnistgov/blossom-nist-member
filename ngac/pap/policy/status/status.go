package status

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
)

const (
	PolicyClassName     = "Status"
	UserAttributeName   = "Status_UA"
	ObjectAttributeName = "Status_OA"
	AccountsOA          = "status_accounts_OA"
	AssetsOA            = "status_assets_OA"
	SwIDsOA             = "status_swids_OA"
	ActiveUA            = "active"
	PendingUA           = "pending"
	InactiveUA          = "inactive"
)

func Configure(graph pip.Graph, adminUA string) error {
	statusPC, err := graph.CreateNode(PolicyClassName, pip.PolicyClass, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating Status policy class node")
	}

	// DAC default nodes
	statusUA, err := graph.CreateNode(UserAttributeName, pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating Status user attribute node")
	}

	if err = graph.Assign(statusUA.Name, statusPC.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", statusUA.Name, statusPC.Name)
	}

	statusOA, err := graph.CreateNode(ObjectAttributeName, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating Status object attribute node")
	}

	if err = graph.Assign(statusOA.Name, statusPC.Name); err != nil {
		return errors.Wrapf(err, "error assigning %q to %q", statusOA.Name, statusPC.Name)
	}

	// associate the admin UA with the default nodes
	if err = graph.Associate(adminUA, statusUA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating admin user attribute with status user attribute")
	}

	if err = graph.Associate(adminUA, statusOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating admin user attribute with status object attribute")
	}

	// create OAs for accounts and licenses
	accountsOA, err := graph.CreateNode(AccountsOA, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating accounts object attribute in status policy class")
	}

	assetsOA, err := graph.CreateNode(AssetsOA, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating licenses object attribute in status policy class")
	}

	swidsOA, err := graph.CreateNode(SwIDsOA, pip.ObjectAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating swids object attribute in status policy class")
	}

	if err = graph.Assign(accountsOA.Name, statusOA.Name); err != nil {
		return errors.Wrapf(err, "error assigning accounts object attribute to status object attribute")
	}

	if err = graph.Assign(assetsOA.Name, statusOA.Name); err != nil {
		return errors.Wrapf(err, "error assigning licenses object attribute to status object attribute")
	}

	if err = graph.Assign(swidsOA.Name, statusOA.Name); err != nil {
		return errors.Wrapf(err, "error assigning swids object attribute to status object attribute")
	}

	// create status user attributes: active, pending, inactive
	activeUA, err := graph.CreateNode(ActiveUA, pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating active user attribute in status policy class")
	}

	pendingUA, err := graph.CreateNode(PendingUA, pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating pending user attribute in status policy class")
	}

	inactiveUA, err := graph.CreateNode(InactiveUA, pip.UserAttribute, nil)
	if err != nil {
		return errors.Wrapf(err, "error creating inactive user attribute in status policy class")
	}

	if err = graph.Assign(activeUA.Name, statusUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning active user attribute to status user attribute")
	}

	if err = graph.Assign(pendingUA.Name, statusUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning pending user attribute to status user attribute")
	}

	// assign inactive to pending because the permissions are the same
	if err = graph.Assign(inactiveUA.Name, pendingUA.Name); err != nil {
		return errors.Wrapf(err, "error assigning inactive user attribute to status user attribute")
	}

	// associate status UAs with account and license OAs
	if err = graph.Associate(activeUA.Name, accountsOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating active user attribute with accounts object attribute")
	}

	if err = graph.Associate(pendingUA.Name, accountsOA.Name, pip.ToOps(operations.ViewAccount, operations.UploadATO,
		operations.ViewATO, operations.ViewMSPID, operations.ViewUsers, operations.ViewStatus, operations.ViewAccountLicenses)); err != nil {
		return errors.Wrapf(err, "error associating active user attribute with accounts object attribute")
	}

	if err = graph.Associate(activeUA.Name, assetsOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating active user attribute with licenses object attribute")
	}
	if err = graph.Associate(activeUA.Name, swidsOA.Name, pip.ToOps(pip.AllOps)); err != nil {
		return errors.Wrapf(err, "error associating active user attribute with swids object attribute")
	}

	return nil
}
