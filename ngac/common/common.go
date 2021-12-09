package common

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/pkg/errors"
)

const (
	GraphKey        = "graph"
	ProhibitionsKey = "prohibitions"
	ObligationsKey  = "obligations"
)

func FormatUsername(user string, mspid string) string {
	return fmt.Sprintf("%s:%s", user, mspid)
}

func GetUser(stub shim.ChaincodeStubInterface) (string, error) {
	cert, err := cid.GetX509Certificate(stub)
	if err != nil {
		return "", err
	}

	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return "", err
	}

	return FormatUsername(cert.Subject.CommonName, mspid), nil
}

func GetPvtCollFunctionalEntity(stub shim.ChaincodeStubInterface, pvtCollName string) (ngac.FunctionalEntity, error) {
	pip := memory.NewPIP()

	// get graph
	bytes, err := stub.GetPrivateData(pvtCollName, GraphKey)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading graph of collection %s", pvtCollName)
	}
	if bytes != nil {
		if err = pip.Graph().UnmarshalJSON(bytes); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling graph bytes")
		}
	}

	// get prohibitions
	bytes, err = stub.GetPrivateData(pvtCollName, ProhibitionsKey)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading graph of collection %s", pvtCollName)
	}
	if bytes != nil {
		if err = pip.Prohibitions().UnmarshalJSON(bytes); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling prohibition bytes")
		}
	}

	// get obligations
	bytes, err = stub.GetPrivateData(pvtCollName, ObligationsKey)
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

func PutPvtCollFunctionalEntity(stub shim.ChaincodeStubInterface, pvtCollName string, fe ngac.FunctionalEntity) error {
	// put graph
	bytes, err := fe.Graph().MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, "error marshaling graph for collection %s", pvtCollName)
	}

	if err = stub.PutPrivateData(pvtCollName, GraphKey, bytes); err != nil {
		return errors.Wrapf(err, "error putting graph for collection %s", pvtCollName)
	}

	// put prohibitions
	bytes, err = fe.Prohibitions().MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, "error marshaling graph for collection %s", pvtCollName)
	}

	if err = stub.PutPrivateData(pvtCollName, ProhibitionsKey, bytes); err != nil {
		return errors.Wrapf(err, "error putting prohibitions for collection %s", pvtCollName)
	}

	// put obligations
	bytes, err = fe.Obligations().MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, "error marshaling obligations for collection %s", pvtCollName)
	}

	if err = stub.PutPrivateData(pvtCollName, ObligationsKey, bytes); err != nil {
		return errors.Wrapf(err, "error putting obligations for collection %s", pvtCollName)
	}

	return nil
}
