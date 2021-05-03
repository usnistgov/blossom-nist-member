package asset

import "github.com/golang/protobuf/ptypes/timestamp"

type (
	// License represents a license on the ledger.
	License struct {
		// ID is the unique identifier of the license
		ID string
		// Name is the common name of the license
		Name string
		// TotalAmount is the total number of keys available to Blossom
		TotalAmount int
		// Available is the number of license keys that are currently available to be checked out
		Available int
		// Cost is the cost of obtaining the license
		Cost float64
		// OnboardingDate is the date in which the license was added to Blossom
		OnboardingDate timestamp.Timestamp
		// Expiration is the date in which the license will expire from Blossom
		Expiration timestamp.Timestamp
		// Keys is a set of license keys associated with this license
		Keys []string
	}
)
