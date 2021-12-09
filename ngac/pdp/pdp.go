package pdp

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/ngac/common"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
)

func InitCatalogNGAC(stub shim.ChaincodeStubInterface, collection string) error {
	fe, err := pap.LoadCatalogPolicy()
	if err != nil {
		return errors.Wrap(err, "error loading catalog policy")
	}

	if err = canInit(stub, fe, collection); err != nil {
		return err
	}

	return common.PutPvtCollFunctionalEntity(stub, collection, fe)
}

func canInit(stub shim.ChaincodeStubInterface, fe ngac.FunctionalEntity, pvtCollName string) error {
	return checkWithFuncEnt(stub, fe, pvtCollName, pap.BlossomObject, "init_blossom")
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

	fe, err := common.GetPvtCollFunctionalEntity(stub, pvtCollName)
	if err != nil {
		return err
	}

	decider := pdp.NewDecider(fe.Graph(), fe.Prohibitions())
	if ok, err := decider.HasPermissions(user, target, permission); err != nil {
		return errors.Wrapf(err, "error checking if user %s can %s on %s", user, permission, target)
	} else if !ok {
		return fmt.Errorf("user %s does not have permission %s on %s", user, permission, target)
	}

	return nil
}

func checkWithFuncEnt(stub shim.ChaincodeStubInterface, fe ngac.FunctionalEntity, pvtCollName, target, permission string) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user")
	}

	decider := pdp.NewDecider(fe.Graph(), fe.Prohibitions())
	if ok, err := decider.HasPermissions(user, target, permission); err != nil {
		return errors.Wrapf(err, "error checking if user %s can %s on %s", user, permission, target)
	} else if !ok {
		return fmt.Errorf("user %s does not have permission %s on %s", user, permission, target)
	}

	return nil
}
