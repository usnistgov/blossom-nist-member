package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type (
	accountTransientInput struct {
		SystemOwner           string `json:"system_owner,omitempty"`
		SystemAdmin           string `json:"system_admin,omitempty"`
		AcquisitionSpecialist string `json:"acquisition_specialist,omitempty"`
	}

	uploadATOTransientInput struct {
		ATO string `json:"ato,omitempty"`
	}

	onboardAssetTransientInput struct {
		Licenses []string `json:"licenses,omitempty"`
	}

	requestCheckoutTransientInput struct {
		AssetID string `json:"asset_id,omitempty"`
		Amount  int    `json:"amount,omitempty"`
	}

	approveCheckoutTransientInput struct {
		Account string `json:"account,omitempty"`
		AssetID string `json:"asset_id,omitempty"`
	}

	initiateCheckinTransientInput struct {
		AssetID  string   `json:"asset_id,omitempty"`
		Licenses []string `json:"licenses,omitempty"`
	}

	processCheckinTransientInput struct {
		Account string `json:"account,omitempty"`
		AssetID string `json:"asset_id,omitempty"`
	}

	reportSwIDTransientInput struct {
		PrimaryTag string `json:"primary_tag,omitempty"`
		Asset      string `json:"asset,omitempty"`
		License    string `json:"license,omitempty"`
		Xml        string `json:"xml,omitempty"`
	}

	swidTransientInput struct {
		Account    string `json:"account,omitempty"`
		PrimaryTag string `json:"primary_tag,omitempty"`
	}

	getSwIDsAssociatedWithAssetTransientInput struct {
		Account string `json:"account,omitempty"`
		AssetID string `json:"asset_id,omitempty"`
	}
)

func getAccountTransientInput(stub shim.ChaincodeStubInterface) (accountTransientInput, error) {
	transientMap, err := stub.GetTransient()
	if err != nil {
		return accountTransientInput{}, fmt.Errorf("error getting transient: %v", err)
	}

	transientAccountJson, ok := transientMap["account"]
	if !ok {
		return accountTransientInput{}, fmt.Errorf("account not found in transient map input")
	}

	var input accountTransientInput
	if err = json.Unmarshal(transientAccountJson, &input); err != nil {
		return accountTransientInput{}, fmt.Errorf("error unmarshaling json: %v", err)
	}

	if len(input.SystemOwner) == 0 {
		return accountTransientInput{}, fmt.Errorf("account system owner cannot be nil")
	}
	if len(input.SystemAdmin) == 0 {
		return accountTransientInput{}, fmt.Errorf("account system admin cannot be nil")
	}
	if len(input.AcquisitionSpecialist) == 0 {
		return accountTransientInput{}, fmt.Errorf("account acquisition specialist cannot be nil")
	}

	return input, nil
}

func getUploadATOTransientInput(stub shim.ChaincodeStubInterface) (uploadATOTransientInput, error) {
	transientMap, err := stub.GetTransient()
	if err != nil {
		return uploadATOTransientInput{}, fmt.Errorf("error getting transient: %v", err)
	}

	transientAccountJson, ok := transientMap["ato"]
	if !ok {
		return uploadATOTransientInput{}, fmt.Errorf("ato not found in transient map input")
	}

	var input uploadATOTransientInput
	if err = json.Unmarshal(transientAccountJson, &input); err != nil {
		return uploadATOTransientInput{}, fmt.Errorf("error unmarshaling json: %v", err)
	}

	if len(input.ATO) == 0 {
		return uploadATOTransientInput{}, fmt.Errorf("ato cannot be nil")
	}

	return input, nil
}

func getOnboardAssetTransientInput(stub shim.ChaincodeStubInterface) (onboardAssetTransientInput, error) {
	transientMap, err := stub.GetTransient()
	if err != nil {
		return onboardAssetTransientInput{}, fmt.Errorf("error getting transient: %v", err)
	}

	transientAccountJson, ok := transientMap["asset"]
	if !ok {
		return onboardAssetTransientInput{}, fmt.Errorf("asset not found in transient map input")
	}

	var input onboardAssetTransientInput
	if err = json.Unmarshal(transientAccountJson, &input); err != nil {
		return onboardAssetTransientInput{}, fmt.Errorf("error unmarshaling json: %v", err)
	}

	if len(input.Licenses) == 0 {
		return onboardAssetTransientInput{}, fmt.Errorf("licenses cannot be empty")
	}

	return input, nil
}

func getRequestCheckoutTransientInput(stub shim.ChaincodeStubInterface) (requestCheckoutTransientInput, error) {
	transientMap, err := stub.GetTransient()
	if err != nil {
		return requestCheckoutTransientInput{}, fmt.Errorf("error getting transient: %v", err)
	}

	transientAccountJson, ok := transientMap["checkout"]
	if !ok {
		return requestCheckoutTransientInput{}, fmt.Errorf("checkout not found in transient map input")
	}

	var input requestCheckoutTransientInput
	if err = json.Unmarshal(transientAccountJson, &input); err != nil {
		return requestCheckoutTransientInput{}, fmt.Errorf("error unmarshaling json: %v", err)
	}

	if input.AssetID == "" {
		return requestCheckoutTransientInput{}, fmt.Errorf("asset id cannot be empty")
	}
	if input.Amount == 0 {
		return requestCheckoutTransientInput{}, fmt.Errorf("amount cannot be empty")
	}

	return input, nil
}

func getApproveCheckoutTransientInput(stub shim.ChaincodeStubInterface) (approveCheckoutTransientInput, error) {
	transientMap, err := stub.GetTransient()
	if err != nil {
		return approveCheckoutTransientInput{}, fmt.Errorf("error getting transient: %v", err)
	}

	transientAccountJson, ok := transientMap["checkout"]
	if !ok {
		return approveCheckoutTransientInput{}, fmt.Errorf("checkout not found in transient map input")
	}

	var input approveCheckoutTransientInput
	if err = json.Unmarshal(transientAccountJson, &input); err != nil {
		return approveCheckoutTransientInput{}, fmt.Errorf("error unmarshaling json: %v", err)
	}

	if input.Account == "" {
		return approveCheckoutTransientInput{}, fmt.Errorf("account cannot be empty")
	}
	if input.AssetID == "" {
		return approveCheckoutTransientInput{}, fmt.Errorf("asset id cannot be empty")
	}

	return input, nil
}

func getInitiateCheckinTransientInput(stub shim.ChaincodeStubInterface) (initiateCheckinTransientInput, error) {
	transientMap, err := stub.GetTransient()
	if err != nil {
		return initiateCheckinTransientInput{}, fmt.Errorf("error getting transient: %v", err)
	}

	transientAccountJson, ok := transientMap["checkin"]
	if !ok {
		return initiateCheckinTransientInput{}, fmt.Errorf("checkin not found in transient map input")
	}

	var input initiateCheckinTransientInput
	if err = json.Unmarshal(transientAccountJson, &input); err != nil {
		return initiateCheckinTransientInput{}, fmt.Errorf("error unmarshaling json: %v", err)
	}

	if input.AssetID == "" {
		return initiateCheckinTransientInput{}, fmt.Errorf("asset id cannot be nil")
	}
	if len(input.Licenses) == 0 {
		return initiateCheckinTransientInput{}, fmt.Errorf("licenses cannot be empty")
	}

	return input, nil
}

func getProcessCheckinTransientInput(stub shim.ChaincodeStubInterface) (processCheckinTransientInput, error) {
	transientMap, err := stub.GetTransient()
	if err != nil {
		return processCheckinTransientInput{}, fmt.Errorf("error getting transient: %v", err)
	}

	transientAccountJson, ok := transientMap["checkin"]
	if !ok {
		return processCheckinTransientInput{}, fmt.Errorf("checkin not found in transient map input")
	}

	var input processCheckinTransientInput
	if err = json.Unmarshal(transientAccountJson, &input); err != nil {
		return processCheckinTransientInput{}, fmt.Errorf("error unmarshaling json: %v", err)
	}

	if input.Account == "" {
		return processCheckinTransientInput{}, fmt.Errorf("account cannot be nil")
	}
	if input.AssetID == "" {
		return processCheckinTransientInput{}, fmt.Errorf("asset id cannot be nil")
	}

	return input, nil
}

func getReportSwIDTransientInput(stub shim.ChaincodeStubInterface) (reportSwIDTransientInput, error) {
	transientMap, err := stub.GetTransient()
	if err != nil {
		return reportSwIDTransientInput{}, fmt.Errorf("error getting transient: %v", err)
	}

	transientAccountJson, ok := transientMap["swid"]
	if !ok {
		return reportSwIDTransientInput{}, fmt.Errorf("swid not found in transient map input")
	}

	var input reportSwIDTransientInput
	if err = json.Unmarshal(transientAccountJson, &input); err != nil {
		return reportSwIDTransientInput{}, fmt.Errorf("error unmarshaling json: %v", err)
	}

	if input.PrimaryTag == "" {
		return reportSwIDTransientInput{}, fmt.Errorf("primary tag cannot be nil")
	}
	if input.Asset == "" {
		return reportSwIDTransientInput{}, fmt.Errorf("asset cannot be nil")
	}
	if input.License == "" {
		return reportSwIDTransientInput{}, fmt.Errorf("license cannot be nil")
	}
	if input.Xml == "" {
		return reportSwIDTransientInput{}, fmt.Errorf("xml cannot be nil")
	}

	return input, nil
}

func getGetSwIDTransientInput(stub shim.ChaincodeStubInterface) (swidTransientInput, error) {
	transientMap, err := stub.GetTransient()
	if err != nil {
		return swidTransientInput{}, fmt.Errorf("error getting transient: %v", err)
	}

	transientAccountJson, ok := transientMap["swid"]
	if !ok {
		return swidTransientInput{}, fmt.Errorf("swid not found in transient map input")
	}

	var input swidTransientInput
	if err = json.Unmarshal(transientAccountJson, &input); err != nil {
		return swidTransientInput{}, fmt.Errorf("error unmarshaling json: %v", err)
	}

	if input.Account == "" {
		return swidTransientInput{}, fmt.Errorf("account cannot be nil")
	}
	if input.PrimaryTag == "" {
		return swidTransientInput{}, fmt.Errorf("primary tag cannot be nil")
	}

	return input, nil
}

func getGetSwIDsAssociatedWithAssetTransientInput(stub shim.ChaincodeStubInterface) (getSwIDsAssociatedWithAssetTransientInput, error) {
	transientMap, err := stub.GetTransient()
	if err != nil {
		return getSwIDsAssociatedWithAssetTransientInput{}, fmt.Errorf("error getting transient: %v", err)
	}

	transientAccountJson, ok := transientMap["swid"]
	if !ok {
		return getSwIDsAssociatedWithAssetTransientInput{}, fmt.Errorf("swid not found in transient map input")
	}

	var input getSwIDsAssociatedWithAssetTransientInput
	if err = json.Unmarshal(transientAccountJson, &input); err != nil {
		return getSwIDsAssociatedWithAssetTransientInput{}, fmt.Errorf("error unmarshaling json: %v", err)
	}

	if input.Account == "" {
		return getSwIDsAssociatedWithAssetTransientInput{}, fmt.Errorf("account cannot be nil")
	}
	if input.AssetID == "" {
		return getSwIDsAssociatedWithAssetTransientInput{}, fmt.Errorf("asset id cannot be nil")
	}

	return input, nil
}
