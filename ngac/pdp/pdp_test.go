package pdp

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/stretchr/testify/require"
	"github.com/usnistgov/blossom/chaincode/api/mocks"
	"testing"
)

// Org1AdminCert returns the certificate for the user Org1 Admin:Org1MSP
func Org1AdminCert() *x509.Certificate {
	str := "-----BEGIN CERTIFICATE-----\n" +
		"MIIBzjCCAXSgAwIBAgIQSb323vGSInSfGJjW8I1gYTAKBggqhkjOPQQDAjASMRAw\n" +
		"DgYDVQQDEwdPcmcxIENBMB4XDTIxMDUwODAwMTExN1oXDTMxMDUwNjAwMTExN1ow\n" +
		"JTEOMAwGA1UECxMFYWRtaW4xEzARBgNVBAMTCk9yZzEgQWRtaW4wWTATBgcqhkjO\n" +
		"PQIBBggqhkjOPQMBBwNCAAR1TpubG65pIYToFf8xvm35XaEpn2ZIZ/vkG9QhX/PW\n" +
		"YYVIJf71XgLJ9cGKjyczoutEIVXqJmQB1WY7ZJ6SftV2o4GYMIGVMA4GA1UdDwEB\n" +
		"/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUHAwEwDAYDVR0TAQH/\n" +
		"BAIwADApBgNVHQ4EIgQgwv4cCnggolWnkkg4rcNumwNU8arzY2VjOMhoZQBMpTQw\n" +
		"KwYDVR0jBCQwIoAgfpDlwmohJG8C6wni21AOOe63BA6eiagEppuLgbd6PvMwCgYI\n" +
		"KoZIzj0EAwIDSAAwRQIhAIk3spHgBaZ9IZdVHK/CT1OwKhPxA1OFD2+keLgv0N1r\n" +
		"AiBECB68AZaKWdgjz3A+6ak8fHScGYAO/RoPdZeZibJ/nA==\n" +
		"-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	return cert
}

// A1SystemOwnerCert returns the certificate for the user a1_system_owner:Org2MSP
func A1SystemOwnerCert() *x509.Certificate {
	str := "-----BEGIN CERTIFICATE-----\n" +
		"MIICETCCAbegAwIBAgIUY46/5h6vqcPhxeLZolshAAWIoOgwCgYIKoZIzj0EAwIw\n" +
		"EjEQMA4GA1UEAxMHT3JnMiBDQTAeFw0yMTA1MTExOTA5MDBaFw0yMjA1MTExOTE0\n" +
		"MDBaMCsxDzANBgNVBAsTBmNsaWVudDEYMBYGA1UEAwwPYTFfc3lzdGVtX293bmVy\n" +
		"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE9dLsVb4kT4fKx6JFBg2fzrZI7Whz\n" +
		"fC7+eATQiq5Q9obxfoyzj8bwnDvN3vYJIoFugkBw/zF+Udaar7hM8G9zLKOB0TCB\n" +
		"zjAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUSpYneyw4\n" +
		"ynGa7NGPV9ImbKApRF0wKwYDVR0jBCQwIoAgPSGaMupMVSlW3LGlaKW2j2jYMwsK\n" +
		"pruo+JJAsevpEl8wYgYIKgMEBQYHCAEEVnsiYXR0cnMiOnsiaGYuQWZmaWxpYXRp\n" +
		"b24iOiIiLCJoZi5FbnJvbGxtZW50SUQiOiJhMV9zeXN0ZW1fb3duZXIiLCJoZi5U\n" +
		"eXBlIjoiY2xpZW50In19MAoGCCqGSM49BAMCA0gAMEUCIQDb4znXNqIPt6gZ4hST\n" +
		"/FhQoFp95QnAsmMh8LevKNeotQIgUvPaXkV6QcA9pcRtaCY3Q3c/ex8VgqrkzHCl\n" +
		"D8R6EnE=\n" +
		"-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	return cert
}

// A1SystemAdminCert returns the certificate for the user a1_system_admin:Org2MSP
func A1SystemAdminCert() *x509.Certificate {
	str := "-----BEGIN CERTIFICATE-----\n" +
		"MIICEDCCAbegAwIBAgIUXsc8SIL/O5DAs7EmtlqDSNUPiaQwCgYIKoZIzj0EAwIw\n" +
		"EjEQMA4GA1UEAxMHT3JnMiBDQTAeFw0yMTA1MTExOTA5MDBaFw0yMjA1MTExOTE0\n" +
		"MDBaMCsxDzANBgNVBAsTBmNsaWVudDEYMBYGA1UEAwwPYTFfc3lzdGVtX2FkbWlu\n" +
		"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAExK0hrd+UL3vMVqu/1o+PxqJ7zd4+\n" +
		"Dp8Ren9J9kwPt66EcEI7bDmLzOfLsra0wrSOI7tDFKCSKcbey2DZCCf9JaOB0TCB\n" +
		"zjAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUaWCsMWt9\n" +
		"EiWJFx8TW/SnRUExQCwwKwYDVR0jBCQwIoAgPSGaMupMVSlW3LGlaKW2j2jYMwsK\n" +
		"pruo+JJAsevpEl8wYgYIKgMEBQYHCAEEVnsiYXR0cnMiOnsiaGYuQWZmaWxpYXRp\n" +
		"b24iOiIiLCJoZi5FbnJvbGxtZW50SUQiOiJhMV9zeXN0ZW1fYWRtaW4iLCJoZi5U\n" +
		"eXBlIjoiY2xpZW50In19MAoGCCqGSM49BAMCA0cAMEQCIHToH9BR+n0jXTtjQ3X6\n" +
		"7VevbvB6akQP5CUa0ESM+jMYAiBztN2WrXvj8aUPJb/4Qdk91USFnJ28vo27xLqh\n" +
		"zaQ05g==\n-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	return cert
}

// A1AcqSpecCert returns the certificate for the user a1_acq_spec:Org2MSP
func A1AcqSpecCert() *x509.Certificate {
	str := "-----BEGIN CERTIFICATE-----\n" +
		"MIICCDCCAa+gAwIBAgIULRCOBjku2WVRVv5gAeySJU1kNhAwCgYIKoZIzj0EAwIw\n" +
		"EjEQMA4GA1UEAxMHT3JnMiBDQTAeFw0yMTA1MTExOTEwMDBaFw0yMjA1MTExOTE1\n" +
		"MDBaMCcxDzANBgNVBAsTBmNsaWVudDEUMBIGA1UEAwwLYTFfYWNxX3NwZWMwWTAT\n" +
		"BgcqhkjOPQIBBggqhkjOPQMBBwNCAATknFtPPYyG6JxZ3BTOP+faPgd7uQtB44f4\n" +
		"fxKikjFyNgCdJafduLlee1X0EuQ06Cgw9y6+Wd9CltyEW684WCc9o4HNMIHKMA4G\n" +
		"A1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBTxV4twzgyBk20b\n" +
		"BHbALaGrJdVi4DArBgNVHSMEJDAigCA9IZoy6kxVKVbcsaVopbaPaNgzCwqmu6j4\n" +
		"kkCx6+kSXzBeBggqAwQFBgcIAQRSeyJhdHRycyI6eyJoZi5BZmZpbGlhdGlvbiI6\n" +
		"IiIsImhmLkVucm9sbG1lbnRJRCI6ImExX2FjcV9zcGVjIiwiaGYuVHlwZSI6ImNs\n" +
		"aWVudCJ9fTAKBggqhkjOPQQDAgNHADBEAiBGWTL0ZJCG1p3GdXCOoxFZ943dTa+k\n" +
		"iTlmPe1QikW40QIgMpMW6WHeq6syroZjHHl+Ju8RfFp3vtJTWD/ALU72zbM=\n" +
		"-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	return cert
}

func TestCerts(t *testing.T) {
	if cert := Org1AdminCert(); cert.Subject.CommonName != "Org1 Admin" {
		t.Fatalf("expected Org1 Admin got %s", cert.Subject.CommonName)
	}

	if cert := A1SystemOwnerCert(); cert.Subject.CommonName != "a1_system_owner" {
		t.Fatalf("expected a1_system_owner got %s", cert.Subject.CommonName)
	}

	if cert := A1SystemAdminCert(); cert.Subject.CommonName != "a1_system_admin" {
		t.Fatalf("expected a1_system_admin got %s", cert.Subject.CommonName)
	}

	if cert := A1AcqSpecCert(); cert.Subject.CommonName != "a1_acq_spec" {
		t.Fatalf("expected a1_acq_spec got %s", cert.Subject.CommonName)
	}
}

func TestInitGraph(t *testing.T) {
	// create the mock chaincode stub and context
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	// create the mock client identity
	clientIdentity := &mocks.ClientIdentity{}
	clientIdentity.GetMSPIDReturns("Org1MSP", nil)
	clientIdentity.GetX509CertificateReturns(Org1AdminCert(), nil)
	transactionContext.GetClientIdentityReturns(clientIdentity)

	adminDecider := NewAdminDecider()
	err := adminDecider.InitGraph(transactionContext)
	// no error should occur since Org1 Admin has permission to init blossom
	require.NoError(t, err)

	// do the same with an unauthorized user - a1_system_owner
	clientIdentity = &mocks.ClientIdentity{}
	clientIdentity.GetMSPIDReturns("Org2MSP", nil)
	clientIdentity.GetX509CertificateReturns(A1SystemOwnerCert(), nil)
	transactionContext.GetClientIdentityReturns(clientIdentity)

	// make a call to init blossom
	err = adminDecider.InitGraph(transactionContext)
	// an error should occur because a1_system_owner cannot init blossom
	require.Error(t, err)
}
