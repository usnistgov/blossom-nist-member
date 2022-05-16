package api

import (
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"testing"
)

func TestSwID(t *testing.T) {
	ctx := newTestStub(t)
	err := ctx.SetClientIdentity(mocks.Super)
	require.NoError(t, err)

	bcc := BlossomSmartContract{}

	onboardTestAsset(t, ctx, "123", "myasset", []string{"1", "2"})
	require.NoError(t, err)

	requestTestAccount(t, ctx, A1MSP)

	err = ctx.SetClientIdentity(mocks.A1SystemAdmin)
	require.NoError(t, err)

	err = ctx.SetTransient("checkout", requestCheckoutTransientInput{"123", 1})
	require.NoError(t, err)
	err = bcc.RequestCheckout(ctx)
	require.NoError(t, err)

	err = ctx.SetClientIdentity(mocks.Super)
	require.NoError(t, err)

	err = ctx.SetTransient("checkout", approveCheckoutTransientInput{A1MSP, "123"})
	require.NoError(t, err)
	err = bcc.ApproveCheckout(ctx)
	require.NoError(t, err)

	err = ctx.SetClientIdentity(mocks.A1SystemAdmin)
	require.NoError(t, err)

	// report swid on license they did not checkout
	t.Run("report swid on license A1 did not checkout", func(t *testing.T) {
		err = ctx.SetTransient("swid", reportSwIDTransientInput{
			PrimaryTag: "primary_tag_1",
			Asset:      "123",
			License:    "2",
			Xml:        "swid_xml",
		})
		require.NoError(t, err)
		err = bcc.ReportSwID(ctx)
		require.Error(t, err)
	})

	err = ctx.SetTransient("swid", reportSwIDTransientInput{
		PrimaryTag: "primary_tag_1",
		Asset:      "123",
		License:    "1",
		Xml:        "swid_xml",
	})
	require.NoError(t, err)
	err = bcc.ReportSwID(ctx)
	require.NoError(t, err)

	// check swid in collection
	err = ctx.SetTransient("swid", swidTransientInput{
		Account:    A1MSP,
		PrimaryTag: "primary_tag_1",
	})
	require.NoError(t, err)

	swid := &model.SwID{}
	swid, err = bcc.GetSwID(ctx)
	require.NoError(t, err)
	require.Equal(t, &model.SwID{
		PrimaryTag: "primary_tag_1",
		XML:        "swid_xml",
		Asset:      "123",
		License:    "1",
	}, swid)

	swids := make([]*model.SwID, 0)
	swids, err = bcc.GetSwIDsAssociatedWithAsset(ctx, A1MSP, "123")
	require.NoError(t, err)
	require.Equal(t, 1, len(swids))
	require.Equal(t, &model.SwID{
		PrimaryTag: "primary_tag_1",
		XML:        "swid_xml",
		Asset:      "123",
		License:    "1",
	}, swids[0])

	err = ctx.SetClientIdentity(mocks.A1SystemOwner)
	require.NoError(t, err)

	// try deleting as unauthorized user
	t.Run("delete swid as unauthorized user", func(t *testing.T) {
		err = ctx.SetTransient("swid", swidTransientInput{
			Account:    A1MSP,
			PrimaryTag: "primary_tag_1",
		})
		require.NoError(t, err)
		err = bcc.DeleteSwID(ctx)
		require.Error(t, err)
	})

	err = ctx.SetClientIdentity(mocks.A1SystemAdmin)
	require.NoError(t, err)

	err = ctx.SetTransient("swid", swidTransientInput{
		Account:    A1MSP,
		PrimaryTag: "primary_tag_1",
	})
	require.NoError(t, err)
	err = bcc.DeleteSwID(ctx)
	require.NoError(t, err)

	err = ctx.SetTransient("swid", swidTransientInput{
		Account:    A1MSP,
		PrimaryTag: "primary_tag_1",
	})
	require.NoError(t, err)
	_, err = bcc.GetSwID(ctx)
	require.Error(t, err)
}
