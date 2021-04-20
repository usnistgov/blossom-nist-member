package agency

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

// Contract for managing agencies
type Contract struct {
	contractapi.Contract
}

// RequestAccount allows agencies to request an account in the Blossom system.  This function will stage the information
// provided in the Agency parameter in a separate structure until the request is accepted or denied.  The agency will
// be identified by the name provided in the request. The MSPID of the agency is needed to distinguish users, who may have
// the same username in a differing MSPs, in the NGAC system.
func (c Contract) RequestAccount(ctx contractapi.TransactionContextInterface, agency Agency) error {
	return nil
}

// ApproveAccountRequest approves a request matching the provided agency name. Only a user that has permission to accept
// a request may do so.  This will most likely be a Blossom administrator.
// Before returning, this function will invoke NGAC chaincode to register the users in the agency request into the
// NGAC system.
//   - The SystemOwner will be granted administrative permissions on the agency's account and perform administrative tasks
//     in the context of the agency.
//   - The AcquisitionSpecialist will be granted permission to authorize transactions within their agency
//   - The SystemAdministrator will be granted permission to query the ledger and checkin/checkout software licenses
func (c Contract) ApproveAccountRequest(ctx contractapi.TransactionContextInterface, name string) error {
	return nil
}

// DenyAccountRequest denies a request matching the provided agency name.  Only a user that has permission to deny
// a request may do so.  This will most likely be a Blossom administrator.
func (c Contract) DenyAccountRequest(ctx contractapi.TransactionContextInterface, name string) error {
	return nil
}

// GetAgencies returns a list of all the agencies that are registered with Blossom.  Any agency in which the requesting
// user does not have access to will not be returned.  Likewise, any fields of any agency the user does not have access
// to will not be returned.
func (c Contract) Agencies(ctx contractapi.TransactionContextInterface) ([]Agency, error) {
	return nil, nil
}

// GetAgency returns the agency information of the agency with the provided name.  Any fields of any agency the user
// does not have access to will not be returned.
func (c Contract) Agency(ctx contractapi.TransactionContextInterface, name string) (Agency, error) {
	return Agency{}, nil
}

// TODO placeholder function until ATO model is finalized
func (c Contract) UploadATO(ctx contractapi.TransactionContextInterface, ato string) error {
	return nil
}
