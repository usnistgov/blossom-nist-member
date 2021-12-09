package pap

import (
	"github.com/PM-Master/policy-machine-go/author"
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/PM-Master/policy-machine-go/pip/memory"
)

func LoadCatalogPolicy() (ngac.FunctionalEntity, error) {
	fe := memory.NewPIP()
	policyAuthor := author.New(fe)
	err := policyAuthor.ReadPAL("ngac/pap/catalog_policy.ngac")
	if err != nil {
		return nil, err
	}

	if err = policyAuthor.Apply(); err != nil {
		return nil, err
	}

	return fe, nil
}

func LoadAccountPolicy() (ngac.FunctionalEntity, error) {
	fe := memory.NewPIP()
	policyAuthor := author.New(fe)
	err := policyAuthor.ReadPAL("ngac/pap/account_policy.ngac")
	if err != nil {
		return nil, err
	}

	if err = policyAuthor.Apply(); err != nil {
		return nil, err
	}

	return fe, nil
}
