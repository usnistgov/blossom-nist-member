package model

import (
	"fmt"
)

type (
	// SwID represents a software identification tag
	SwID struct {
		// PrimaryTag identifies the software asset
		PrimaryTag string `json:"primary_tag"`
		// XML is the contents of the SwID document in xml format
		XML string `json:"xml"`
		// Asset is the ID of the associated asset
		Asset string `json:"asset"`
		// License is the ID of the associated license
		License string `json:"license"`
	}
)

const SwIDPrefix = "swid:"

// SwIDKey returns the key for a swid tag on the ledger.  SwIDs are stored with the format: "swid:<primary_tag>".
func SwIDKey(name string) string {
	return fmt.Sprintf("%s%s", string(SwIDPrefix), name)
}
