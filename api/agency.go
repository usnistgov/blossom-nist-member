package api

import (
	"encoding/json"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type (
	// AgencyInterface provides the functions to interact with Agencies in blossom.
	AgencyInterface interface {
		// RequestAccount allows agencies to request an account in the Blossom system.  This function will stage the information
		// provided in the Agency parameter in a separate structure until the request is accepted or denied.  The agency will
		// be identified by the name provided in the request. The MSPID of the agency is needed to distinguish users, who may have
		// the same username in a differing MSPs, in the NGAC system.
		RequestAccount(ctx contractapi.TransactionContextInterface, agency model.Agency) error

		// UploadATO updates the ATO field of the Agency with the given name.
		// TODO placeholder function until ATO model is finalized
		UploadATO(ctx contractapi.TransactionContextInterface, agency string, ato string) error

		// UpdateAgencyStatus updates the status of an agency in Blossom.
		// Updating the status to Approved allows the agency to read and write to blossom.
		// Updating the status to Pending allows the agency to read write only agency related information such as ATOs.
		// Updating the status to Inactive provides the same NGAC consequences as Pending
		UpdateAgencyStatus(ctx contractapi.TransactionContextInterface, agency string, status model.Status) error

		// Agencies returns a list of all the agencies that are registered with Blossom.  Any agency in which the requesting
		// user does not have access to will not be returned.  Likewise, any fields of any agency the user does not have access
		// to will not be returned.
		Agencies(ctx contractapi.TransactionContextInterface) ([]*model.Agency, error)

		// Agency returns the agency information of the agency with the provided name.  Any fields of any agency the user
		// does not have access to will not be returned.
		Agency(ctx contractapi.TransactionContextInterface, agency string) (*model.Agency, error)
	}
)

func NewAgencyContract() AgencyInterface {
	return &BlossomSmartContract{}
}

func (b *BlossomSmartContract) agencyExists(ctx contractapi.TransactionContextInterface, agencyName string) (bool, error) {
	data, err := ctx.GetStub().GetState(model.AgencyKey(agencyName))
	if err != nil {
		return false, errors.Wrapf(err, "error checking if agency %q already exists on the ledger", agencyName)
	}

	return data != nil, nil
}

func (b *BlossomSmartContract) RequestAccount(ctx contractapi.TransactionContextInterface, agency model.Agency) error {
	// check that an agency doesn't already exist with the same name
	if ok, err := b.agencyExists(ctx, agency.Name); err != nil {
		return errors.Wrapf(err, "error requesting account")
	} else if ok {
		return errors.Errorf("an agency with the name %q already exists", agency.Name)
	}

	// begin NGAC
	if err := pdp.NewAgencyDecider().RequestAccount(ctx, agency); err != nil {
		return errors.Wrapf(err, "error adding agency to NGAC")
	}
	// end NGAC

	// add agency to ledger with pending status
	agency.Status = model.PendingApproval
	agency.Assets = make(map[string]map[string]time.Time)

	// convert agency to bytes
	bytes, err := json.Marshal(agency)
	if err != nil {
		return errors.Wrapf(err, "error marshaling agency %q", agency.Name)
	}

	// add agency to world state
	if err = ctx.GetStub().PutState(model.AgencyKey(agency.Name), bytes); err != nil {
		return errors.Wrapf(err, "error adding agency to ledger")
	}

	return nil
}

func (b *BlossomSmartContract) UploadATO(ctx contractapi.TransactionContextInterface, agencyName string, ato string) error {
	if ok, err := b.agencyExists(ctx, agencyName); err != nil {
		return errors.Wrapf(err, "error checking if agency %q exists", agencyName)
	} else if !ok {
		return errors.Errorf("an agency with the name %q does not exist", agencyName)
	}

	// begin NGAC
	if err := pdp.NewAgencyDecider().UploadATO(ctx, agencyName); errors.Is(err, pdp.ErrAccessDenied) {
		return err
	} else if err != nil {
		return errors.Wrapf(err, "error checking if user can update ATO")
	}
	// end NGAC

	bytes, err := ctx.GetStub().GetState(model.AgencyKey(agencyName))
	if err != nil {
		return errors.Wrapf(err, "error getting agency %q from world state", agencyName)
	}

	ledgerAgency := &model.Agency{}
	if err = json.Unmarshal(bytes, ledgerAgency); err != nil {
		return errors.Wrapf(err, "error unmarshaling agency %q", agencyName)
	}

	// update ATO value
	ledgerAgency.ATO = ato

	// marshal back to json
	if bytes, err = json.Marshal(ledgerAgency); err != nil {
		return errors.Wrapf(err, "error marshaling agency %q", agencyName)
	}

	// update world state
	if err = ctx.GetStub().PutState(model.AgencyKey(agencyName), bytes); err != nil {
		return errors.Wrapf(err, "error updating ATO for agency %q", agencyName)
	}

	return nil
}

func (b *BlossomSmartContract) UpdateAgencyStatus(ctx contractapi.TransactionContextInterface, agencyName string, status model.Status) error {
	if ok, err := b.agencyExists(ctx, agencyName); err != nil {
		return errors.Wrapf(err, "error checking if agency %q exists", agencyName)
	} else if !ok {
		return errors.Errorf("an agency with the name %q does not exist", agencyName)
	}

	// begin NGAC
	if err := pdp.NewAgencyDecider().UpdateAgencyStatus(ctx, agencyName, status); errors.Is(err, pdp.ErrAccessDenied) {
		return err
	} else if err != nil {
		return errors.Wrapf(err, "error checking if user can update agency status")
	}
	// end NGAC

	bytes, err := ctx.GetStub().GetState(model.AgencyKey(agencyName))
	if err != nil {
		return errors.Wrapf(err, "error getting agency %q from world state", agencyName)
	}

	ledgerAgency := &model.Agency{}
	if err = json.Unmarshal(bytes, ledgerAgency); err != nil {
		return errors.Wrapf(err, "error unmarshaling agency %q", agencyName)
	}

	// update ATO value
	ledgerAgency.Status = status

	// marshal back to json
	if bytes, err = json.Marshal(ledgerAgency); err != nil {
		return errors.Wrapf(err, "error marshaling agency %q", agencyName)
	}

	// update world state
	if err = ctx.GetStub().PutState(model.AgencyKey(agencyName), bytes); err != nil {
		return errors.Wrapf(err, "error updating status of agency %q", agencyName)
	}

	return nil
}

func (b *BlossomSmartContract) Agencies(ctx contractapi.TransactionContextInterface) ([]*model.Agency, error) {
	agencies, err := agencies(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error getting agencies")
	}

	// begin NGAC
	if agencies, err = pdp.NewAgencyDecider().FilterAgencies(ctx, agencies); err != nil {
		return nil, errors.Wrapf(err, "error filtering agencies")
	}
	// end NGAC

	return agencies, nil
}

func agencies(ctx contractapi.TransactionContextInterface) ([]*model.Agency, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	agencies := make([]*model.Agency, 0)
	for resultsIterator.HasNext() {
		var queryResponse *queryresult.KV
		if queryResponse, err = resultsIterator.Next(); err != nil {
			return nil, err
		}

		// agencies on the ledger begin with the agency prefix -- ignore other assets
		if !strings.HasPrefix(queryResponse.Key, model.AgencyPrefix) {
			continue
		}

		agency := &model.Agency{}
		if err = json.Unmarshal(queryResponse.Value, agency); err != nil {
			return nil, err
		}

		agencies = append(agencies, agency)
	}

	return agencies, nil
}

func (b *BlossomSmartContract) Agency(ctx contractapi.TransactionContextInterface, agencyName string) (*model.Agency, error) {
	var (
		agency = &model.Agency{}
		bytes  []byte
		err    error
	)

	if bytes, err = ctx.GetStub().GetState(model.AgencyKey(agencyName)); err != nil {
		return nil, errors.Wrapf(err, "error getting agency from ledger")
	}

	if err = json.Unmarshal(bytes, agency); err != nil {
		return nil, errors.Wrapf(err, "error deserializing agency")
	}

	// begin NGAC
	// filter agency object removing any fields the user does not have access to
	if err = pdp.NewAgencyDecider().FilterAgency(ctx, agency); err != nil {
		return nil, errors.Wrapf(err, "error filtering agency %s", agencyName)
	}
	// end NGAC

	return agency, nil
}
