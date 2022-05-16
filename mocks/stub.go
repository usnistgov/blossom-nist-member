package mocks

import (
	"encoding/json"
	"github.com/PM-Master/policy-machine-go/policy"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/usnistgov/blossom/chaincode/ngac/common"
)

type (
	stub struct {
		state     map[string][]byte
		args      [][]byte
		function  string
		user      *ClientIdentity
		iterator  []*kv
		transient map[string][]byte
		pvtData   *PvtData
	}

	kv struct {
		k string
		v []byte
	}
)

func (s *stub) HasNext() bool {
	return len(s.iterator) != 0
}

func (s *stub) Close() error {
	s.iterator = nil
	return nil
}

func (s *stub) Next() (*queryresult.KV, error) {
	kv := s.iterator[0]
	s.iterator = s.iterator[1:]
	return &queryresult.KV{Key: kv.k, Value: kv.v}, nil
}

func newStub() *stub {
	return &stub{
		state:     make(map[string][]byte),
		args:      make([][]byte, 0),
		function:  "",
		user:      &ClientIdentity{},
		transient: make(map[string][]byte),
		pvtData:   NewPvtData(),
	}
}

func (s *stub) PutNGAC(collection string, policyStore policy.Store) error {
	bytes, err := policyStore.Graph().MarshalJSON()
	if err != nil {
		return err
	}

	if err = s.PutPrivateData(collection, common.GraphKey, bytes); err != nil {
		return err
	}

	bytes, err = policyStore.Prohibitions().MarshalJSON()
	if err != nil {
		return err
	}

	if err = s.PutPrivateData(collection, common.ProhibitionsKey, bytes); err != nil {
		return err
	}

	bytes, err = policyStore.Obligations().MarshalJSON()
	if err != nil {
		return err
	}

	if err = s.PutPrivateData(collection, common.ObligationsKey, bytes); err != nil {
		return err
	}

	return nil
}

func (s *stub) SetFunctionAndArgs(function string, args ...string) {
	s.function = function

	bytes := make([][]byte, 0)
	for _, a := range args {
		bytes = append(bytes, []byte(a))
	}

	s.args = bytes
}

func (s *stub) SetClientIdentity(clientIdentity *ClientIdentity) {
	s.user = clientIdentity
}

func (s *stub) SetTransient(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	s.transient[key] = bytes

	return nil
}

func (s *stub) GetArgs() [][]byte {
	return s.args
}

func (s *stub) GetStringArgs() []string {
	strArgs := make([]string, 0)
	for _, b := range s.args {
		strArgs = append(strArgs, string(b))
	}
	return strArgs
}

func (s *stub) GetFunctionAndParameters() (string, []string) {
	return s.function, s.GetStringArgs()
}

func (s *stub) GetArgsSlice() ([]byte, error) {
	panic("implement me")
}

func (s *stub) GetTxID() string {
	panic("implement me")
}

func (s *stub) GetChannelID() string {
	panic("implement me")
}

func (s *stub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) peer.Response {
	panic("implement me")
}

func (s *stub) GetState(key string) ([]byte, error) {
	return s.state[key], nil
}

func (s *stub) PutState(key string, value []byte) error {
	s.state[key] = value
	return nil
}

func (s *stub) DelState(key string) error {
	delete(s.state, key)
	return nil
}

func (s *stub) SetStateValidationParameter(key string, ep []byte) error {
	panic("implement me")
}

func (s *stub) GetStateValidationParameter(key string) ([]byte, error) {
	panic("implement me")
}

func (s *stub) GetStateByRange(startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	kvs := make([]*kv, 0)
	for k, v := range s.state {
		kvs = append(kvs, &kv{k, v})
	}
	s.iterator = kvs
	return s, nil
}

func (s *stub) GetStateByRangeWithPagination(startKey, endKey string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *peer.QueryResponseMetadata, error) {
	panic("implement me")
}

func (s *stub) GetStateByPartialCompositeKey(objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
	panic("implement me")
}

func (s *stub) GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *peer.QueryResponseMetadata, error) {
	panic("implement me")
}

func (s *stub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	panic("implement me")
}

func (s *stub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	panic("implement me")
}

func (s *stub) GetQueryResult(query string) (shim.StateQueryIteratorInterface, error) {
	panic("implement me")
}

func (s *stub) GetQueryResultWithPagination(query string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *peer.QueryResponseMetadata, error) {
	panic("implement me")
}

func (s *stub) GetHistoryForKey(key string) (shim.HistoryQueryIteratorInterface, error) {
	panic("implement me")
}

func (s *stub) CreateCollection(collection string, readers []string, writers []string) {
	s.pvtData.CreateNewCollection(collection, readers, writers)
}

func (s *stub) GetPrivateData(collection, key string) ([]byte, error) {
	mspid, err := s.user.GetMSPID()
	if err != nil {
		return nil, err
	}
	return s.pvtData.GetPrivateData(mspid, collection, key)
}

func (s *stub) GetPrivateDataHash(collection, key string) ([]byte, error) {
	panic("implement me")
}

func (s *stub) PutPrivateData(collection string, key string, value []byte) error {
	mspid, err := s.user.GetMSPID()
	if err != nil {
		return err
	}
	return s.pvtData.PutPrivateData(mspid, collection, key, value)
}

func (s *stub) DelPrivateData(collection, key string) error {
	mspid, err := s.user.GetMSPID()
	if err != nil {
		return err
	}

	return s.pvtData.DelPrivateData(mspid, collection, key)
}

func (s *stub) SetPrivateDataValidationParameter(collection, key string, ep []byte) error {
	panic("implement me")
}

func (s *stub) GetPrivateDataValidationParameter(collection, key string) ([]byte, error) {
	panic("implement me")
}

func (s *stub) GetPrivateDataByRange(collection, startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	mspid, err := s.user.GetMSPID()
	if err != nil {
		return nil, err
	}
	return s.pvtData.GetPrivateDataByRange(mspid, collection, startKey, endKey)
}

func (s *stub) GetPrivateDataByPartialCompositeKey(collection, objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
	panic("implement me")
}

func (s *stub) GetPrivateDataQueryResult(collection, query string) (shim.StateQueryIteratorInterface, error) {
	panic("implement me")
}

func (s *stub) GetCreator() ([]byte, error) {
	return nil, nil
}

func (s *stub) GetTransient() (map[string][]byte, error) {
	return s.transient, nil
}

func (s *stub) GetBinding() ([]byte, error) {
	panic("implement me")
}

func (s *stub) GetDecorations() map[string][]byte {
	panic("implement me")
}

func (s *stub) GetSignedProposal() (*peer.SignedProposal, error) {
	panic("implement me")
}

func (s *stub) GetTxTimestamp() (*timestamp.Timestamp, error) {
	panic("implement me")
}

func (s *stub) SetEvent(name string, payload []byte) error {
	panic("implement me")
}
