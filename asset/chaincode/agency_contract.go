package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/asset/chaincode/ngac"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type (
	// AgencyInterface provides the functions to interact with Agencies in blossom.
	AgencyInterface interface {
		// RequestAccount allows agencies to request an account in the Blossom system.  This function will stage the information
		// provided in the Agency parameter in a separate structure until the request is accepted or denied.  The agency will
		// be identified by the name provided in the request. The MSPID of the agency is needed to distinguish users, who may have
		// the same username in a differing MSPs, in the NGAC system.
		RequestAccount(ctx contractapi.TransactionContextInterface, agency Agency) error

		// UploadATO updates the ATO field of the Agency with the given name.
		// TODO placeholder function until ATO model is finalized
		UploadATO(ctx contractapi.TransactionContextInterface, agency string, ato string) error

		// UpdateAgencyStatus updates the status of an agency in Blossom.
		UpdateAgencyStatus(ctx contractapi.TransactionContextInterface, agency string, status Status) error

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
		Agencies(ctx contractapi.TransactionContextInterface) ([]*Agency, error)

		// Agency returns the agency information of the agency with the provided name.  Any fields of any agency the user
		// does not have access to will not be returned.
		Agency(ctx contractapi.TransactionContextInterface, agency string) (*Agency, error)
	}
)

const AgencyPrefix = "agency:"

func NewAgencyContract() AgencyInterface {
	return &BlossomContract{}
}

// AgencyKey returns the key for an agency on the ledger.  Agencies are stored with the format: "agency:<agency_name>".
func AgencyKey(name string) string {
	return fmt.Sprintf("%s%s", AgencyPrefix, name)
}

func (b *BlossomContract) agencyExists(ctx contractapi.TransactionContextInterface, agencyName string) (bool, error) {
	data, err := ctx.GetStub().GetState(AgencyKey(agencyName))
	if err != nil {
		return false, errors.Wrapf(err, "error checking if agency %q aleady exists on the ledger", agencyName)
	}

	return data != nil, nil
}

func (b *BlossomContract) RequestAccount(ctx contractapi.TransactionContextInterface, agency Agency) error {
	// check that an agency doesn't already exist with the same name
	if ok, err := b.agencyExists(ctx, agency.Name); err != nil {
		return errors.Wrapf(err, "error requesting account")
	} else if ok {
		return errors.Errorf("an agency with the name %q already exists", agency.Name)
	}

	// check for permissions and add agency to NGAC
	agencyPEP := new(ngac.AgencyPEP)
	if err := agencyPEP.RequestAccount(ctx, agency); err != nil {
		return errors.Wrapf(err, "error adding agency to NGAC")
	}

	// add agency to ledger with pending status
	agency.Status = PendingApproval

	// convert agency to bytes
	bytes, err := json.Marshal(agency)
	if err != nil {
		return errors.Wrapf(err, "error marshaling agency %q", agency.Name)
	}

	// add agency to world state
	if err = ctx.GetStub().PutState(AgencyKey(agency.Name), bytes); err != nil {
		return errors.Wrapf(err, "error adding agency to ledger")
	}

	return nil
}

func (b *BlossomContract) UploadATO(ctx contractapi.TransactionContextInterface, agencyName string, ato string) error {
	if ok, err := b.agencyExists(ctx, agencyName); err != nil {
		return errors.Wrapf(err, "error checking if agency %q exists", agencyName)
	} else if !ok {
		return errors.Errorf("an agency with the name %q does not exist", agencyName)
	}

	// check permissions in NGAC
	agencyPEP := new(ngac.AgencyPEP)
	if err := agencyPEP.UploadATO(ctx, agencyName); err != nil {
		return errors.Wrapf(err, "error checking if user can update ATO")
	}

	bytes, err := ctx.GetStub().GetState(AgencyKey(agencyName))
	if err != nil {
		return errors.Wrapf(err, "error getting agency %q from world state", agencyName)
	}

	agency := &Agency{}
	if err = json.Unmarshal(bytes, agency); err != nil {
		return errors.Wrapf(err, "error unmarshaling agency %q", agency)
	}

	// update ATO value
	agency.ATO = ato

	// marshal back to json
	if bytes, err = json.Marshal(agency); err != nil {
		return errors.Wrapf(err, "error marshaling agency %q", agency)
	}

	// update world state
	if err = ctx.GetStub().PutState(AgencyKey(agencyName), bytes); err != nil {
		return errors.Wrapf(err, "error updating ATO for agency %q", agencyName)
	}

	return nil
}

func (b *BlossomContract) UpdateAgencyStatus(ctx contractapi.TransactionContextInterface, agencyName string, status Status) error {
	if ok, err := b.agencyExists(ctx, agencyName); err != nil {
		return errors.Wrapf(err, "error checking if agency %q exists", agencyName)
	} else if !ok {
		return errors.Errorf("an agency with the name %q does not exist", agencyName)
	}

	// check permissions in NGAC
	agencyPEP := new(ngac.AgencyPEP)
	if err := agencyPEP.UpdateAgencyStatus(ctx, agencyName); err != nil {
		return errors.Wrapf(err, "error checking if user can update agency status")
	}

	bytes, err := ctx.GetStub().GetState(AgencyKey(agencyName))
	if err != nil {
		return errors.Wrapf(err, "error getting agency %q from world state", agencyName)
	}

	agency := &Agency{}
	if err = json.Unmarshal(bytes, agency); err != nil {
		return errors.Wrapf(err, "error unmarshaling agency %q", agency)
	}

	// update ATO value
	agency.Status = status

	// marshal back to json
	if bytes, err = json.Marshal(agency); err != nil {
		return errors.Wrapf(err, "error marshaling agency %q", agency)
	}

	// update world state
	if err = ctx.GetStub().PutState(AgencyKey(agencyName), bytes); err != nil {
		return errors.Wrapf(err, "error updating ATO for agency %q", agencyName)
	}

	return nil
}

func (b *BlossomContract) ApproveAccountRequest(ctx contractapi.TransactionContextInterface, agency string) error {
	// check that an agency doesn't already exist with the same name
	if ok, err := b.agencyExists(ctx, agency); err != nil {
		return fmt.Errorf("error checking if agency %q exists: %w", agency, err)
	} else if !ok {
		return fmt.Errorf("agency %q does not exist", agency)
	}

	// get agency from ledger
	bytes, err := ctx.GetStub().GetState(AgencyKey(agency))
	if err != nil {
		return fmt.Errorf("error retrieving agency %q from world state: %w", agency, err)
	}

	ledgerAgency := &Agency{}
	if err = json.Unmarshal(bytes, ledgerAgency); err != nil {
		return fmt.Errorf("error deserializing agency %q: %w", agency, err)
	}

	// update agency status
	ledgerAgency.Status = Approved

	// serialize the agency
	if bytes, err = json.Marshal(agency); err != nil {
		return fmt.Errorf("error marshaling agency %q: %w", agency, err)
	}

	// add agency to world state
	if err = ctx.GetStub().PutState(AgencyKey(agency), bytes); err != nil {
		return fmt.Errorf("error updating agency %q in world state: %w", agency, err)
	}

	// add agency to NGAC
	err = ngacApproveAgency(ctx, ledgerAgency)
	if err != nil {
		return fmt.Errorf("error adding agency to NGAC: %w", err)
	}

	return nil
}

func (b *BlossomContract) DenyAccountRequest(ctx contractapi.TransactionContextInterface, agency string) error {
	// check that an agency doesn't already exist with the same name
	if ok, err := b.agencyExists(ctx, agency); err != nil {
		return fmt.Errorf("error checking if agency %q exists: %w", agency, err)
	} else if !ok {
		return fmt.Errorf("agency %q does not exist", agency)
	}

	// get agency from ledger
	bytes, err := ctx.GetStub().GetState(AgencyKey(agency))
	if err != nil {
		return fmt.Errorf("error retrieving agency %q from world state: %w", agency, err)
	}

	ledgerAgency := &Agency{}
	if err = json.Unmarshal(bytes, ledgerAgency); err != nil {
		return fmt.Errorf("error deserializing agency %q: %w", agency, err)
	}

	// update agency status
	ledgerAgency.Status = PendingDenied

	// serialize the agency
	if bytes, err = json.Marshal(agency); err != nil {
		return fmt.Errorf("error marshaling agency %q: %w", agency, err)
	}

	// add agency to world state
	if err = ctx.GetStub().PutState(AgencyKey(agency), bytes); err != nil {
		return fmt.Errorf("error updating agency %q in world state: %w", agency, err)
	}

	return nil
}

func (b *BlossomContract) Agencies(ctx contractapi.TransactionContextInterface) ([]*Agency, error) {
	// add agency to NGAC
	agencyNames, err := ngacAgencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("error adding agency to NGAC: %w", err)
	}

	agencies := make([]*Agency, 0)
	for _, agencyName := range agencyNames {
		bytes, err := ctx.GetStub().GetState(AgencyKey(agencyName))
		if err != nil {
			return nil, fmt.Errorf("error getting agency %q from ledger: %w", agencyName, err)
		}

		agency := &Agency{}
		if err = json.Unmarshal(bytes, agency); err != nil {
			return nil, fmt.Errorf("error deserializing agency %q: %w", agencyName, err)
		}

		agencies = append(agencies, agency)
	}

	return agencies, nil
}

func (b *BlossomContract) Agency(ctx contractapi.TransactionContextInterface, agency string) (*Agency, error) {
	if ok, err := b.agencyExists(ctx, agency); err != nil {
		return &Agency{}, fmt.Errorf("error checking if agency %q exists: %w", agency, err)
	} else if !ok {
		return &Agency{}, fmt.Errorf("agency %q does not exist", agency)
	}

	if ok, err := ngacAgency(ctx, agency); err != nil {
		return &Agency{}, fmt.Errorf("error getting agency from NGAC: %w", err)
	} else if !ok {
		return &Agency{}, fmt.Errorf("agency %q does not exist", agency)
	}

	bytes, err := ctx.GetStub().GetState(AgencyKey(agency))
	if err != nil {
		return &Agency{}, fmt.Errorf("error getting agency %q from ledger: %w", agency, err)
	}

	ledgerAgency := &Agency{}
	if err = json.Unmarshal(bytes, ledgerAgency); err != nil {
		return &Agency{}, fmt.Errorf("error deserializing agency %q: %w", agency, err)
	}

	return ledgerAgency, nil
}
