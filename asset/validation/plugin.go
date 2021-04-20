package validation

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	validation "github.com/hyperledger/fabric/core/handlers/validation/api"
	. "github.com/hyperledger/fabric/core/handlers/validation/api/identities"
	. "github.com/hyperledger/fabric/core/handlers/validation/api/state"
	defaultvscc "github.com/hyperledger/fabric/core/handlers/validation/builtin"
	"github.com/hyperledger/fabric/protoutil"
)

type (
	NGACValidationFactory struct{}

	validator struct {
		idDeserializer   IdentityDeserializer
		ngacValidator    *NGACVSCC
		defaultValidator validation.Plugin
	}
)

func NewPluginFactory() validation.PluginFactory {
	return &NGACValidationFactory{}
}

func (*NGACValidationFactory) New() validation.Plugin {
	return validator{}
}

func (v validator) Validate(block *common.Block, namespace string, txPosition int, actionPosition int, contextData ...validation.ContextDatum) error {
	var err error

	// perform the default validation first
	if err = v.defaultValidator.Validate(block, namespace, txPosition, actionPosition, contextData); err != nil {
		return err
	}

	// deserialize the transaction from the block
	envelope := &common.Envelope{}
	if envelope, err = protoutil.GetEnvelopeFromBlock(block.Data.Data[txPosition]); err != nil {
		return err
	}

	payload := &common.Payload{}
	if err = proto.Unmarshal(envelope.Payload, payload); err != nil {
		return err
	}

	tx := &peer.Transaction{}
	if tx, err = protoutil.UnmarshalTransaction(payload.Data); err != nil {
		return err
	}

	// process each action in the transaction, checking if the user can perform the invoked chaincode function
	for _, action := range tx.Actions {
		actPayload := &peer.ChaincodeActionPayload{}
		if actPayload, err = protoutil.UnmarshalChaincodeActionPayload(action.Payload); err != nil {
			return err
		}

		propResPayload := &peer.ProposalResponsePayload{}
		if propResPayload, err = protoutil.UnmarshalProposalResponsePayload(actPayload.Action.ProposalResponsePayload); err != nil {
			return err
		}

		ccAction := &peer.ChaincodeAction{}
		if err = proto.Unmarshal(propResPayload.Extension, ccAction); err != nil {
			return err
		}

		ccPropPayload := &peer.ChaincodeProposalPayload{}
		if ccPropPayload, err = protoutil.UnmarshalChaincodeProposalPayload(actPayload.ChaincodeProposalPayload); err != nil {
			return err
		}

		ccInvokeSpec := &peer.ChaincodeInvocationSpec{}
		if err = proto.Unmarshal(ccPropPayload.Input, ccInvokeSpec); err != nil {
			return err
		}

		// get user identity
		id := &IdentityIdentifier{}
		if id, err = v.GetIdentity(block); err != nil {
			return err
		}

		var ok bool
		if ok, err = v.ngacValidator.CanInvoke(id, namespace, ccInvokeSpec); err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("user cannot invoke")
		}
	}

	return nil
}

func (v validator) Init(dependencies ...validation.Dependency) error {
	var sf StateFetcher

	for _, d := range dependencies {
		// initialize the StateFetcher which will be used by the ngac validator
		if stateFetcher, ok := d.(StateFetcher); ok {
			sf = stateFetcher
		}

		// initialize the IdentityDeserializer which will be used by the plugin to deserialize user info from v block
		if deserializer, isIdentityDeserializer := d.(IdentityDeserializer); isIdentityDeserializer {
			v.idDeserializer = deserializer
		}
	}

	if sf == nil {
		return fmt.Errorf("error initializing VSCC: StateFetcher not passed")
	}

	if v.idDeserializer == nil {
		return fmt.Errorf("error initializing VSCC: IdentityDeserializer not passed")
	}

	v.ngacValidator = New(sf)

	factory := &defaultvscc.DefaultValidationFactory{}
	v.defaultValidator = factory.New()
	if err := v.defaultValidator.Init(dependencies...); err != nil {
		return fmt.Errorf("error creating default vscc: %w", err)
	}

	return nil
}

// MarshaledSignedData is used to store information about a block signature
// This struct follows the sample validation plugin provided in the fabric core
type MarshaledSignedData struct {
	Data      []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	Identity  []byte `protobuf:"bytes,2,opt,name=identity,proto3" json:"identity,omitempty"`
}

func (sd *MarshaledSignedData) Reset() {
	*sd = MarshaledSignedData{}
}

func (*MarshaledSignedData) String() string {
	return ""
}

func (*MarshaledSignedData) ProtoMessage() {

}

func (v validator) GetIdentity(block *common.Block) (*IdentityIdentifier, error) {
	txData := block.Data.Data[0]
	txn := &MarshaledSignedData{}
	err := proto.Unmarshal(txData, txn)

	// Check if the chaincode is instantiated
	state, err := v.ngacValidator.sf.FetchState()
	if err != nil {
		return nil, err
	}
	defer state.Done()

	// Check the identity can be properly deserialized
	identity, err := v.idDeserializer.DeserializeIdentity(txn.Identity)
	if err != nil {
		return nil, err
	}

	return identity.GetIdentityIdentifier(), nil
}
