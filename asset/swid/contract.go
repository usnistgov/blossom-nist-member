package swid

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

// Contract to manage SwIDs on the the ledger
type Contract struct {
	contractapi.Contract
}

// ReportSwID is used by Agencies to report to Blossom when a software user has installed a piece of software associated
// with a license that agency has out. This function will invoke NGAc chaincode to add the SwID to the NGAC graph.
func (c Contract) ReportSwID(ctx contractapi.TransactionContextInterface, swid SwID) error {
	return nil
}

// GetSwID returns the SwID object including the XML that matches the provided primaryTag parameter.
func (c Contract) GetSwID(ctx contractapi.TransactionContextInterface, primaryTag string) (SwID, error) {
	return SwID{}, nil
}

// GetLicenseSwIDs returns the SwIDs that are associated with the given license ID.
func (c Contract) GetLicenseSwIDs(ctx contractapi.TransactionContextInterface) ([]SwID, error) {
	return nil, nil
}
