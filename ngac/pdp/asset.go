package pdp

import (
	"time"

	"github.com/PM-Master/policy-machine-go/pdp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/operations"
	"github.com/usnistgov/blossom/chaincode/ngac/pap"
	assetpap "github.com/usnistgov/blossom/chaincode/ngac/pap/asset"
	rbacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/policy/rbac"
)

type AssetDecider struct {
	// user is the user that is currently executing a function
	user string
	// pap is the policy administration point for licenses
	pap *pap.AssetAdmin
	// decider is the NGAC decider used to make decisions
	decider pdp.Decider
}

// NewAssetDecider creates a new AssetDecider with the user from the ctx and a NGAC Decider using the NGAC graph
// from the ledger.
func NewAssetDecider() *AssetDecider {
	return &AssetDecider{}
}

func (l *AssetDecider) setup(ctx contractapi.TransactionContextInterface) error {
	user, err := GetUser(ctx)
	if err != nil {
		return errors.Wrapf(err, "error getting user from request")
	}

	l.user = user

	// initialize the license policy administration point
	l.pap, err = pap.NewAssetAdmin(ctx)
	if err != nil {
		return errors.Wrapf(err, "error initializing agency administraion point")
	}

	l.decider = pdp.NewDecider(l.pap.Graph())

	return nil
}

func (l *AssetDecider) FilterAsset(ctx contractapi.TransactionContextInterface, asset *model.Asset) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up asset decider")
	}

	return l.filterAsset(asset)
}

func (l *AssetDecider) filterAsset(asset *model.Asset) error {
	permissions, err := l.decider.ListPermissions(l.user, assetpap.ObjectAttribute(asset.ID))
	if err != nil {
		return errors.Wrapf(err, "error getting permissions for user %s on asset %s", l.user, asset.Name)
	}

	if !permissions.Contains(operations.ViewAsset) {
		asset.ID = ""
		asset.Name = ""
		asset.TotalAmount = 0
		asset.Available = 0
		asset.Cost = 0
		asset.OnboardingDate = time.Time{}
		asset.Expiration = time.Time{}
		asset.Licenses = make([]string, 0)
		asset.AvailableLicenses = make([]string, 0)
		asset.CheckedOut = make(map[string]map[string]time.Time)
		return nil
	}

	if !permissions.Contains(operations.ViewAllLicenses) {
		asset.Licenses = make([]string, 0)
	}

	if !permissions.Contains(operations.ViewAvailableLicenses) {
		asset.AvailableLicenses = make([]string, 0)
	}

	if !permissions.Contains(operations.ViewCheckedOut) {
		asset.CheckedOut = make(map[string]map[string]time.Time)
	}

	return nil
}

func (l *AssetDecider) FilterAssets(ctx contractapi.TransactionContextInterface, assets []*model.Asset) ([]*model.Asset, error) {
	if err := l.setup(ctx); err != nil {
		return nil, errors.Wrapf(err, "error setting up asset decider")
	}

	filteredAssets := make([]*model.Asset, 0)
	for _, asset := range assets {
		if err := l.filterAsset(asset); err != nil {
			return nil, errors.Wrapf(err, "error filtering asset")
		}

		if asset.ID == "" {
			continue
		}

		filteredAssets = append(filteredAssets, asset)
	}

	return filteredAssets, nil
}

func (l *AssetDecider) OnboardAsset(ctx contractapi.TransactionContextInterface, asset *model.Asset) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up asset decider")
	}

	// check user can onboard license
	if ok, err := l.decider.HasPermissions(l.user, rbacpolicy.AssetsOA, operations.OnboardLicense); err != nil {
		return errors.Wrapf(err, "error checking if user %s can onboard a license", l.user)
	} else if !ok {
		return ErrAccessDenied
	}

	return l.pap.OnboardAsset(ctx, asset)
}

func (l *AssetDecider) OffboardAsset(ctx contractapi.TransactionContextInterface, licenseID string) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up agency decider")
	}

	// check user can onboard license
	if ok, err := l.decider.HasPermissions(l.user, licenseID, operations.OffboardLicense); err != nil {
		return errors.Wrapf(err, "error checking if user %s can offboard a license", l.user)
	} else if !ok {
		return ErrAccessDenied
	}

	return l.pap.OffboardAsset(ctx, licenseID)
}

func (l *AssetDecider) Checkout(ctx contractapi.TransactionContextInterface, agencyName string, assetID string,
	licenses map[string]time.Time) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up asset decider")
	}

	// check user can checkout license
	if ok, err := l.decider.HasPermissions(l.user, assetID, operations.CheckOut); err != nil {
		return errors.Wrapf(err, "error checking if user %s can checkout an asset", l.user)
	} else if !ok {
		return ErrAccessDenied
	}

	return l.pap.Checkout(ctx, agencyName, assetID, licenses)
}

func (l *AssetDecider) Checkin(ctx contractapi.TransactionContextInterface, agencyName string, licenseID string,
	keys []string) error {
	if err := l.setup(ctx); err != nil {
		return errors.Wrapf(err, "error setting up asset decider")
	}

	// check user can checkin license
	if ok, err := l.decider.HasPermissions(l.user, licenseID, operations.CheckIn); err != nil {
		return errors.Wrapf(err, "error checking if user %s can checkin an asset", l.user)
	} else if !ok {
		return ErrAccessDenied
	}

	return l.pap.Checkin(ctx, agencyName, licenseID, keys)
}
