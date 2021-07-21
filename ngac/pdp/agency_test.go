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
	decider := NewAgencyDecider()

	t.Run("test a1 system owner", func(t *testing.T) {
		mock := mocks.New()

		err := initAgencyTestGraph(t, mock)
		require.NoError(t, err)

		mock.SetUser(mocks.A1SystemOwner)

		err = decider.UploadATO(mock.Stub, "A1")
		require.NoError(t, err)
	})

	t.Run("test a1 system admin", func(t *testing.T) {
		mock := mocks.New()

		err := initAgencyTestGraph(t, mock)
		require.NoError(t, err)

		mock.SetUser(mocks.A1SystemAdmin)

		err = decider.UploadATO(mock.Stub, "A1")
		require.Error(t, err)
	})
}

func TestUpdateAgencyStatus(t *testing.T) {
	decider := NewAgencyDecider()

	t.Run("test a1 system owner", func(t *testing.T) {
		mock := mocks.New()

		err := initAgencyTestGraph(t, mock)
		require.NoError(t, err)

		mock.SetUser(mocks.A1SystemOwner)

		err = decider.UpdateAgencyStatus(mock.Stub, "A1", "test")
		require.Error(t, err)
	})

	t.Run("test a1 system admin", func(t *testing.T) {
		mock := mocks.New()

		err := initAgencyTestGraph(t, mock)
		require.NoError(t, err)

		mock.SetUser(mocks.A1SystemAdmin)

		err = decider.UpdateAgencyStatus(mock.Stub, "A1", "test")
		require.Error(t, err)
	})

	t.Run("test super", func(t *testing.T) {
		mock := mocks.New()

		err := initAgencyTestGraph(t, mock)
		require.NoError(t, err)

		mock.SetUser(mocks.Super)

		err = decider.UpdateAgencyStatus(mock.Stub, "A1", "test")
		require.NoError(t, err)
	})
}

func TestFilterAgency(t *testing.T) {
	decider := NewAgencyDecider()

	t.Run("test a1 system owner", func(t *testing.T) {
		mock := mocks.New()

		err := initAgencyTestGraph(t, mock)
		require.NoError(t, err)

		exp := time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
		mock.SetUser(mocks.A1SystemOwner)

		agency := &model.Agency{
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
		err = decider.filterAgency(agency)
		require.NoError(t, err)
		require.Equal(t, "A1", agency.Name)
		require.Equal(t, "ato", agency.ATO)
		require.Equal(t, "A1MSP", agency.MSPID)
		require.Equal(t, model.Status("status"), agency.Status)
		require.Equal(t, model.Users{
			SystemOwner:           "a1_system_owner",
			SystemAdministrator:   "a1_system_admin",
			AcquisitionSpecialist: "a1_acq_spec",
		}, agency.Users)
		require.Equal(t, map[string]map[string]time.Time{
			"license1": {
				"k1": exp,
				"k2": exp,
			},
		}, agency.Assets)
	})

	t.Run("test a1 system admin", func(t *testing.T) {
		mock := mocks.New()

		err := initAgencyTestGraph(t, mock)
		require.NoError(t, err)

		exp := time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
		mock.SetUser(mocks.A1SystemAdmin)

		agency := &model.Agency{
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
		err = decider.filterAgency(agency)
		require.NoError(t, err)
		require.Equal(t, "A1", agency.Name)
		require.Equal(t, "", agency.ATO)
		require.Equal(t, "", agency.MSPID)
		require.Equal(t, model.Status(""), agency.Status)
		require.Equal(t, model.Users{}, agency.Users)
		require.Equal(t, map[string]map[string]time.Time{
			"license1": {
				"k1": exp,
				"k2": exp,
			},
		}, agency.Assets)
	})

	t.Run("test super", func(t *testing.T) {
		mock := mocks.New()

		err := initAgencyTestGraph(t, mock)
		require.NoError(t, err)

		exp := time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
		mock.SetUser(mocks.Super)

		agency := &model.Agency{
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
		err = decider.filterAgency(agency)
		require.NoError(t, err)
		require.Equal(t, "A1", agency.Name)
		require.Equal(t, "ato", agency.ATO)
		require.Equal(t, "A1MSP", agency.MSPID)
		require.Equal(t, model.Status("status"), agency.Status)
		require.Equal(t, model.Users{
			SystemOwner:           "a1_system_owner",
			SystemAdministrator:   "a1_system_admin",
			AcquisitionSpecialist: "a1_acq_spec",
		}, agency.Users)
		require.Equal(t, map[string]map[string]time.Time{
			"license1": {
				"k1": exp,
				"k2": exp,
			},
		}, agency.Assets)
	})
}

func initAgencyTestGraph(t *testing.T, mock mocks.Mock) error {
	graph := memory.NewGraph()

	// configure the policy
	err := policy.Configure(graph)
	require.NoError(t, err)

	// add an account
	agency := model.Agency{
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
	mock.SetUser(mocks.A1SystemOwner)
	agencyDecider := NewAgencyDecider()
	err = agencyDecider.RequestAccount(mock.Stub, agency)
	require.NoError(t, err)

	mock.SetGraphState(agencyDecider.pap.Graph())

	return nil
}
