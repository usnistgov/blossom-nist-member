package pap

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/PM-Master/policy-machine-go/policy"
	"github.com/PM-Master/policy-machine-go/policy/author"
	"github.com/PM-Master/policy-machine-go/policy/author/assign"
	"github.com/PM-Master/policy-machine-go/policy/author/create"
	"github.com/PM-Master/policy-machine-go/policy/author/deassign"
	"github.com/PM-Master/policy-machine-go/policy/author/grant"
	"github.com/PM-Master/policy-machine-go/policy/author/remove"
)

const (
	BlossomObject = "blossom_object"
	BlossomOA     = "Blossom_OA"
)

func AccountObjectName(accountName string) string {
	return fmt.Sprintf("%s_object", accountName)
}

func AccountUA(accountName string) string {
	return fmt.Sprintf("%s_UA", accountName)
}

func adminUA(adminMSP string) string {
	return fmt.Sprintf("%s_UA", adminMSP)
}

func LoadCatalogPolicy(adminUser string, adminMSP string) (policy.Store, error) {
	policyStore := memory.NewPolicyStore()

	adminUA := adminUA(adminMSP)

	const (
		RbacPolicyClass          = "RBAC_PC"
		RbacUserAttr             = "RBAC_UA"
		RbacObjectAttr           = "RBAC_OA"
		SystemOwnerAttr          = "SystemOwner"
		SystemAdminAttr          = "SystemAdministrator"
		AcqSpecAttr              = "AcquisitionSpecialist"
		AccountsObjectAttrInRBAC = "accounts_OA.RBAC_PC"
		AccountsUserAttrInRBAC   = "accounts_UA.RBAC_PC"

		AssetsPolicyClass          = "Assets_PC"
		AssetsBaseUserAttr         = "Assets_UA"
		AssetsBaseObjectAttr       = "Assets_OA"
		AssetManagersAttr          = "asset_managers"
		AssetsObjectAttr           = "assets"
		AccountsUserAttrInAssetsPC = "accounts_UA.Assets_PC"
		AllAssetsObj               = "all_assets"

		StatusPolicyClass            = "Status_PC"
		StatusBaseUserAttr           = "Status_UA"
		StatusBaseObjectAttr         = "Status_OA"
		ActiveAttr                   = "active"
		PendingAttr                  = "pending"
		InactiveAttr                 = "inactive"
		AccountsObjectAttrInStatusPC = "accounts_OA.Status_PC"
		CatalogObjectAttrInStatusPC  = "catalog_OA.Status_PC"
	)

	err := author.Author(policyStore,
		// RBAC policy
		create.PolicyClass(RbacPolicyClass),

		// admin policy
		create.UserAttribute(RbacUserAttr).In(RbacPolicyClass),
		create.ObjectAttribute(RbacObjectAttr).In(RbacPolicyClass),
		create.UserAttribute(adminUA).In(RbacPolicyClass),
		create.User(adminUser).In(adminUA),

		create.ObjectAttribute(BlossomOA).In(RbacObjectAttr),
		create.Object(BlossomObject).In(BlossomOA),

		// grant admin permission on user and object attributes in this policy class
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On(RbacUserAttr),
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On(RbacObjectAttr),

		// create roles in RBAC ua
		create.UserAttribute(SystemOwnerAttr).In(RbacUserAttr),
		create.UserAttribute(SystemAdminAttr).In(RbacUserAttr),
		create.UserAttribute(AcqSpecAttr).In(RbacUserAttr),

		create.ObjectAttribute(AccountsObjectAttrInRBAC).In(RbacObjectAttr),
		create.UserAttribute(AccountsUserAttrInRBAC).In(RbacUserAttr),

		grant.UserAttribute(SystemOwnerAttr).
			Permissions("upload_ato").
			On(AccountsObjectAttrInRBAC),
		grant.UserAttribute(SystemAdminAttr).
			Permissions("check_out", "initiate_check_in", "report_swid", "delete_swid").
			On(AccountsObjectAttrInRBAC),

		// assets policy
		create.PolicyClass(AssetsPolicyClass),

		// admin policy
		create.UserAttribute(AssetsBaseUserAttr).In(AssetsPolicyClass),
		create.ObjectAttribute(AssetsBaseObjectAttr).In(AssetsPolicyClass),
		assign.UserAttribute(adminUA).To(AssetsPolicyClass),

		assign.ObjectAttribute(BlossomOA).To(AssetsBaseObjectAttr),

		// grant admin permission on user and object attributes in this policy class
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On(AssetsBaseUserAttr),
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On(AssetsBaseObjectAttr),

		// onboarding/offboarding policy
		create.UserAttribute(AssetManagersAttr).In(AssetsBaseUserAttr),

		create.ObjectAttribute(AssetsObjectAttr).In(AssetsBaseObjectAttr),
		create.Object(AllAssetsObj).In(AssetsObjectAttr),

		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On(AssetsObjectAttr),
		grant.UserAttribute(AssetManagersAttr).Permissions("onboard_asset", "offboard_asset", "view_assets", "view_asset_private", "view_asset_public").On(AssetsObjectAttr),

		create.Obligation("onboard_asset").
			When(policy.AnyUserSubject).
			Performs("onboard_asset", "asset_id").
			Do(create.Object("<asset_id>").In(AssetsObjectAttr)),
		create.Obligation("offboard_asset").
			When(policy.AnyUserSubject).
			Performs("offboard_asset", "asset_id").
			Do(remove.Object("<asset_id>")),

		// view catalog policy
		create.UserAttribute(AccountsUserAttrInAssetsPC).In(AssetsBaseUserAttr),
		grant.UserAttribute(AccountsUserAttrInAssetsPC).Permissions("view_assets", "view_asset_public").On(AssetsObjectAttr),

		// status policy
		create.PolicyClass(StatusPolicyClass),

		// admin policy
		create.UserAttribute(StatusBaseUserAttr).In(StatusPolicyClass),
		create.ObjectAttribute(StatusBaseObjectAttr).In(StatusPolicyClass),
		assign.UserAttribute(adminUA).To(StatusPolicyClass),

		assign.ObjectAttribute(BlossomOA).To(StatusBaseObjectAttr),

		// grant admin permission on user and object attributes in this policy class
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On(StatusBaseUserAttr),
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On(StatusBaseObjectAttr),

		// ua
		create.UserAttribute(ActiveAttr).In(StatusBaseUserAttr),
		create.UserAttribute(PendingAttr).In(StatusBaseUserAttr),
		create.UserAttribute(InactiveAttr).In(PendingAttr),

		// oa
		create.ObjectAttribute(AccountsObjectAttrInStatusPC).In(StatusBaseObjectAttr),
		create.ObjectAttribute(CatalogObjectAttrInStatusPC).In(StatusBaseObjectAttr),
		assign.Object(BlossomObject).To(CatalogObjectAttrInStatusPC),
		assign.Object(AllAssetsObj).To(CatalogObjectAttrInStatusPC),

		// grants
		grant.UserAttribute(ActiveAttr).Permissions(policy.AllOps).On(AccountsObjectAttrInStatusPC),
		grant.UserAttribute(PendingAttr).Permissions("upload_ato").On(AccountsObjectAttrInStatusPC),
		grant.UserAttribute(ActiveAttr).Permissions("view_assets", "view_asset_public").On(CatalogObjectAttrInStatusPC),

		create.Obligation("set_account_active").
			When(policy.AnyUserSubject).
			Performs("set_account_active", "accountName").
			Do(
				deassign.UserAttribute(AccountUA("<accountName>")).From(PendingAttr),
				deassign.UserAttribute(AccountUA("<accountName>")).From(InactiveAttr),
				assign.UserAttribute(AccountUA("<accountName>")).To(ActiveAttr),
			),
		create.Obligation("set_account_pending").
			When(policy.AnyUserSubject).
			Performs("set_account_pending", "accountName").
			Do(
				assign.UserAttribute(AccountUA("<accountName>")).To(PendingAttr),
				deassign.UserAttribute(AccountUA("<accountName>")).From(InactiveAttr),
				deassign.UserAttribute(AccountUA("<accountName>")).From(ActiveAttr),
			),
		create.Obligation("set_account_inactive").
			When(policy.AnyUserSubject).
			Performs("set_account_inactive", "accountName").
			Do(
				deassign.UserAttribute(AccountUA("<accountName>")).From(PendingAttr),
				assign.UserAttribute(AccountUA("<accountName>")).To(InactiveAttr),
				deassign.UserAttribute(AccountUA("<accountName>")).From(ActiveAttr),
			),

		// general obligations
		create.Obligation("approve_account").
			When(policy.AnyUserSubject).
			Performs("approve_account", "accountName").
			Do(
				// create account obj and assign to rbac and status policies
				create.Object(AccountObjectName("<accountName>")).
					WithProperties("account", "<accountName>", "type", "account").
					In(AccountsObjectAttrInRBAC, AccountsObjectAttrInStatusPC),

				// create a UA for the account
				create.UserAttribute(AccountUA("<accountName>")).In(AccountsUserAttrInAssetsPC, PendingAttr),
			),
	)

	if err != nil {
		return nil, fmt.Errorf("error building policy: %v", err)
	}

	return policyStore, nil
}
