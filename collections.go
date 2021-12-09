package main

import "fmt"

func AccountCollectionName(account string) string {
	return fmt.Sprintf("%s_account_coll", account)
}

func CatalogCollectionName() string {
	return fmt.Sprintf("catalog_coll")
}

func LicensesCollectionName() string {
	return fmt.Sprintf("licenses_coll")
}
