package agency

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type (
	// Contract for managing agencies
	Contract struct {
		contractapi.Contract
	}

	// Interface provides the functions to interact with Agencies in fabric.
	Interface interface {
		// RequestAccount allows agencies to request an account in the Blossom system.  This function will stage the information
		// provided in the Agency parameter in a separate structure until the request is accepted or denied.  The agency will
		// be identified by the name provided in the request. The MSPID of the agency is needed to distinguish users, who may have
		// the same username in a differing MSPs, in the NGAC system.
		RequestAccount(ctx contractapi.TransactionContextInterface, agency Agency) error

		// UploadATO updates the ATO field of the Agency with the given name.
		// TODO placeholder function until ATO model is finalized
		UploadATO(ctx contractapi.TransactionContextInterface, agency string, ato string) error

		// ApproveAccountRequest approves a request matching the provided agency name. Only a user that has permission to accept
		// a request may do so.  This will most likely be a Blossom administrator.
		// Before returning, this function will invoke NGAC chaincode to register the users in the agency request into the
		// NGAC system.
		//   - The SystemOwner will be granted administrative permissions on the agency's account and perform administrative tasks
		//     in the context of the agency.
		//   - The AcquisitionSpecialist will be granted permission to authorize transactions within their agency
		//   - The SystemAdministrator will be granted permission to query the ledger and checkin/checkout software licenses
		ApproveAccountRequest(ctx contractapi.TransactionContextInterface, agency string) error

		// DenyAccountRequest denies a request matching the provided agency name.  Only a user that has permission to deny
		// a request may do so.  This will most likely be a Blossom administrator.
		DenyAccountRequest(ctx contractapi.TransactionContextInterface, agency string) error

		// Agencies returns a list of all the agencies that are registered with Blossom.  Any agency in which the requesting
		// user does not have access to will not be returned.  Likewise, any fields of any agency the user does not have access
		// to will not be returned.
		Agencies(ctx contractapi.TransactionContextInterface) ([]Agency, error)

		// Agency returns the agency information of the agency with the provided name.  Any fields of any agency the user
		// does not have access to will not be returned.
		Agency(ctx contractapi.TransactionContextInterface, agency string) (Agency, error)
	}
)

const AgencyPrefix = "agency:"

func New() Interface {
	return Contract{}
}

// AgencyKey returns the key for an agency on the ledger.  Agencies are stored with the format: "agency:<agency_name>".
func AgencyKey(name string) string {
	return fmt.Sprintf("%s%s", AgencyPrefix, name)
}

func (c *Contract) agencyExists(ctx contractapi.TransactionContextInterface, agency string) (bool, error) {
	data, err := ctx.GetStub().GetState(agency)
	if err != nil {
		return false, fmt.Errorf("error checking if agency %q aleady exists on the ledger: %w", agency, err)
	}

	return data != nil, nil
}

func (c Contract) RequestAccount(ctx contractapi.TransactionContextInterface, agency Agency) error {
	// check that an agency doesn't already exist with the same name
	if ok, err := c.agencyExists(ctx, agency.Name); err != nil {
		return fmt.Errorf("error requesting account: %w", err)
	} else if ok {
		return fmt.Errorf("an agency with the name %q already exists", agency.Name)
	}

	// add agency to ledger with pending status
	agency.Status = Pending

	// convert agency to bytes
	bytes, err := json.Marshal(agency)
	if err != nil {
		return fmt.Errorf("error marshaling agency %q: %w", agency.Name, err)
	}

	// add agency to world state
	if err = ctx.GetStub().PutState(AgencyKey(agency.Name), bytes); err != nil {
		return fmt.Errorf("error adding agency to ledger: %w", err)
	}

	// add agency to NGAC
	err = createAgency(ctx, agency)
	if err != nil {
		return fmt.Errorf("error adding agency to NGAC: %w", err)
	}

	return nil
}

func (c Contract) UploadATO(ctx contractapi.TransactionContextInterface, agency string, ato string) error {
	return nil
}

func (c Contract) ApproveAccountRequest(ctx contractapi.TransactionContextInterface, agency string) error {
	return nil
}

func (c Contract) DenyAccountRequest(ctx contractapi.TransactionContextInterface, agency string) error {
	return nil
}

func (c Contract) Agencies(ctx contractapi.TransactionContextInterface) ([]Agency, error) {
	return nil, nil
}

func (c Contract) Agency(ctx contractapi.TransactionContextInterface, agency string) (Agency, error) {
	return Agency{}, nil
}
