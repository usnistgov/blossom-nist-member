package main

import (
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
)

type BlossomSmartContract struct {
}

func (b *BlossomSmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	adminPDP := pdp.NewAdminDecider()
	if err := adminPDP.InitGraph(stub); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (b *BlossomSmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, _ := stub.GetFunctionAndParameters()
	args := stub.GetArgs()

	var (
		result []byte
		err    error
	)

	switch fn {
	case "RequestAccount":
		err = b.handleRequestAccount(stub, args)
	case "UploadATO":
		err = b.handleUploadATO(stub, args)
	case "UpdateAgencyStatus":
		err = b.handleUpdateAgencyStatus(stub, args)
	case "Agencies":
		result, err = b.handleAgencies(stub)
	case "Agency":
		result, err = b.handleAgency(stub, args)
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
		result, err = b.handlegetSwIDsAssociatedWithAsset(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(result)
}

func (b *BlossomSmartContract) handleRequestAccount(stub shim.ChaincodeStubInterface, args [][]byte) error {
	agency := &model.Agency{}
	if err := json.Unmarshal(args[0], agency); err != nil {
		return err
	}

	return b.RequestAccount(stub, agency)
}

func (b *BlossomSmartContract) handleUploadATO(stub shim.ChaincodeStubInterface, args [][]byte) error {
	agencyName := string(args[0])
	ato := string(args[1])

	return b.UploadATO(stub, agencyName, ato)
}

func (b *BlossomSmartContract) handleUpdateAgencyStatus(stub shim.ChaincodeStubInterface, args [][]byte) error {
	agencyName := string(args[0])
	status := model.Status(args[1])

	return b.UpdateAgencyStatus(stub, agencyName, status)
}

func (b *BlossomSmartContract) handleAgencies(stub shim.ChaincodeStubInterface) ([]byte, error) {
	agencies, err := b.Agencies(stub)
	if err != nil {
		return nil, err
	}

	return json.Marshal(agencies)
}

func (b *BlossomSmartContract) handleAgency(stub shim.ChaincodeStubInterface, args [][]byte) ([]byte, error) {
	agencyName := string(args[0])

	agency, err := b.Agency(stub, agencyName)
	if err != nil {
		return nil, err
	}

	return json.Marshal(agency)
}

func (b *BlossomSmartContract) handleOnboardAsset(stub shim.ChaincodeStubInterface, args [][]byte) error {
	asset := &model.Asset{}
	if err := json.Unmarshal(args[0], asset); err != nil {
		return err
	}

	return b.OnboardAsset(stub, asset)
}

func (b *BlossomSmartContract) handleOffboardAsset(stub shim.ChaincodeStubInterface, args [][]byte) error {
	assetID := string(args[0])
	return b.OffboardAsset(stub, assetID)
}

func (b *BlossomSmartContract) handleAssets(stub shim.ChaincodeStubInterface) ([]byte, error) {
	assets, err := b.Assets(stub)
	if err != nil {
		return nil, err
	}

	return json.Marshal(assets)
}

func (b *BlossomSmartContract) handleAssetInfo(stub shim.ChaincodeStubInterface, args [][]byte) ([]byte, error) {
	assetID := string(args[0])

	asset, err := b.AssetInfo(stub, assetID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(asset)
}

func (b *BlossomSmartContract) handleCheckout(stub shim.ChaincodeStubInterface, args [][]byte) ([]byte, error) {
	assetID := string(args[0])
	agencyName := string(args[1])
	amount, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return nil, err
	}

	result, err := b.Checkout(stub, assetID, agencyName, amount)
	if err != nil {
		return nil, err
	}

	return json.Marshal(result)
}

func (b *BlossomSmartContract) handleCheckin(stub shim.ChaincodeStubInterface, args [][]byte) error {
	assetID := string(args[0])
	licenses := make([]string, 0)
	if err := json.Unmarshal(args[1], &licenses); err != nil {
		return err
	}
	agencyName := string(args[2])

	return b.Checkin(stub, assetID, licenses, agencyName)
}

func (b *BlossomSmartContract) handleReportSwID(stub shim.ChaincodeStubInterface, args [][]byte) error {
	swid := &model.SwID{}
	if err := json.Unmarshal(args[0], swid); err != nil {
		return err
	}
	agencyName := string(args[1])

	return b.ReportSwID(stub, swid, agencyName)
}

func (b *BlossomSmartContract) handleGetSwID(stub shim.ChaincodeStubInterface, args [][]byte) ([]byte, error) {
	primaryTag := string(args[0])
	swid, err := b.GetSwID(stub, primaryTag)
	if err != nil {
		return nil, err
	}

	return json.Marshal(swid)
}

func (b *BlossomSmartContract) handlegetSwIDsAssociatedWithAsset(stub shim.ChaincodeStubInterface, args [][]byte) ([]byte, error) {
	asset := string(args[0])
	swids, err := b.GetSwIDsAssociatedWithAsset(stub, asset)
	if err != nil {
		return nil, err
	}

	return json.Marshal(swids)
}
