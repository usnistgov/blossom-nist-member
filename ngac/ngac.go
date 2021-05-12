package ngac

import (
	"fmt"
)

func FormatUsername(user string, mspid string) string {
	return fmt.Sprintf("%s:%s", user, mspid)
}
