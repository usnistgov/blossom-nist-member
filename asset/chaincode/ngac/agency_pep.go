package ngac

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/asset/chaincode"
)

type AgencyPEP struct {}

func (a AgencyPEP) RequestAccount(ctx contractapi.TransactionContextInterface, agency chaincode.Agency) error {
	bytes := []byte(agency.Name)
	response := ctx.GetStub().InvokeChaincode("ngac", [][]byte{[]byte("RequestAccount"), bytes}, "mychannel")
	if response.Status != 200 {
		return fmt.Errorf("error requesting account: %v", response.Message)
	}

	return nil
}

func (a AgencyPEP) UploadATO(ctx contractapi.TransactionContextInterface, agency string) error {
	bytes := []byte(agency)
	response := ctx.GetStub().InvokeChaincode("ngac", [][]byte{[]byte("UploadATO"), bytes}, "mychannel")
	if response.Status != 200 {
		return fmt.Errorf("error uploading ATO: %v", response.Message)
	}

	return nil
}

func (a AgencyPEP) UpdateAgencyStatus(ctx contractapi.TransactionContextInterface, agency string) error {
	bytes := []byte(agency)
	response := ctx.GetStub().InvokeChaincode("ngac", [][]byte{[]byte("UpdateAgencyStatus"), bytes}, "mychannel")
	if response.Status != 200 {
		return fmt.Errorf("error updating agency status: %v", response.Message)
	}

	return nil
}

func (a AgencyPEP) ApproveAccountRequest(ctx contractapi.TransactionContextInterface, agency string) error {
	bytes := []byte(agency)
	response := ctx.GetStub().InvokeChaincode("ngac", [][]byte{[]byte("ApproveAccountRequest"), bytes}, "mychannel")
	if response.Status != 200 {
		return fmt.Errorf("error approving account request: %v", response.Message)
	}

	return nil
}

func (a AgencyPEP) DenyAccountRequest(ctx contractapi.TransactionContextInterface, agency string) error {
	bytes := []byte(agency)
	response := ctx.GetStub().InvokeChaincode("ngac", [][]byte{[]byte("DenyAccountRequest"), bytes}, "mychannel")
	if response.Status != 200 {
		return fmt.Errorf("error denying account request: %v", response.Message)
	}

	return nil
}

func (a AgencyPEP) Agencies(ctx contractapi.TransactionContextInterface) ([]*chaincode.Agency, error) {
	response := ctx.GetStub().InvokeChaincode("ngac", [][]byte{[]byte("Agencies")}, "mychannel")
	if response.Status != 200 {
		return nil, fmt.Errorf("error getting agencies: %v", response.Message)
	}

	agencyNames := make([]string, 0)
	if err := json.Unmarshal(response.Payload, &agencyNames); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling agencies response")
	}

	agencies := make([]*chaincode.Agency, 0)
	for _, agencyName := range agencyNames {
		bytes, err := ctx.GetStub().GetState(chaincode.AgencyKey(agencyName))
		if err != nil {
			return nil, errors.Wrap(err, "error getting agency %q from world state")
		}

		// TODO filter fields of agency

		agency := &chaincode.Agency{}
		if err = json.Unmarshal(bytes, agency); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling agency")
		}

		agencies = append(agencies, agency)
	}

	return agencies, nil
}

func (a AgencyPEP) Agency(ctx contractapi.TransactionContextInterface, agencyName string) (*chaincode.Agency, error) {
	response := ctx.GetStub().InvokeChaincode("ngac", [][]byte{[]byte("Agency"), []byte(agencyName)}, "mychannel")
	if response.Status != 200 {
		return nil, fmt.Errorf("error getting agency %q: %v", agencyName, response.Message)
	}

	bytes, err := ctx.GetStub().GetState(chaincode.AgencyKey(agencyName))
	if err != nil {
		return nil, errors.Wrap(err, "error getting agency %q from world state")
	}

	agency := &chaincode.Agency{}
	if err = json.Unmarshal(bytes, agency); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling agency")
	}

	return agency, nil
}




