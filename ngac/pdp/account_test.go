package pdp

import (
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy"
	"testing"
	"time"
)

func TestUploadATO(t *testing.T) {
	decider := NewAccountDecider()

	t.Run("test a1 system owner", func(t *testing.T) {
		mock := mocks.New()

		err := initAccountTestGraph(t, mock)
		require.NoError(t, err)

		mock.SetUser(mocks.A1SystemOwner)

		err = decider.UploadATO(mock.Stub, "A1")
		require.NoError(t, err)
	})

	t.Run("test a1 system admin", func(t *testing.T) {
		mock := mocks.New()

		err := initAccountTestGraph(t, mock)
		require.NoError(t, err)

		mock.SetUser(mocks.A1SystemAdmin)

		err = decider.UploadATO(mock.Stub, "A1")
		require.Error(t, err)
	})
}

func TestUpdateAccountStatus(t *testing.T) {
	decider := NewAccountDecider()

	t.Run("test a1 system owner", func(t *testing.T) {
		mock := mocks.New()

		err := initAccountTestGraph(t, mock)
		require.NoError(t, err)

		mock.SetUser(mocks.A1SystemOwner)

		err = decider.UpdateAccountStatus(mock.Stub, "A1", "test")
		require.Error(t, err)
	})

	t.Run("test a1 system admin", func(t *testing.T) {
		mock := mocks.New()

		err := initAccountTestGraph(t, mock)
		require.NoError(t, err)

		mock.SetUser(mocks.A1SystemAdmin)

		err = decider.UpdateAccountStatus(mock.Stub, "A1", "test")
		require.Error(t, err)
	})

	t.Run("test super", func(t *testing.T) {
		mock := mocks.New()

		err := initAccountTestGraph(t, mock)
		require.NoError(t, err)

		mock.SetUser(mocks.Super)

		err = decider.UpdateAccountStatus(mock.Stub, "A1", "test")
		require.NoError(t, err)
	})
}

func TestFilterAccount(t *testing.T) {
	decider := NewAccountDecider()

	t.Run("test a1 system owner", func(t *testing.T) {
		mock := mocks.New()

		err := initAccountTestGraph(t, mock)
		require.NoError(t, err)

		exp := time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
		mock.SetUser(mocks.A1SystemOwner)

		account := &model.Account{
			Name:  "A1",
			ATO:   "ato",
			MSPID: "A1MSP",
			Users: model.Users{
				SystemOwner:           "a1_system_owner",
				SystemAdministrator:   "a1_system_admin",
				AcquisitionSpecialist: "a1_acq_spec",
			},
			Status: "status",
			Assets: map[string]map[string]time.Time{
				"license1": {
					"k1": exp,
					"k2": exp,
				},
			},
		}

		err = decider.setup(mock.Stub)
		require.NoError(t, err)
		err = decider.filterAccount(account)
		require.NoError(t, err)
		require.Equal(t, "A1", account.Name)
		require.Equal(t, "ato", account.ATO)
		require.Equal(t, "A1MSP", account.MSPID)
		require.Equal(t, model.Status("status"), account.Status)
		require.Equal(t, model.Users{
			SystemOwner:           "a1_system_owner",
			SystemAdministrator:   "a1_system_admin",
			AcquisitionSpecialist: "a1_acq_spec",
		}, account.Users)
		require.Equal(t, map[string]map[string]time.Time{
			"license1": {
				"k1": exp,
				"k2": exp,
			},
		}, account.Assets)
	})

	t.Run("test a1 system admin", func(t *testing.T) {
		mock := mocks.New()

		err := initAccountTestGraph(t, mock)
		require.NoError(t, err)

		exp := time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
		mock.SetUser(mocks.A1SystemAdmin)

		account := &model.Account{
			Name:  "A1",
			ATO:   "ato",
			MSPID: "A1MSP",
			Users: model.Users{
				SystemOwner:           "a1_system_owner",
				SystemAdministrator:   "a1_system_admin",
				AcquisitionSpecialist: "a1_acq_spec",
			},
			Status: "status",
			Assets: map[string]map[string]time.Time{
				"license1": {
					"k1": exp,
					"k2": exp,
				},
			},
		}

		err = decider.setup(mock.Stub)
		require.NoError(t, err)
		err = decider.filterAccount(account)
		require.NoError(t, err)
		require.Equal(t, "A1", account.Name)
		require.Equal(t, "", account.ATO)
		require.Equal(t, "", account.MSPID)
		require.Equal(t, model.Status(""), account.Status)
		require.Equal(t, model.Users{}, account.Users)
		require.Equal(t, map[string]map[string]time.Time{
			"license1": {
				"k1": exp,
				"k2": exp,
			},
		}, account.Assets)
	})

	t.Run("test super", func(t *testing.T) {
		mock := mocks.New()

		err := initAccountTestGraph(t, mock)
		require.NoError(t, err)

		exp := time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
		mock.SetUser(mocks.Super)

		account := &model.Account{
			Name:  "A1",
			ATO:   "ato",
			MSPID: "A1MSP",
			Users: model.Users{
				SystemOwner:           "a1_system_owner",
				SystemAdministrator:   "a1_system_admin",
				AcquisitionSpecialist: "a1_acq_spec",
			},
			Status: "status",
			Assets: map[string]map[string]time.Time{
				"license1": {
					"k1": exp,
					"k2": exp,
				},
			},
		}

		err = decider.setup(mock.Stub)
		require.NoError(t, err)
		err = decider.filterAccount(account)
		require.NoError(t, err)
		require.Equal(t, "A1", account.Name)
		require.Equal(t, "ato", account.ATO)
		require.Equal(t, "A1MSP", account.MSPID)
		require.Equal(t, model.Status("status"), account.Status)
		require.Equal(t, model.Users{
			SystemOwner:           "a1_system_owner",
			SystemAdministrator:   "a1_system_admin",
			AcquisitionSpecialist: "a1_acq_spec",
		}, account.Users)
		require.Equal(t, map[string]map[string]time.Time{
			"license1": {
				"k1": exp,
				"k2": exp,
			},
		}, account.Assets)
	})
}

func initAccountTestGraph(t *testing.T, mock mocks.Mock) error {
	graph := memory.NewGraph()

	// configure the policy
	err := policy.Configure(graph)
	require.NoError(t, err)

	// add an account
	account := &model.Account{
		Name:  "A1",
		ATO:   "ato",
		MSPID: "A1MSP",
		Users: model.Users{
			SystemOwner:           "a1_system_owner",
			SystemAdministrator:   "a1_system_admin",
			AcquisitionSpecialist: "a1_acq_spec",
		},
		Status: "status",
		Assets: make(map[string]map[string]time.Time),
	}

	mock.SetGraphState(graph)

	// add account as the a1 system owner
	err = mock.SetUser(mocks.A1SystemOwner)
	require.NoError(t, err)

	accountDecider := NewAccountDecider()
	err = accountDecider.RequestAccount(mock.Stub, account)
	require.NoError(t, err)

	mock.SetGraphState(accountDecider.pap.Graph())

	return nil
}
