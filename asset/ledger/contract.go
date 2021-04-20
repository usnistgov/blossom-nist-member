package ledger

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

// Contract for managing the ledger
type Contract struct {
	contractapi.Contract
}

// InitLedger initializes the ledger components including: Agencies, Licenses, and SwID tags. This method also
// invokes NGAC chaincode to initialize the NGAC components.
func (c Contract) InitLedger(ctx contractapi.TransactionContextInterface) {

}
