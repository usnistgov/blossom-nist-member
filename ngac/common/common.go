package common

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/PM-Master/policy-machine-go/policy"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/collections"
)

const (
	GraphKey        = "graph"
	ProhibitionsKey = "prohibitions"
	ObligationsKey  = "obligations"
)

func FormatUsername(user string, mspid string) string {
	return fmt.Sprintf("%s:%s", user, mspid)
}

func GetUser(ctx contractapi.TransactionContextInterface) (string, error) {
	cid := ctx.GetClientIdentity()

	cert, err := cid.GetX509Certificate()
	if err != nil {
		return "", err
	}

	mspid, err := cid.GetMSPID()
	if err != nil {
		return "", err
	}

	return FormatUsername(cert.Subject.CommonName, mspid), nil
}

func GetUsername(ctx contractapi.TransactionContextInterface) (string, error) {
	cert, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return "", err
	}

	return cert.Subject.CommonName, nil
}

func GetPvtCollPolicyStore(ctx contractapi.TransactionContextInterface, pvtCollName string) (policy.Store, error) {
	pip := memory.NewPolicyStore()

	// get graph
	bytes, err := ctx.GetStub().GetPrivateData(pvtCollName, GraphKey)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading graph of collection %s", pvtCollName)
	} else if bytes == nil {
		return nil, fmt.Errorf("catalog collection NGAC graph has not been initialized with InitNGAC")
	}

	if err = pip.Graph().UnmarshalJSON(bytes); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling graph bytes")
	}

	// get prohibitions
	bytes, err = ctx.GetStub().GetPrivateData(pvtCollName, ProhibitionsKey)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading graph of collection %s", pvtCollName)
	}
	if bytes != nil {
		if err = pip.Prohibitions().UnmarshalJSON(bytes); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling prohibition bytes")
		}
	}

	// get obligations
	bytes, err = ctx.GetStub().GetPrivateData(pvtCollName, ObligationsKey)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading graph of collection %s", pvtCollName)
	}

	if bytes != nil {
		if err = pip.Obligations().UnmarshalJSON(bytes); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling obligation bytes")
		}
	}

	return pip, nil
}

func PutPvtCollPolicyStore(ctx contractapi.TransactionContextInterface, policyStore policy.Store) error {
	coll := collections.Catalog()

	// put graph
	bytes, err := policyStore.Graph().MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, "error marshaling graph for collection %s", coll)
	}

	if err = ctx.GetStub().PutPrivateData(coll, GraphKey, bytes); err != nil {
		return errors.Wrapf(err, "error putting graph for collection %s", coll)
	}

	// put prohibitions
	bytes, err = policyStore.Prohibitions().MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, "error marshaling graph for collection %s", coll)
	}

	if err = ctx.GetStub().PutPrivateData(coll, ProhibitionsKey, bytes); err != nil {
		return errors.Wrapf(err, "error putting prohibitions for collection %s", coll)
	}

	// put obligations
	bytes, err = policyStore.Obligations().MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, "error marshaling obligations for collection %s", coll)
	}

	if err = ctx.GetStub().PutPrivateData(coll, ObligationsKey, bytes); err != nil {
		return errors.Wrapf(err, "error putting obligations for collection %s", coll)
	}

	return nil
}

func IsNGACInitialized(ctx contractapi.TransactionContextInterface, collName string) (bool, error) {
	bytes, err := ctx.GetStub().GetPrivateData(collName, GraphKey)
	if err != nil {
		return false, errors.Wrapf(err, "error reading graph of collection %s", collName)
	}

	return bytes != nil, nil
}
