package main

import "fmt"

func AccountCollection(account string) string {
	return fmt.Sprintf("%s_account_coll", account)
}

func CatalogCollection() string {
	return fmt.Sprintf("catalog_coll")
}

func LicensesCollection() string {
	return fmt.Sprintf("licenses_coll")
}
