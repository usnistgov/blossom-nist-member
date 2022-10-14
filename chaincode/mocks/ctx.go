package mocks

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type (
	Ctx struct {
		stub shim.ChaincodeStubInterface
		user cid.ClientIdentity
	}
)

func NewCtx() *Ctx {
	return &Ctx{
		stub: newStub(),
		user: &ClientIdentity{},
	}
}

func (c *Ctx) CreateCollection(collection string, readers []string, writers []string) {
	c.stub.(*stub).CreateCollection(collection, readers, writers)
}

func (c *Ctx) SetTransient(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.stub.(*stub).transient[key] = bytes

	return nil
}

func (c *Ctx) GetStub() shim.ChaincodeStubInterface {
	return c.stub
}

func (c *Ctx) GetClientIdentity() cid.ClientIdentity {
	return c.user
}

func (c *Ctx) SetClientIdentity(userFun func() (*ClientIdentity, error)) error {
	clientIdentity, err := userFun()
	if err != nil {
		return fmt.Errorf("error setting mock user: %w", err)
	}

	/*sid := &msp.SerializedIdentity{IdBytes: []byte(cert), Mspid: mspid}
	bytes, err := proto.Marshal(sid)
	if err != nil {
		return errors.Wrap(err, "error marshaling serialized identity")
	}*/

	c.user = clientIdentity
	c.stub.(*stub).SetClientIdentity(clientIdentity)

	return nil
}
