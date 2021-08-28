package account

import "fmt"

func ObjectAttributeName(accountName string) string {
	return fmt.Sprintf("%s_OA", accountName)
}

func InfoObjectName(accountName string) string {
	return fmt.Sprintf("%s_info", accountName)
}

func UserAttributeName(accountName string) string {
	return fmt.Sprintf("%s_UA", accountName)
}
