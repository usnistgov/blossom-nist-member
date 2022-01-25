package collections

import "fmt"

func Account(account string) string {
	return fmt.Sprintf("%s_account_coll", account)
}

func Catalog() string {
	return fmt.Sprintf("catalog_coll_v2")
}

func Licenses() string {
	return fmt.Sprintf("licenses_coll")
}
