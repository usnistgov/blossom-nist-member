package epp

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/epp"
	"github.com/PM-Master/policy-machine-go/policy"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/common"
)

func process(stub shim.ChaincodeStubInterface, collection string, evtCtx epp.EventContext, policyStore policy.Store) error {
	eventProcessor := epp.NewEPP(policyStore)

	if err := eventProcessor.ProcessEvent(evtCtx); err != nil {
		return err
	}

	return common.PutPvtCollPolicyStore(stub, collection, policyStore)
}

func ProcessApproveAccount(stub shim.ChaincodeStubInterface, pvtCollName, account string, acctPvt *model.AccountPrivate, store policy.Store) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return err
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "approve_account",
		Args: map[string]string{
			"accountName": account,
			"sysOwner":    common.FormatUsername(acctPvt.Users.SystemOwner, account),
			"sysAdmin":    common.FormatUsername(acctPvt.Users.SystemAdministrator, account),
			"acqSpec":     common.FormatUsername(acctPvt.Users.AcquisitionSpecialist, account),
		},
	}

	return process(stub, pvtCollName, evtCtx, store)
}

func UpdateAccountStatusEvent(stub shim.ChaincodeStubInterface, accountName, pvtColl string, status model.Status) error {
	store, err := common.GetPvtCollPolicyStore(stub, pvtColl)
	if err != nil {
		return err
	}

	var f func(shim.ChaincodeStubInterface, string, string, policy.Store) error
	switch status {
	case model.PendingApproval:
		f = ProcessSetAccountPending
	case model.PendingATO:
		f = ProcessSetAccountPending
	case model.Active:
		f = ProcessSetAccountActive
	case model.InactiveDenied:
		f = ProcessSetAccountInactive
	case model.InactiveATO:
		f = ProcessSetAccountInactive
	case model.InactiveOptOut:
		f = ProcessSetAccountInactive
	case model.InactiveSecurityRisk:
		f = ProcessSetAccountInactive
	case model.InactiveROB:
		f = ProcessSetAccountInactive
	default:
		return fmt.Errorf("unknown status: %s", status)
	}

	return f(stub, pvtColl, accountName, store)
}

func ProcessSetAccountActive(stub shim.ChaincodeStubInterface, pvtCollName, account string, store policy.Store) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user from stub")
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "set_account_active",
		Args: map[string]string{
			"account": account,
		},
	}

	return process(stub, pvtCollName, evtCtx, store)
}

func ProcessSetAccountPending(stub shim.ChaincodeStubInterface, pvtCollName, account string, store policy.Store) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user from stub")
	}

	policyStore, err := common.GetPvtCollPolicyStore(stub, pvtCollName)
	if err != nil {
		return errors.Wrap(err, "error getting ngac components")
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "set_account_pending",
		Args: map[string]string{
			"account": account,
		},
	}

	return process(stub, pvtCollName, evtCtx, policyStore)
}

func ProcessSetAccountInactive(stub shim.ChaincodeStubInterface, pvtCollName, account string, store policy.Store) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user from stub")
	}

	policyStore, err := common.GetPvtCollPolicyStore(stub, pvtCollName)
	if err != nil {
		return errors.Wrap(err, "error getting ngac components")
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "set_account_inactive",
		Args: map[string]string{
			"account": account,
		},
	}

	return process(stub, pvtCollName, evtCtx, policyStore)
}

func ProcessOnboardAsset(stub shim.ChaincodeStubInterface, pvtCollName, assetID string) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user from stub")
	}

	policyStore, err := common.GetPvtCollPolicyStore(stub, pvtCollName)
	if err != nil {
		return errors.Wrap(err, "error getting ngac components")
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "onboard_asset",
		Args: map[string]string{
			"asset_id": assetID,
		},
	}

	return process(stub, pvtCollName, evtCtx, policyStore)
}

func ProcessOffboardAsset(stub shim.ChaincodeStubInterface, pvtCollName, assetID string) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user from stub")
	}

	policyStore, err := common.GetPvtCollPolicyStore(stub, pvtCollName)
	if err != nil {
		return errors.Wrap(err, "error getting ngac components")
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "offboard_asset",
		Args: map[string]string{
			"asset_id": assetID,
		},
	}

	return process(stub, pvtCollName, evtCtx, policyStore)
}
