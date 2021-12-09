package model

import (
	"fmt"
)

type (
	AccountPrivate struct {
		ATO    string                         `json:"ato"`
		Users  Users                          `json:"users"`
		Assets map[string]map[string]DateTime `json:"assets" json:"assets"`
	}

	AccountPublic struct {
		Name   string `json:"name"`
		MSPID  string `json:"mspid"`
		Status Status `json:"status"`
	}

	Account struct {
		Name   string                         `json:"name"`
		MSPID  string                         `json:"mspid"`
		Status Status                         `json:"status"`
		ATO    string                         `json:"ato"`
		Users  Users                          `json:"users"`
		Assets map[string]map[string]DateTime `json:"assets" json:"assets"`
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

var (
	statusUpdates = map[string]Status{
		"PENDING_APPROVAL":       PendingApproval,
		"PENDING_ATO":            PendingATO,
		"ACTIVE":                 Active,
		"INACTIVE_DENIED":        InactiveDenied,
		"INACTIVE_ATO":           InactiveATO,
		"INACTIVE_OPTOUT":        InactiveOptOut,
		"INACTIVE_SECURITY_RISK": InactiveSecurityRisk,
		"INACTIVE_ROB":           InactiveROB,
	}
)

const (
	PendingApproval      Status = "Pending: waiting for approval"
	PendingATO           Status = "Pending: waiting for ATO"
	Active               Status = "Active"
	InactiveDenied       Status = "Inactive: request denied"
	InactiveATO          Status = "Inactive: waiting for ATO renewal"
	InactiveOptOut       Status = "Inactive: opted out"
	InactiveSecurityRisk Status = "Inactive: security risk"
	InactiveROB          Status = "Inactive: breach in rules of behavior"

	AccountPrefix = 'a'
)

func GetStatusUpdate(s string) (Status, error) {
	status, ok := statusUpdates[s]
	if !ok {
		return "", fmt.Errorf("unknown status: %s", s)
	}

	return status, nil
}

// AccountKey returns the key for an account on the ledger.  Accounts are stored with the format: "account:<account_name>".
func AccountKey(name string) string {
	return fmt.Sprintf("%s%s", string(AccountPrefix), name)
}
