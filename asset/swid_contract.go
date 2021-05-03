package asset

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

type (
	// SwIDInterface provides the functions to interact with SwID tags in fabric.
	SwIDInterface interface {
		// ReportSwID is used by Agencies to report to Blossom when a software user has installed a piece of software associated
		// with a license that agency has out. This function will invoke NGAc chaincode to add the SwID to the NGAC graph.
		ReportSwID(ctx contractapi.TransactionContextInterface, swid *SwID) error

		// GetSwID returns the SwID object including the XML that matches the provided primaryTag parameter.
		GetSwID(ctx contractapi.TransactionContextInterface, primaryTag string) (*SwID, error)

		// GetLicenseSwIDs returns the primary tags of the SwIDs that are associated with the given license ID.
		GetLicenseSwIDs(ctx contractapi.TransactionContextInterface) ([]string, error)
	}
)

func NewSwIDContract() SwIDInterface {
	return &BlossomContract{}
}

func (b *BlossomContract) ReportSwID(ctx contractapi.TransactionContextInterface, swid *SwID) error {
	return nil
}

func (b *BlossomContract) GetSwID(ctx contractapi.TransactionContextInterface, primaryTag string) (*SwID, error) {
	return &SwID{}, nil
}

func (b *BlossomContract) GetLicenseSwIDs(ctx contractapi.TransactionContextInterface) ([]string, error) {
	return nil, nil
}
