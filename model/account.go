package model

import (
	"fmt"
)

type (
	// Account stores information for a Blossom account
	Account struct {
		// Name is the unique name of the account
		Name string `json:"name"`
		// ATO is the Authority To Operate document
		ATO string `json:"ato"`
		// MSPID is the Membership Service Provider ID
		MSPID string `json:"mspid"`
		// Users contains the users that will access Blossom
		Users Users `json:"users"`
		// Status of an account within the Blossom system
		Status Status `json:"status"`
		// Assets stores the assets that an account has checked out. The first map key is the asset ID, the second map
		// key is the license ID and the value is the expiration of the license lease.
		Assets map[string]map[string]DateTime `json:"assets"`
	}

	// Status represents the status of an account within the blossom system
	Status string

	// Users that will access blossom on behalf of an account
	Users struct {
		// SystemOwner is responsible for administrative tasks for the account system
		SystemOwner string `json:"system_owner"`
		// AcquisitionSpecialist authorizes transaction requests for the account
		AcquisitionSpecialist string `json:"acquisition_specialist"`
		// SystemAdministrator interacts with the smart contracts to checkin and checkout software licenses for the account
		SystemAdministrator string `json:"system_administrator"`
	}
)

const (
	PendingApproval           Status = "Pending: waiting for approval"
	PendingATO                Status = "Pending: waiting for ATO"
	PendingDenied             Status = "Pending: request denied"
	Approved                  Status = "Approved"
	InactiveATO               Status = "Inactive: waiting for ATO renewal"
	InactiveOptOut            Status = "Inactive: opted out"
	InactiveSecurityRisk      Status = "Inactive: security risk"
	InactiveRulesOfEngagement Status = "Inactive: breach in rules of engagement"

	AccountPrefix = "account:"
)

// AccountKey returns the key for an account on the ledger.  Accounts are stored with the format: "account:<account_name>".
func AccountKey(name string) string {
	return fmt.Sprintf("%s%s", AccountPrefix, name)
}
