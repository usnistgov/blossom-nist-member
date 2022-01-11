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

func adminUA(adminMSP string) string {
	return fmt.Sprintf("%s_UA", adminMSP)
}

func LoadCatalogPolicy(adminUser string, adminMSP string) (policy.Store, error) {
	policyStore := memory.NewPolicyStore()
	err := author.Author(policyStore,
		// catalog policy
		create.PolicyClass("catalog"),

		// ua
		create.UserAttribute(adminUA(adminMSP)).In("catalog"),
		create.User(adminUser).In(adminUA(adminMSP)),

		// oa
		create.ObjectAttribute(BlossomOA).In("catalog"),
		create.Object(BlossomObject).In(BlossomOA),

		// grants
		grant.UserAttribute(adminUA(adminMSP)).Permissions(policy.AllOps).On(BlossomOA),

		// onboarding/offboarding policy
		create.UserAttribute("Onboarders").In("catalog"),

		create.ObjectAttribute("assets").In("catalog"),

		grant.UserAttribute(adminUA(adminMSP)).Permissions(policy.AllOps).On("assets"),
		grant.UserAttribute("Onboarders").Permissions("onboard_asset", "offboard_asset").On("assets"),

		create.Obligation("onboard_asset").
			When(policy.AnyUserSubject).
			Performs("onboard_asset", "asset_id").
			Do(create.Object("<asset_id>").In("assets")),
		create.Obligation("offboard_asset").
			When(policy.AnyUserSubject).
			Performs("offboard_asset", "asset_id").
			Do(remove.Object("<asset_id>")),
	)
	if err != nil {
		return nil, err
	}

	return policyStore, nil
}

func LoadAccountPolicy(adminUser string, adminMSP string) (policy.Store, error) {
	policyStore := memory.NewPolicyStore()
	err := author.Author(policyStore,
		create.PolicyClass("super"),
		create.UserAttribute(adminUA(adminMSP)).In("super"),
		create.User(adminUser).In(adminUA(adminMSP)),

		create.ObjectAttribute(BlossomOA).In("super"),
		create.Object(BlossomObject).In(BlossomOA),

		grant.UserAttribute(adminUA(adminMSP)).Permissions(policy.AllOps).On(BlossomOA),

		create.Obligation("approve_account").
			When(policy.AnyUserSubject).
			Performs("approve_account", "accountName", "sysOwner", "sysAdmin", "acqSpec").
			Do(
				create.PolicyClass("RBAC"),
				// ua
				create.UserAttribute("RBAC_UA").In("RBAC"),
				create.UserAttribute("Account_UA").In("RBAC_UA"),
				create.UserAttribute("SystemOwner").In("RBAC_UA"),
				create.UserAttribute("SystemAdministrator").In("RBAC_UA"),
				create.UserAttribute("AcquisitionSpecialist").In("RBAC_UA"),
				create.UserAttribute("Approvers").In("RBAC_UA"),
				// assign.UserAttribute(BlossomAdminUser).To("Approvers"),

				create.User("<sysOwner>").In("SystemOwner", "Account_UA"),
				create.User("<sysAdmin>").In("SystemAdministrator", "Account_UA"),
				create.User("<acqSpec>").In("AcquisitionSpecialist", "Account_UA"),

				// oa
				create.ObjectAttribute("RBAC_OA").In("RBAC"),
				create.ObjectAttribute("Account_OA").In("RBAC_OA"),
				create.Object(AccountObjectName("<accountName>")).
					WithProperties("account", "<accountName>", "type", "account").
					In("Account_OA"),

				// grants
				grant.UserAttribute(adminUA(adminMSP)).Permissions(policy.AllOps).On("RBAC_UA"),
				grant.UserAttribute(adminUA(adminMSP)).Permissions(policy.AllOps).On("RBAC_OA"),

				grant.UserAttribute("SystemOwner").Permissions("upload_ato").On("Account_OA"),
				grant.UserAttribute("SystemAdministrator").
					Permissions("check_out", "initiate_check_in", "report_swid", "delete_swid").
					On("Account_OA"),
				grant.UserAttribute("Approvers").
					Permissions("approve_checkout", "process_check_in").
					On("Account_OA"),

				// status policy
				create.PolicyClass("Status"),
				// ua
				create.UserAttribute("Status_UA").In("Status"),
				create.UserAttribute("active").In("Status_UA"),
				create.UserAttribute("pending").In("Status_UA"),
				create.UserAttribute("inactive").In("pending"),

				assign.UserAttribute("Account_UA").To("pending"),

				// oa
				create.ObjectAttribute("Status_OA").In("Status"),
				create.ObjectAttribute("status_accounts_OA").In("Status_OA"),
				create.ObjectAttribute("status_assets_OA").In("Status_OA"),
				create.ObjectAttribute("status_swids_OA").In("Status_OA"),

				assign.Object(AccountObjectName("<accountName>")).To("status_accounts_OA"),

				// grants
				grant.UserAttribute(adminUA(adminMSP)).Permissions(policy.AllOps).On("Status_UA"),
				grant.UserAttribute(adminUA(adminMSP)).Permissions(policy.AllOps).On("Status_OA"),
				grant.UserAttribute("active").Permissions(policy.AllOps).On("status_accounts_OA"),
				grant.UserAttribute("pending").Permissions("upload_ato").On("status_accounts_OA"),
				grant.UserAttribute("active").Permissions(policy.AllOps).On("status_assets_OA"),
				grant.UserAttribute("active").Permissions(policy.AllOps).On("status_swids_OA"),
			),

		create.Obligation("set_account_active").
			When(policy.AnyUserSubject).
			Performs("set_account_active").
			Do(
				deassign.UserAttribute("Account_UA").From("pending"),
				deassign.UserAttribute("Account_UA").From("inactive"),
				assign.UserAttribute("Account_UA").To("active"),
			),
		create.Obligation("set_account_pending").
			When(policy.AnyUserSubject).
			Performs("set_account_pending").
			Do(
				assign.UserAttribute("Account_UA").To("pending"),
				deassign.UserAttribute("Account_UA").From("inactive"),
				deassign.UserAttribute("Account_UA").From("active"),
			),
		create.Obligation("set_account_inactive").
			When(policy.AnyUserSubject).
			Performs("set_account_inactive").
			Do(
				deassign.UserAttribute("Account_UA").From("pending"),
				assign.UserAttribute("Account_UA").To("inactive"),
				deassign.UserAttribute("Account_UA").From("active"),
			),
	)

	if err != nil {
		return nil, err
	}

	return policyStore, nil
}
