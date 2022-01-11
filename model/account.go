package model

import (
	"fmt"
)

type (
	AccountPrivate struct {
		ATO    string                       `json:"ato"`
		Users  Users                        `json:"users"`
		Assets map[string]map[string]string `json:"assets" json:"assets"`
	}

	AccountPublic struct {
		Name   string `json:"name"`
		MSPID  string `json:"mspid"`
		Status Status `json:"status"`
	}

	Account struct {
		Name   string                       `json:"name"`
		MSPID  string                       `json:"mspid"`
		Status Status                       `json:"status"`
		ATO    string                       `json:"ato"`
		Users  Users                        `json:"users"`
		Assets map[string]map[string]string `json:"assets" json:"assets"`
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
		"PENDING_APPROVAL":           PendingApproval,
		"PENDING_ATO":                PendingATO,
		"AUTHORIZED":                 Authorized,
		"UNAUTHORIZED_DENIED":        UnauthorizedDenied,
		"UNAUTHORIZED_ATO":           UnauthorizedATO,
		"UNAUTHORIZED_OPTOUT":        UnauthorizedOptOut,
		"UNAUTHORIZED_SECURITY_RISK": UnauthorizedSecurityRisk,
		"UNAUTHORIZED_ROB":           UnauthorizedROB,
	}
)

const (
	PendingApproval          Status = "Pending: waiting for approval"
	PendingATO               Status = "Pending: waiting for ATO"
	Authorized               Status = "Authorized"
	UnauthorizedDenied       Status = "Unauthorized: request denied"
	UnauthorizedATO          Status = "Unauthorized: waiting for ATO renewal"
	UnauthorizedOptOut       Status = "Unauthorized: opted out"
	UnauthorizedSecurityRisk Status = "Unauthorized: security risk"
	UnauthorizedROB          Status = "Unauthorized: breach in rules of behavior"

	AccountPrefix = "account:"
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
	return fmt.Sprintf("%s%s", AccountPrefix, name)
}
