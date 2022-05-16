package api

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pdp"
)

type (
	// BlossomSmartContract is the struct that will implement the below interfaces and expose the functions as
	// smart contract functions
	BlossomSmartContract struct {
		contractapi.Contract
	}

	// AccountInterface provides the functions to interact with Accounts in blossom.
	AccountInterface interface {
		// RequestAccount allows accounts to request an account in the Blossom system. The systemOwner, systemAdmin, and
		// acqSpec will be added as users to the NGAC graph, and given the appropriate permissions on the account. The
		// ato can be empty and uploaded via UploadATO later. The name of the acount is the MSPID of the requesting
		// user's member.
		// TRANSIENT MAP: export ACCOUNT=$(echo -n "{\"system_owner\":\"\",\"system_admin\":\"\",\"acquisition_specialist\":\"\",\"ato\":\"\"}" | base64 | tr -d \\n)
		RequestAccount(ctx contractapi.TransactionContextInterface) error

		// ApproveAccount initializes the account's NGAC graph in the account's PDC, with the user invoking this function
		// being the admin in the graph.  The status of the account will be Pending after execution.  The admin user can
		// call UpdateAccountStatus to update the status of the account.
		ApproveAccount(ctx contractapi.TransactionContextInterface, account string) error

		// UploadATO updates the ATO field of the Account with the given name. This is just the ATO attestation not a full ATO.
		// TRANSIENT MAP: export ATO=$(echo -n "{\"ato\":\"\"}" | base64 | tr -d \\n)
		UploadATO(ctx contractapi.TransactionContextInterface) error

		// UpdateAccountStatus updates the status of an account in Blossom. The status is one of:
		//		"PENDING_APPROVAL",
		//		"PENDING_ATO",
		//		"AUTHORIZED",
		//		"UNAUTHORIZED_DENIED",
		//		"UNAUTHORIZED_ATO",
		//		"UNAUTHORIZED_OPTOUT",
		//		"UNAUTHORIZED_SECURITY_RISK",
		//		"UNAUTHORIZED_ROB"
		// Updating the status to Authorized allows the account to read and write to blossom.
		// Updating the status to Pending allows the account to read write only account related information such as ATOs.
		// Updating the status to Inactive provides the same NGAC consequences as Pending
		UpdateAccountStatus(ctx contractapi.TransactionContextInterface, account string, status string) error

		// GetAccounts returns the public info of all accounts that are registered with Blossom.
		GetAccounts(ctx contractapi.TransactionContextInterface) ([]*model.AccountPublic, error)

		// GetAccount returns the account information of the account with the provided name.  Any fields of any account the user
		// does not have access to will not be returned.
		GetAccount(ctx contractapi.TransactionContextInterface, account string) (*model.Account, error)

		// GetHistory returns the transaction history of the account.
		GetHistory(ctx contractapi.TransactionContextInterface, account string) ([]model.HistorySnapshot, error)
	}

	// AssetsInterface provides the functions to interact with Assets in fabric.
	AssetsInterface interface {
		// OnboardAsset adds a new software asset to Blossom.  This will create a new asset object on the ledger and in the
		// NGAC graph. Assets are identified by the ID field. The user performing the request will need to
		// have permission to add an asset to the ledger. The asset will be an object attribute in NGAC and the
		// asset licenses will be objects that are assigned to the asset.
		// TRANSIENT MAP: export ATO=$(echo -n "{\"licenses\":\"\"}" | base64 | tr -d \\n)
		OnboardAsset(ctx contractapi.TransactionContextInterface, id string, name string, onboardDate string, expiration string) error

		// OffboardAsset removes an existing asset in Blossom.  This will remove the license from the ledger
		// and from NGAC. An error will be returned if there are any accounts that have checked out the asset
		// and the licenses are not returned
		OffboardAsset(ctx contractapi.TransactionContextInterface, id string) error

		// GetAssets returns all software assets in Blossom. This information includes which accounts have licenses for each
		// asset.
		GetAssets(ctx contractapi.TransactionContextInterface) ([]*model.AssetPublic, error)

		// GetAsset returns the info for the asset with the given asset ID.
		GetAsset(ctx contractapi.TransactionContextInterface, id string) (*model.Asset, error)

		// RequestCheckout requests software licenses for an account.  The requesting user must have permission to request
		// (i.e. System Administrator). The amount parameter is the amount of software licenses the account is requesting.
		// This number is subtracted from the total available for the asset. Returns the set of licenses that are now assigned to
		// the account.
		// TRANSIENT MAP: export CHECKOUT=$(echo -n "{\"asset_id\":\"\", \"amount\":}" | base64 | tr -d \\n)
		RequestCheckout(ctx contractapi.TransactionContextInterface) error

		// GetCheckoutRequests returns an array of checkout requests made by the account.
		GetCheckoutRequests(ctx contractapi.TransactionContextInterface, account string) ([]CheckoutRequest, error)

		// ApproveCheckout approves a checkout request made by an account.  The requested licenses for the asset will be
		// added to the account's private data collection. A user on the account can then call Licenses to get the approved
		// license keys.
		// TRANSIENT MAP: export CHECKOUT=$(echo -n "{\"account\":\"\", \"asset_id\":\"\"}" | base64 | tr -d \\n)
		ApproveCheckout(ctx contractapi.TransactionContextInterface) error

		// GetLicenses get the license keys for an asset that an account has access to in their private data collection.
		// The account is extracted from the requesting identity.
		GetLicenses(ctx contractapi.TransactionContextInterface, account, assetID string) (map[string]string, error)

		// InitiateCheckin starts the process of returning licenses to Blossom. This is serves as a request to the blossom
		// admin to process the return of the licenses. This is because only the blossom admin can write to the licenses
		// private data collection to return the licenses to the available pool.
		// TRANSIENT MAP: export CHECKIN=$(echo -n "{\"asset_id\":\"\", \"licenses\":[]}" | base64 | tr -d \\n)
		InitiateCheckin(ctx contractapi.TransactionContextInterface) error

		// GetInitiatedCheckins returns the list of checkins initiated by the given account.
		GetInitiatedCheckins(ctx contractapi.TransactionContextInterface, account string) ([]CheckinRequest, error)

		// ProcessCheckin processes an account's checkin request (from InitiateCheckin) and returns the licenses to the
		// available pool in the licenses private data collection.
		// TRANSIENT MAP: export CHECKIN=$(echo -n "{\"asset_id\":\"\", \"account\":\"\"}" | base64 | tr -d \\n)
		ProcessCheckin(ctx contractapi.TransactionContextInterface) error
	}

	// SwIDInterface provides the functions to interact with SwID tags in fabric.
	SwIDInterface interface {
		// ReportSwID is used by Accounts to report to Blossom when a software user has installed a piece of software associated
		// with an asset that account has checked out. The account is extracted from the requesting identity.  The account
		// must have checked out the defined license or this function will fail.
		// TRANSIENT MAP: export ATO=$(echo -n "{\"primary_tag\":\"123\",\"asset\":\"101\",\"license\":\"asset1-license-1\",\"xml\":\"<swid></swid>\"}" | base64 | tr -d \\n)
		ReportSwID(ctx contractapi.TransactionContextInterface) error

		// DeleteSwID deletes a swid from the ledger. This would happen in the case of an organziation returning licenses,
		// and the swid no longer being valid.  The requesting user will need to have the correct permissions in NGAC
		// to do so.  The user with pemrission is the system_owner as defined in the account info.
		// TRANSIENT MAP: export swid=$(echo -n "{\"primary_tag\":\"\",\"account\":\"\"}" | base64 | tr -d \\n)
		DeleteSwID(ctx contractapi.TransactionContextInterface) error

		// GetSwID returns the SwID object including the XML that matches the provided primaryTag parameter.
		// TRANSIENT MAP: export swid=$(echo -n "{\"primary_tag\":\"\",\"account\":\"\"}" | base64 | tr -d \\n)
		GetSwID(ctx contractapi.TransactionContextInterface) (*model.SwID, error)

		// GetSwIDsAssociatedWithAsset returns the SwIDs that are associated with the given asset for an account.
		GetSwIDsAssociatedWithAsset(ctx contractapi.TransactionContextInterface, account string, assetID string) ([]*model.SwID, error)
	}
)

func (b *BlossomSmartContract) InitNGAC(ctx contractapi.TransactionContextInterface) error {
	return pdp.InitCatalogNGAC(ctx)
}
