package mocks

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
)

type MemChaincodeStub struct {
	state    map[string][]byte
	args     [][]byte
	function string
	user     []byte
	iterator []*kv
}

type kv struct {
	k string
	v []byte
}

func (m *MemChaincodeStub) HasNext() bool {
	return len(m.iterator) != 0
}

func (m *MemChaincodeStub) Close() error {
	m.iterator = nil
	return nil
}

func (m *MemChaincodeStub) Next() (*queryresult.KV, error) {
	kv := m.iterator[0]
	m.iterator = m.iterator[1:]
	return &queryresult.KV{Key: kv.k, Value: kv.v}, nil
}

func NewMemCCStub() MemChaincodeStub {
	return MemChaincodeStub{
		state:    make(map[string][]byte),
		args:     make([][]byte, 0),
		function: "",
		user:     make([]byte, 0),
	}
}

func (m *MemChaincodeStub) SetFunctionAndArgs(function string, args ...[]byte) {
	m.function = function
	m.args = args
}

func (m *MemChaincodeStub) SetUser(userFun func() (string, string)) error {
	cert, mspid := userFun()
	sid := &msp.SerializedIdentity{IdBytes: []byte(cert), Mspid: mspid}
	bytes, err := proto.Marshal(sid)
	if err != nil {
		return errors.Wrap(err, "error marshaling serialized identity")
	}

	m.user = bytes

	return nil
}

func (m *MemChaincodeStub) GetArgs() [][]byte {
	return m.args
}

func (m *MemChaincodeStub) GetStringArgs() []string {
	strArgs := make([]string, 0)
	for _, b := range m.args {
		strArgs = append(strArgs, string(b))
	}
	return strArgs
}

func (m *MemChaincodeStub) GetFunctionAndParameters() (string, []string) {
	return m.function, m.GetStringArgs()
}

func (m *MemChaincodeStub) GetArgsSlice() ([]byte, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetTxID() string {
	panic("implement me")
}

func (m *MemChaincodeStub) GetChannelID() string {
	panic("implement me")
}

func (m *MemChaincodeStub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) peer.Response {
	panic("implement me")
}

func (m *MemChaincodeStub) GetState(key string) ([]byte, error) {
	return m.state[key], nil
}

func (m *MemChaincodeStub) PutState(key string, value []byte) error {
	m.state[key] = value
	return nil
}

func (m *MemChaincodeStub) DelState(key string) error {
	delete(m.state, key)
	return nil
}

func (m *MemChaincodeStub) SetStateValidationParameter(key string, ep []byte) error {
	panic("implement me")
}

func (m *MemChaincodeStub) GetStateValidationParameter(key string) ([]byte, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetStateByRange(startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	kvs := make([]*kv, 0)
	for k, v := range m.state {
		kvs = append(kvs, &kv{k, v})
	}
	m.iterator = kvs
	return m, nil
}

func (m *MemChaincodeStub) GetStateByRangeWithPagination(startKey, endKey string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *peer.QueryResponseMetadata, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetStateByPartialCompositeKey(objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *peer.QueryResponseMetadata, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetQueryResult(query string) (shim.StateQueryIteratorInterface, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetQueryResultWithPagination(query string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *peer.QueryResponseMetadata, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetHistoryForKey(key string) (shim.HistoryQueryIteratorInterface, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetPrivateData(collection, key string) ([]byte, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetPrivateDataHash(collection, key string) ([]byte, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) PutPrivateData(collection string, key string, value []byte) error {
	panic("implement me")
}

func (m *MemChaincodeStub) DelPrivateData(collection, key string) error {
	panic("implement me")
}

func (m *MemChaincodeStub) SetPrivateDataValidationParameter(collection, key string, ep []byte) error {
	panic("implement me")
}

func (m *MemChaincodeStub) GetPrivateDataValidationParameter(collection, key string) ([]byte, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetPrivateDataByRange(collection, startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetPrivateDataByPartialCompositeKey(collection, objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetPrivateDataQueryResult(collection, query string) (shim.StateQueryIteratorInterface, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetCreator() ([]byte, error) {
	return m.user, nil
}

func (m *MemChaincodeStub) GetTransient() (map[string][]byte, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetBinding() ([]byte, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetDecorations() map[string][]byte {
	panic("implement me")
}

func (m *MemChaincodeStub) GetSignedProposal() (*peer.SignedProposal, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) GetTxTimestamp() (*timestamp.Timestamp, error) {
	panic("implement me")
}

func (m *MemChaincodeStub) SetEvent(name string, payload []byte) error {
	panic("implement me")
}
