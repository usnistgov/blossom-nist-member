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

	err := author.Author(policyStore,
		// RBAC policy
		create.PolicyClass("RBAC_PC"),

		// admin policy
		create.UserAttribute("RBAC_UA").In("RBAC_PC"),
		create.ObjectAttribute("RBAC_OA").In("RBAC_PC"),
		create.UserAttribute(adminUA).In("RBAC_PC"),
		create.User(adminUser).In(adminUA),

		create.ObjectAttribute(BlossomOA).In("RBAC_OA"),
		create.Object(BlossomObject).In(BlossomOA),

		// grant admin permission on user and object attributes in this policy class
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On("RBAC_UA"),
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On("RBAC_OA"),

		// create roles in RBAC ua
		create.UserAttribute("SystemOwner").In("RBAC_UA"),
		create.UserAttribute("SystemAdministrator").In("RBAC_UA"),
		create.UserAttribute("AcquisitionSpecialist").In("RBAC_UA"),

		create.ObjectAttribute("accounts_OA.RBAC_PC").In("RBAC_OA"),
		create.UserAttribute("accounts_UA.RBAC_PC").In("RBAC_UA"),

		grant.UserAttribute("SystemOwner").
			Permissions("upload_ato").
			On("accounts_OA.RBAC_PC"),
		grant.UserAttribute("SystemAdministrator").
			Permissions("check_out", "initiate_check_in", "report_swid", "delete_swid").
			On("accounts_OA.RBAC_PC"),

		// assets policy
		create.PolicyClass("Assets_PC"),

		// admin policy
		create.UserAttribute("Assets_UA").In("Assets_PC"),
		create.ObjectAttribute("Assets_OA").In("Assets_PC"),
		assign.UserAttribute(adminUA).To("Assets_PC"),

		assign.ObjectAttribute(BlossomOA).To("Assets_OA"),

		// grant admin permission on user and object attributes in this policy class
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On("Assets_UA"),
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On("Assets_OA"),

		// onboarding/offboarding policy
		create.UserAttribute("asset_managers").In("Assets_UA"),

		create.ObjectAttribute("assets").In("Assets_OA"),
		create.Object("all_assets").In("assets"),

		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On("assets"),
		grant.UserAttribute("asset_managers").Permissions("onboard_asset", "offboard_asset", "view_assets", "view_asset_private", "view_asset_public").On("assets"),

		create.Obligation("onboard_asset").
			When(policy.AnyUserSubject).
			Performs("onboard_asset", "asset_id").
			Do(create.Object("<asset_id>").In("assets")),
		create.Obligation("offboard_asset").
			When(policy.AnyUserSubject).
			Performs("offboard_asset", "asset_id").
			Do(remove.Object("<asset_id>")),

		// view catalog policy
		create.UserAttribute("accounts_UA.Assets_PC").In("Assets_UA"),
		grant.UserAttribute("accounts_UA.Assets_PC").Permissions("view_assets", "view_asset_public").On("assets"),

		// status policy
		create.PolicyClass("Status_PC"),

		// admin policy
		create.UserAttribute("Status_UA").In("Status_PC"),
		create.ObjectAttribute("Status_OA").In("Status_PC"),
		assign.UserAttribute(adminUA).To("Status_PC"),

		assign.ObjectAttribute(BlossomOA).To("Status_OA"),

		// grant admin permission on user and object attributes in this policy class
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On("Status_UA"),
		grant.UserAttribute(adminUA).Permissions(policy.AllOps).On("Status_OA"),

		// ua
		create.UserAttribute("active").In("Status_UA"),
		create.UserAttribute("pending").In("Status_UA"),
		create.UserAttribute("inactive").In("pending"),

		// oa
		create.ObjectAttribute("accounts_OA.Status_PC").In("Status_OA"),
		create.ObjectAttribute("catalog_OA.Status_PC").In("Status_OA"),
		assign.Object(BlossomObject).To("catalog_OA.Status_PC"),
		assign.Object("all_assets").To("catalog_OA.Status_PC"),

		// grants
		grant.UserAttribute("active").Permissions(policy.AllOps).On("accounts_OA.Status_PC"),
		grant.UserAttribute("pending").Permissions("upload_ato").On("accounts_OA.Status_PC"),
		grant.UserAttribute("active").Permissions("view_assets", "view_asset_public").On("catalog_OA.Status_PC"),

		create.Obligation("set_account_active").
			When(policy.AnyUserSubject).
			Performs("set_account_active", "accountName").
			Do(
				deassign.UserAttribute(AccountUA("<accountName>")).From("pending"),
				deassign.UserAttribute(AccountUA("<accountName>")).From("inactive"),
				assign.UserAttribute(AccountUA("<accountName>")).To("active"),
			),
		create.Obligation("set_account_pending").
			When(policy.AnyUserSubject).
			Performs("set_account_pending", "accountName").
			Do(
				assign.UserAttribute(AccountUA("<accountName>")).To("pending"),
				deassign.UserAttribute(AccountUA("<accountName>")).From("inactive"),
				deassign.UserAttribute(AccountUA("<accountName>")).From("active"),
			),
		create.Obligation("set_account_inactive").
			When(policy.AnyUserSubject).
			Performs("set_account_inactive", "accountName").
			Do(
				deassign.UserAttribute(AccountUA("<accountName>")).From("pending"),
				assign.UserAttribute(AccountUA("<accountName>")).To("inactive"),
				deassign.UserAttribute(AccountUA("<accountName>")).From("active"),
			),

		// general obligations
		create.Obligation("approve_account").
			When(policy.AnyUserSubject).
			Performs("approve_account", "accountName").
			Do(
				// create account obj and assign to rbac and status policies
				create.Object(AccountObjectName("<accountName>")).
					WithProperties("account", "<accountName>", "type", "account").
					In("accounts_OA.RBAC_PC", "accounts_OA.Status_PC"),

				// create a UA for the account
				create.UserAttribute(AccountUA("<accountName>")).In("accounts_UA.Assets_PC", "pending"),
			),
	)

	if err != nil {
		return nil, fmt.Errorf("error building policy: %v", err)
	}

	return policyStore, nil
}
