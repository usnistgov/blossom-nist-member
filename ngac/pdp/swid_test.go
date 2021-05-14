package pdp

import (
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/api/mocks"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/policy"
	"testing"
	"time"
)

func TestReportSwID(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	initSwidTestGraph(t, transactionContext, chaincodeStub)

	SetUser(transactionContext, Org1AdminCert(), "Org1MSP")
	agencyDecider := NewAgencyDecider()
	err := agencyDecider.UpdateAgencyStatus(transactionContext, "Org2", model.Approved)
	require.NoError(t, err)

	SetGraphState(t, chaincodeStub, agencyDecider.pap.Graph())

	SetUser(transactionContext, A1SystemAdminCert(), "Org2MSP")
	licenseDecider := NewLicenseDecider()
	err = licenseDecider.CheckoutLicense(transactionContext, "Org2", "test-license-id",
		map[string]time.Time{"1": time.Now()})
	require.NoError(t, err)

	SetGraphState(t, chaincodeStub, licenseDecider.pap.Graph())

	// report swid
	swidDecider := NewSwIDDecider()
	swid := &model.SwID{
		PrimaryTag:      "pt1",
		XML:             "xml",
		License:         "test-license-id",
		LicenseKey:      "1",
		LeaseExpiration: time.Time{},
	}
	err = swidDecider.ReportSwID(transactionContext, swid, "Org2")
	require.NoError(t, err)

	// report swid on license key that the user does not have access to
	swid = &model.SwID{
		PrimaryTag:      "pt1",
		XML:             "xml",
		License:         "test-license-id",
		LicenseKey:      "2",
		LeaseExpiration: time.Time{},
	}
	err = swidDecider.ReportSwID(transactionContext, swid, "Org2")
	require.Error(t, err)
}

func initSwidTestGraph(t *testing.T, ctx *mocks.TransactionContext, stub *mocks.ChaincodeStub) {
	graph := memory.NewGraph()

	// configure the policy
	err := policy.Configure(graph)
	require.NoError(t, err)

	// add an account
	agency := model.Agency{
		Name:  "Org2",
		ATO:   "ato",
		MSPID: "Org2MSP",
		Users: model.Users{
			SystemOwner:           "a1_system_owner",
			SystemAdministrator:   "a1_system_admin",
			AcquisitionSpecialist: "a1_acq_spec",
		},
		Status:   "status",
		Licenses: make(map[string]map[string]time.Time),
	}

	SetGraphState(t, stub, graph)

	// add account as the a1 system owner
	SetUser(ctx, A1SystemOwnerCert(), "Org2MSP")
	agencyDecider := NewAgencyDecider()
	err = agencyDecider.RequestAccount(ctx, agency)
	require.NoError(t, err)

	SetGraphState(t, stub, agencyDecider.pap.Graph())

	// set up the mock identity as the org1 admin
	SetUser(ctx, Org1AdminCert(), "Org1MSP")

	// create a test license
	license := &model.License{
		ID:             "test-license-id",
		Name:           "test-license",
		TotalAmount:    5,
		Available:      5,
		Cost:           20,
		OnboardingDate: time.Date(2021, 5, 12, 12, 0, 0, 0, time.Local),
		Expiration:     time.Date(2026, 5, 12, 12, 0, 0, 0, time.Local),
		AllKeys:        []string{"1", "2", "3", "4", "5"},
		AvailableKeys:  []string{"1", "2", "3", "4", "5"},
		CheckedOut:     make(map[string]map[string]time.Time),
	}

	licenseDecider := NewLicenseDecider()
	err = licenseDecider.OnboardLicense(ctx, license)
	require.NoError(t, err)

	SetGraphState(t, stub, licenseDecider.pap.Graph())
}
