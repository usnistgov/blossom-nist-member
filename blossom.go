package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
)

type BlossomSmartContract struct {
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(BlossomSmartContract)); err != nil {
		fmt.Printf("Error starting Blossom chaincode: %s", err)
	}
}

func (b *BlossomSmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (b *BlossomSmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	var (
		result []byte
		err    error
	)

	if fn != "InitNGAC" && !isNGACInitialized(stub) {
		return shim.Error("ngac not initialized")
	}

	switch fn {
	case "InitNGAC":
		err = b.handleInitNGAC(stub)
	case "RequestAccount":
		err = b.handleRequestAccount(stub, args)
	case "UploadATO":
		err = b.handleUploadATO(stub, args)
	case "UpdateAccountStatus":
		err = b.handleUpdateAccountStatus(stub, args)
	case "Accounts":
		result, err = b.handleAccounts(stub)
	case "Account":
		result, err = b.handleAccount(stub, args)
	case "OnboardAsset":
		err = b.handleOnboardAsset(stub, args)
	case "OffboardAsset":
		err = b.handleOffboardAsset(stub, args)
	case "Assets":
		result, err = b.handleAssets(stub)
	case "AssetInfo":
		result, err = b.handleAssetInfo(stub, args)
	case "Checkout":
		result, err = b.handleCheckout(stub, args)
	case "Checkin":
		err = b.handleCheckin(stub, args)
	case "ReportSwID":
		err = b.handleReportSwID(stub, args)
	case "GetSwID":
		result, err = b.handleGetSwID(stub, args)
	case "GetSwIDsAssociatedWithAsset":
		result, err = b.handleGetSwIDsAssociatedWithAsset(stub, args)
	case "test":
		result = []byte(args[0])
	case "GetHistory":
		result, err = b.handleGetHistory(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(result)
}

func isNGACInitialized(stub shim.ChaincodeStubInterface) bool {
	graphBytes, err := stub.GetState("graph")
	if err != nil {
		return false
	}

	return graphBytes != nil
}

func (b *BlossomSmartContract) handleInitNGAC(stub shim.ChaincodeStubInterface) error {
	adminPDP := pdp.NewAdminDecider()
	if err := adminPDP.InitGraph(stub); err != nil {
		return err
	}

	return nil
}

func (b *BlossomSmartContract) handleGetHistory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("incorrect number of arguments, expecting 1")
	}

	accountName := args[0]

	history, err := b.GetHistory(stub, accountName)
	if err != nil {
		return nil, err
	}

	return json.Marshal(history)
}

func (b *BlossomSmartContract) handleRequestAccount(stub shim.ChaincodeStubInterface, args []string) error {
	account := &model.Account{}
	if err := json.Unmarshal([]byte(args[0]), account); err != nil {
		return err
	}

	return b.RequestAccount(stub, account)
}

func (b *BlossomSmartContract) handleUploadATO(stub shim.ChaincodeStubInterface, args []string) error {
	accountName := args[0]
	ato := args[1]

	return b.UploadATO(stub, accountName, ato)
}

func (b *BlossomSmartContract) handleUpdateAccountStatus(stub shim.ChaincodeStubInterface, args []string) error {
	accountName := args[0]
	status := model.Status(args[1])

	return b.UpdateAccountStatus(stub, accountName, status)
}

func (b *BlossomSmartContract) handleAccounts(stub shim.ChaincodeStubInterface) ([]byte, error) {
	accounts, err := b.Accounts(stub)
	if err != nil {
		return nil, err
	}

	return json.Marshal(accounts)
}

func (b *BlossomSmartContract) handleAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	accountName := args[0]

	account, err := b.Account(stub, accountName)
	if err != nil {
		return nil, err
	}

	return json.Marshal(account)
}

func (b *BlossomSmartContract) handleOnboardAsset(stub shim.ChaincodeStubInterface, args []string) error {
	asset := &model.Asset{}
	if err := json.Unmarshal([]byte(args[0]), asset); err != nil {
		return err
	}

	return b.OnboardAsset(stub, asset)
}

func (b *BlossomSmartContract) handleOffboardAsset(stub shim.ChaincodeStubInterface, args []string) error {
	assetID := args[0]
	return b.OffboardAsset(stub, assetID)
}

func (b *BlossomSmartContract) handleAssets(stub shim.ChaincodeStubInterface) ([]byte, error) {
	assets, err := b.Assets(stub)
	if err != nil {
		return nil, err
	}

	return json.Marshal(assets)
}

func (b *BlossomSmartContract) handleAssetInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	assetID := args[0]

	asset, err := b.AssetInfo(stub, assetID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(asset)
}

func (b *BlossomSmartContract) handleCheckout(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	assetID := args[0]
	accountName := args[1]
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, err
	}

	result, err := b.Checkout(stub, assetID, accountName, amount)
	if err != nil {
		return nil, err
	}

	return json.Marshal(result)
}

func (b *BlossomSmartContract) handleCheckin(stub shim.ChaincodeStubInterface, args []string) error {
	assetID := args[0]
	licenses := make([]string, 0)
	if err := json.Unmarshal([]byte(args[1]), &licenses); err != nil {
		return err
	}
	accountName := args[2]

	return b.Checkin(stub, assetID, licenses, accountName)
}

func (b *BlossomSmartContract) handleReportSwID(stub shim.ChaincodeStubInterface, args []string) error {
	swid := &model.SwID{}
	if err := json.Unmarshal([]byte(args[0]), swid); err != nil {
		return err
	}
	accountName := args[1]

	return b.ReportSwID(stub, swid, accountName)
}

func (b *BlossomSmartContract) handleGetSwID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	primaryTag := args[0]
	swid, err := b.GetSwID(stub, primaryTag)
	if err != nil {
		return nil, err
	}

	return json.Marshal(swid)
}

func (b *BlossomSmartContract) handleGetSwIDsAssociatedWithAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	asset := args[0]
	swids, err := b.GetSwIDsAssociatedWithAsset(stub, asset)
	if err != nil {
		return nil, err
	}

	return json.Marshal(swids)
}
