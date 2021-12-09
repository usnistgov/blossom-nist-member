package pap

import "fmt"

const (
	BlossomObject = "blossom_object"
	Assets        = "assets"
)

func AccountObjectName(accountName string) string {
	return fmt.Sprintf("%s_object", accountName)
}
