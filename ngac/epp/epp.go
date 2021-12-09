package epp

import (
	"github.com/PM-Master/policy-machine-go/epp"
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/ngac/common"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
)

func process(stub shim.ChaincodeStubInterface, collection string, evtCtx epp.EventContext, fe ngac.FunctionalEntity) error {
	eventProcessor := epp.NewEPP(fe)

	if err := eventProcessor.ProcessEvent(evtCtx); err != nil {
		return err
	}

	return common.PutPvtCollFunctionalEntity(stub, collection, fe)
}

func ProcessRequestAccount(stub shim.ChaincodeStubInterface, pvtCollName, account, sysOwner, sysAdmin, acqSpec string) error {
	fe, err := pap.LoadAccountPolicy()
	if err != nil {
		return err
	}

	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user from stub")
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "request_account",
		Args: map[string]string{
			"account_name": account,
			"sysOwner":     common.FormatUsername(sysOwner, account),
			"sysAdmin":     common.FormatUsername(sysAdmin, account),
			"acqSpec":      common.FormatUsername(acqSpec, account),
		},
	}

	return process(stub, pvtCollName, evtCtx, fe)
}

func ProcessSetAccountActive(stub shim.ChaincodeStubInterface, pvtCollName, account string) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user from stub")
	}

	fe, err := common.GetPvtCollFunctionalEntity(stub, pvtCollName)
	if err != nil {
		return errors.Wrap(err, "error getting ngac components")
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "set_account_active",
		Args: map[string]string{
			"account": account,
		},
	}

	return process(stub, pvtCollName, evtCtx, fe)
}

func ProcessSetAccountPending(stub shim.ChaincodeStubInterface, pvtCollName, account string) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user from stub")
	}

	fe, err := common.GetPvtCollFunctionalEntity(stub, pvtCollName)
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

	return process(stub, pvtCollName, evtCtx, fe)
}

func ProcessSetAccountInactive(stub shim.ChaincodeStubInterface, pvtCollName, account string) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user from stub")
	}

	fe, err := common.GetPvtCollFunctionalEntity(stub, pvtCollName)
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

	return process(stub, pvtCollName, evtCtx, fe)
}

func ProcessOnboardAsset(stub shim.ChaincodeStubInterface, pvtCollName, assetID string) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user from stub")
	}

	fe, err := common.GetPvtCollFunctionalEntity(stub, pvtCollName)
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

	return process(stub, pvtCollName, evtCtx, fe)
}

func ProcessOffboardAsset(stub shim.ChaincodeStubInterface, pvtCollName, assetID string) error {
	user, err := common.GetUser(stub)
	if err != nil {
		return errors.Wrap(err, "error getting user from stub")
	}

	fe, err := common.GetPvtCollFunctionalEntity(stub, pvtCollName)
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

	return process(stub, pvtCollName, evtCtx, fe)
}
