package epp

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/epp"
	"github.com/PM-Master/policy-machine-go/policy"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/usnistgov/blossom/chaincode/collections"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/common"
)

func process(ctx contractapi.TransactionContextInterface, evtCtx epp.EventContext, policyStore policy.Store) error {
	eventProcessor := epp.NewEPP(policyStore)

	if err := eventProcessor.ProcessEvent(evtCtx); err != nil {
		return err
	}

	return common.PutPvtCollPolicyStore(ctx, policyStore)
}

func ProcessApproveAccount(ctx contractapi.TransactionContextInterface, account string) error {
	user, err := common.GetUser(ctx)
	if err != nil {
		return err
	}

	store, err := common.GetPvtCollPolicyStore(ctx, collections.Catalog())
	if err != nil {
		return err
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "approve_account",
		Args: map[string]string{
			"accountName": account,
		},
	}

	return process(ctx, evtCtx, store)
}

func UpdateAccountStatusEvent(ctx contractapi.TransactionContextInterface, accountName, pvtColl string, status model.Status) error {
	store, err := common.GetPvtCollPolicyStore(ctx, collections.Catalog())
	if err != nil {
		return err
	}

	var f func(contractapi.TransactionContextInterface, string, string, policy.Store) error
	switch status {
	case model.PendingApproval:
		f = ProcessSetAccountPending
	case model.PendingATO:
		f = ProcessSetAccountPending
	case model.Authorized:
		f = ProcessSetAccountActive
	case model.UnauthorizedDenied:
		f = ProcessSetAccountInactive
	case model.UnauthorizedATO:
		f = ProcessSetAccountInactive
	case model.UnauthorizedOptOut:
		f = ProcessSetAccountInactive
	case model.UnauthorizedSecurityRisk:
		f = ProcessSetAccountInactive
	case model.UnauthorizedROB:
		f = ProcessSetAccountInactive
	default:
		return fmt.Errorf("unknown status: %s", status)
	}

	return f(ctx, pvtColl, accountName, store)
}

func ProcessSetAccountActive(ctx contractapi.TransactionContextInterface, pvtCollName, account string, store policy.Store) error {
	user, err := common.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("error getting user from stub: %w", err)
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "set_account_active",
		Args: map[string]string{
			"accountName": account,
		},
	}

	return process(ctx, evtCtx, store)
}

func ProcessSetAccountPending(ctx contractapi.TransactionContextInterface, pvtCollName, account string, store policy.Store) error {
	user, err := common.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("error getting user from stub: %w", err)
	}

	policyStore, err := common.GetPvtCollPolicyStore(ctx, pvtCollName)
	if err != nil {
		return fmt.Errorf("error getting ngac components: %w", err)
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "set_account_pending",
		Args: map[string]string{
			"accountName": account,
		},
	}

	return process(ctx, evtCtx, policyStore)
}

func ProcessSetAccountInactive(ctx contractapi.TransactionContextInterface, pvtCollName, account string, store policy.Store) error {
	user, err := common.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("error getting user from stub: %w", err)
	}

	policyStore, err := common.GetPvtCollPolicyStore(ctx, pvtCollName)
	if err != nil {
		return fmt.Errorf("error getting ngac components: %w", err)
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "set_account_inactive",
		Args: map[string]string{
			"accountName": account,
		},
	}

	return process(ctx, evtCtx, policyStore)
}

func ProcessOnboardAsset(ctx contractapi.TransactionContextInterface, pvtCollName, assetID string) error {
	user, err := common.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("error getting user from stub: %w", err)
	}

	policyStore, err := common.GetPvtCollPolicyStore(ctx, pvtCollName)
	if err != nil {
		return fmt.Errorf("error getting ngac components: %w", err)
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "onboard_asset",
		Args: map[string]string{
			"asset_id": assetID,
		},
	}

	return process(ctx, evtCtx, policyStore)
}

func ProcessOffboardAsset(ctx contractapi.TransactionContextInterface, pvtCollName, assetID string) error {
	user, err := common.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("error getting user from stub: %w", err)
	}

	policyStore, err := common.GetPvtCollPolicyStore(ctx, pvtCollName)
	if err != nil {
		return fmt.Errorf("error getting ngac components: %w", err)
	}

	evtCtx := epp.EventContext{
		User:  user,
		Event: "offboard_asset",
		Args: map[string]string{
			"asset_id": assetID,
		},
	}

	return process(ctx, evtCtx, policyStore)
}
