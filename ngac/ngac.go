// Package ngac provides the Policy Decision Point (PDP) and Policy Administration Point (PAP) for Blossom.  The PDP
// provides methods to apply NGAC access control policies to the blossom assets.  The PAP controls the NGAC graph and underlying
// policy that the PDP relies on to make access decision.
package ngac

import (
	"fmt"
)

func FormatUsername(user string, mspid string) string {
	return fmt.Sprintf("%s:%s", user, mspid)
}
