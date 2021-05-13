package agency

import "fmt"

func ObjectAttributeName(agencyName string) string {
	return fmt.Sprintf("%s_OA", agencyName)
}

func InfoObjectName(agencyName string) string {
	return fmt.Sprintf("%s_info", agencyName)
}

func UserAttributeName(agencyName string) string {
	return fmt.Sprintf("%s_UA", agencyName)
}
