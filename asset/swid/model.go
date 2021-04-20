package swid

import "github.com/golang/protobuf/ptypes/timestamp"

type (
	// SwID represents a software identification tag
	SwID struct {
		// PrimaryTag identifies the software asset
		PrimaryTag string
		// XML is the contents of the SwID document in xml format
		XML string
		// License is the associated license
		License string
		// LeaseExpiration is the date when the lease associated with this SwID expires
		LeaseExpiration timestamp.Timestamp
	}
)
