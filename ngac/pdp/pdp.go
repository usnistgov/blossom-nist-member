package pdp

import (
	"encoding/json"
	"fmt"
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/adminmsp"
	"github.com/usnistgov/blossom/chaincode/collections"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/common"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
)

func InitCatalogNGAC(stub shim.ChaincodeStubInterface) error {
	// check if this has already been called.  An error is thrown if this has not been called before
	if _, err := common.GetPvtCollPolicyStore(stub, collections.Catalog()); err == nil {
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

	return common.PutPvtCollPolicyStore(stub, policyStore)
}

func CanApproveAccount(stub shim.ChaincodeStubInterface) error {
	return check(stub, pap.BlossomObject, "approve_account")
}

func CanUploadATO(stub shim.ChaincodeStubInterface, account string) error {
	return check(stub, pap.AccountObjectName(account), "upload_ato")
}

func CanUpdateAccountStatus(stub shim.ChaincodeStubInterface, account string) error {
	return check(stub, pap.AccountObjectName(account), "update_account_status")
}

func CanRequestCheckout(stub shim.ChaincodeStubInterface, account string) error {
	return check(stub, pap.AccountObjectName(account), "check_out")
}

func CanApproveCheckout(stub shim.ChaincodeStubInterface, account string) error {
	return check(stub, pap.BlossomObject, "approve_checkout")
}

func CanInitiateCheckIn(stub shim.ChaincodeStubInterface, account string) error {
	return check(stub, pap.AccountObjectName(account), "initiate_check_in")
}

func CanProcessCheckIn(stub shim.ChaincodeStubInterface, account string) error {
	return check(stub, pap.AccountObjectName(account), "process_check_in")
}

func CanReportSwID(stub shim.ChaincodeStubInterface, account string) error {
	return check(stub, pap.AccountObjectName(account), "report_swid")
}

func CanDeleteSwID(stub shim.ChaincodeStubInterface, account string) error {
	return check(stub, pap.AccountObjectName(account), "delete_swid")
}

func CanOnboardAsset(stub shim.ChaincodeStubInterface) error {
	return check(stub, "assets", "onboard_asset")
}

func CanOffboardAsset(stub shim.ChaincodeStubInterface) error {
	return check(stub, "assets", "offboard_asset")
}

func check(stub shim.ChaincodeStubInterface, target, permission string) error {
	user, err := common.GetUsername(stub)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}

	policyStore, err := common.GetPvtCollPolicyStore(stub, collections.Catalog())
	if err != nil {
		return err
	}

	account, err := cid.GetMSPID(stub)
	if err != nil {
		return err
	}

	// skip this step for users in the adminmsp as they dont have account roles
	if account != adminmsp.AdminMSP {
		role, err := getRole(stub, user, account)
		if err != nil {
			return err
		}

		// assign the user to the account and role
		if err = policyStore.Graph().Assign(user, pap.AccountUA(account)); err != nil {
			return fmt.Errorf("error assigning user %s to account UA %s: %v", user, account, err)
		}
		if err = policyStore.Graph().Assign(user, role); err != nil {
			return fmt.Errorf("error assigning user %s to role UA %s: %v", user, role, err)
		}
	} else {
		user = common.FormatUsername(user, account)
	}

	decider := pdp.NewDecider(policyStore.Graph(), policyStore.Prohibitions())
	if ok, err := decider.HasPermissions(user, target, permission); err != nil {
		return errors.Wrapf(err, "error checking if user %s can %s on %s", user, permission, target)
	} else if !ok {
		return fmt.Errorf("user %q does not have permission %q on %q", user, permission, target)
	}

	return nil
}

func getRole(stub shim.ChaincodeStubInterface, user, account string) (role string, err error) {
	// get role from pvtcol of account if mspid == adminmsp skip this step
	acctColl := collections.Account(account)

	data, err := stub.GetPrivateData(acctColl, model.AccountKey(account))
	if err != nil {
		return "", fmt.Errorf("error getting the users of account: %v", account)
	}

	acctPvt := model.AccountPrivate{}
	if err = json.Unmarshal(data, &acctPvt); err != nil {
		return "", fmt.Errorf("error unmarshaling private account details: %v", err)
	}

	if acctPvt.Users.SystemOwner == user {
		role = "SystemOwner"
	} else if acctPvt.Users.SystemOwner == user {
		role = "SystemAdministrator"
	} else if acctPvt.Users.SystemOwner == user {
		role = "AcquisitionSpecialist"
	} else {
		return "", fmt.Errorf("user %s is not registered with the account %s", user, account)
	}

	return
}
