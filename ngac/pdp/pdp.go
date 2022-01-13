package pdp

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/PM-Master/policy-machine-go/policy"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/adminmsp"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/common"
	events "github.com/usnistgov/blossom/chaincode/ngac/epp"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
)

func IsNGACInitialized(stub shim.ChaincodeStubInterface, collName string) (bool, error) {
	return common.IsNGACInitialized(stub, collName)
}

func InitCatalogNGAC(stub shim.ChaincodeStubInterface, collection string) error {
	// check if this has already been called.  An error is thrown if this has not been called before
	if _, err := common.GetPvtCollPolicyStore(stub, collection); err == nil {
		return fmt.Errorf("ngac initialization function has already been called")
	}

	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return err
	}

	// only member of the predefined AdminMSP can initialize the Catalog ngac store
	if mspid != adminmsp.AdminMSP {
		return fmt.Errorf("users in MSP %s do not have pemrission to initialize ngac graphs", mspid)
	}

	// the admin user for the graph will be the user that performs the initialization
	adminUser, err := common.GetUser(stub)
	if err != nil {
		return err
	}

	policyStore, err := pap.LoadCatalogPolicy(adminUser, adminmsp.AdminMSP)
	if err != nil {
		return errors.Wrap(err, "error loading catalog policy")
	}

	return common.PutPvtCollPolicyStore(stub, collection, policyStore)
}

func InitAccountNGAC(stub shim.ChaincodeStubInterface, collection, account string, acctPvt model.AccountPrivate) error {
	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return err
	}

	// only member of the predefined AdminMSP can initialize the account ngac store
	if mspid != adminmsp.AdminMSP {
		return fmt.Errorf("users in MSP %s do not have pemrission to initialize ngac graphs", mspid)
	}

	// the admin user for the graph will be the user that performs the initialization
	adminUser, err := common.GetUser(stub)
	if err != nil {
		return err
	}

	// load the policy into memory
	policyStore, err := pap.LoadAccountPolicy(adminUser, adminmsp.AdminMSP)
	if err != nil {
		return errors.Wrap(err, "error loading account policy")
	}

	// process the approve account
	if err = events.ProcessApproveAccount(stub, collection, account, acctPvt, policyStore); err != nil {
		return fmt.Errorf("error processing approve_account event: %v", err)
	}

	// store the policy in the collection
	return common.PutPvtCollPolicyStore(stub, collection, policyStore)
}

func CanUploadATO(stub shim.ChaincodeStubInterface, pvtCollName, account string) error {
	return check(stub, pvtCollName, pap.AccountObjectName(account), "upload_ato")
}

func CanUpdateAccountStatus(stub shim.ChaincodeStubInterface, pvtCollName, account string) error {
	return check(stub, pvtCollName, pap.AccountObjectName(account), "update_account_status")
}

func CanRequestCheckout(stub shim.ChaincodeStubInterface, pvtCollName, account string) error {
	return check(stub, pvtCollName, pap.AccountObjectName(account), "check_out")
}

func CanApproveCheckout(stub shim.ChaincodeStubInterface, pvtCollName, account string) error {
	return check(stub, pvtCollName, pap.BlossomObject, "approve_checkout")
}

func CanInitiateCheckIn(stub shim.ChaincodeStubInterface, pvtCollName, account string) error {
	return check(stub, pvtCollName, pap.AccountObjectName(account), "initiate_check_in")
}

func CanProcessCheckIn(stub shim.ChaincodeStubInterface, pvtCollName, account string) error {
	return check(stub, pvtCollName, pap.AccountObjectName(account), "process_check_in")
}

func CanReportSwID(stub shim.ChaincodeStubInterface, pvtCollName, account string) error {
	return check(stub, pvtCollName, pap.AccountObjectName(account), "report_swid")
}

func CanDeleteSwID(stub shim.ChaincodeStubInterface, pvtCollName, account string) error {
	return check(stub, pvtCollName, pap.AccountObjectName(account), "delete_swid")
}

func CanOnboardAsset(stub shim.ChaincodeStubInterface, pvtCollName string) error {
	return check(stub, pvtCollName, "assets", "onboard_asset")
}

func CanOffboardAsset(stub shim.ChaincodeStubInterface, pvtCollName string) error {
	return check(stub, pvtCollName, "assets", "offboard_asset")
}

func check(stub shim.ChaincodeStubInterface, pvtCollName, target, permission string) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user")
	}

	policyStore, err := common.GetPvtCollPolicyStore(stub, pvtCollName)
	if err != nil {
		return err
	}

	decider := pdp.NewDecider(policyStore.Graph(), policyStore.Prohibitions())
	if ok, err := decider.HasPermissions(user, target, permission); err != nil {
		return errors.Wrapf(err, "error checking if user %s can %s on %s", user, permission, target)
	} else if !ok {
		return fmt.Errorf("user %q does not have permission %q on %q", user, permission, target)
	}

	return nil
}

func checkWithPolicyStore(stub shim.ChaincodeStubInterface, policyStore policy.Store, target, permission string) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user")
	}

	decider := pdp.NewDecider(policyStore.Graph(), policyStore.Prohibitions())
	if ok, err := decider.HasPermissions(user, target, permission); err != nil {
		return errors.Wrapf(err, "error checking if user %s can %s on %s", user, permission, target)
	} else if !ok {
		return fmt.Errorf("user %s does not have permission %s on %s", user, permission, target)
	}

	return nil
}
