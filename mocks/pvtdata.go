package mocks

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"strings"
)

type (
	PvtData struct {
		collections map[string]*PvtDataCollection
	}

	PvtDataCollection struct {
		readers map[string]bool
		writers map[string]bool
		pvtData map[string][]byte
	}
)

func NewPvtData() *PvtData {
	return &PvtData{collections: make(map[string]*PvtDataCollection)}
}

func (p *PvtData) CreateNewCollection(collName string, readers []string, writers []string) {
	readersMap := make(map[string]bool)
	writersMap := make(map[string]bool)
	for _, r := range readers {
		readersMap[r] = true
	}
	for _, w := range writers {
		readersMap[w] = true
		writersMap[w] = true
	}

	p.collections[collName] = &PvtDataCollection{
		readers: readersMap,
		writers: writersMap,
		pvtData: make(map[string][]byte),
	}
}

func (p *PvtData) GetPrivateData(mspid, coll, key string) ([]byte, error) {
	collection, ok := p.collections[coll]
	if !ok {
		return nil, fmt.Errorf("collection %q does not exist", coll)
	} else if !collection.readers[mspid] {
		return nil, fmt.Errorf("msp %q does not have read access to collection %q", mspid, coll)
	}

	return collection.pvtData[key], nil
}

func (p *PvtData) PutPrivateData(mspid, coll, key string, bytes []byte) error {
	collection, ok := p.collections[coll]
	if !ok {
		return fmt.Errorf("collection %q does not exist", coll)
	} else if !collection.writers[mspid] {
		return fmt.Errorf("msp %q does not have write to collection %q", mspid, coll)
	}

	collection.pvtData[key] = bytes
	p.collections[coll] = collection

	return nil
}

func (p *PvtData) DelPrivateData(mspid string, coll string, key string) error {
	collection, ok := p.collections[coll]
	if !ok {
		return fmt.Errorf("collection %q does not exist", coll)
	} else if !collection.writers[mspid] {
		return fmt.Errorf("msp %q does not have write to collection %q", mspid, coll)
	}

	delete(collection.pvtData, key)
	p.collections[coll] = collection

	return nil
}

func (p *PvtData) GetPrivateDataByRange(mspid string, coll string, start string, end string) (shim.StateQueryIteratorInterface, error) {
	collection, ok := p.collections[coll]
	if !ok {
		return nil, fmt.Errorf("collection %q does not exist", coll)
	} else if !collection.readers[mspid] {
		return nil, fmt.Errorf("msp %q does not have read access to collection %q", mspid, coll)
	}

	return newIter(coll, collection.pvtData), nil
}

type iter struct {
	kvs   []*queryresult.KV
	index int
	start string
	end   string
}

func newIter(collection string, data map[string][]byte) *iter {
	kvs := make([]*queryresult.KV, 0)
	for key, value := range data {
		kvs = append(kvs, &queryresult.KV{
			Namespace: collection,
			Key:       key,
			Value:     value,
		})
	}

	return &iter{
		kvs:   kvs,
		index: 0,
	}
}

func (i *iter) HasNext() bool {
	for c := i.index; c < len(i.kvs); c++ {
		kv := i.kvs[c]
		if strings.HasPrefix(kv.Key, i.start) {
			// this isn't correct but works for our use
			return true
		}
	}

	return false
}

func (i *iter) Close() error {
	return nil
}

func (i *iter) Next() (*queryresult.KV, error) {
	kv := i.kvs[i.index]
	for !strings.HasPrefix(kv.Key, i.start) {
		i.index++
		kv = i.kvs[i.index]
	}
	i.index++
	return kv, nil
}
