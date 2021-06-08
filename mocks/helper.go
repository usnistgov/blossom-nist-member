package mocks

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

type (
	Mock struct {
		Stub *ChaincodeStub
		Ctx  *TransactionContext
	}
)

func New() Mock {
	chaincodeStub := &ChaincodeStub{}
	transactionContext := &TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	return Mock{
		Stub: chaincodeStub,
		Ctx:  transactionContext,
	}
}

func (c *Mock) SetUser(clientIdentity cid.ClientIdentity) {
	c.Ctx.GetClientIdentityReturns(clientIdentity)
}

// Super returns the certificate for the user Org1 Admin:Org1MSP
func Super() *ClientIdentity {
	str := "-----BEGIN CERTIFICATE-----\nMIIB/zCCAaagAwIBAgIUB2Mv6dyzlt/E4gswDwcZdmcJzfAwCgYIKoZIzj0EAwIw\nFTETMBEGA1UEAxMKQmxvc3NvbSBDQTAeFw0yMTA2MDgxNjM1MDBaFw0yMjA2MDgx\nNjQwMDBaMCExDzANBgNVBAsTBmNsaWVudDEOMAwGA1UEAxMFc3VwZXIwWTATBgcq\nhkjOPQIBBggqhkjOPQMBBwNCAAQzCY03xo1c9UqApo0fXtyfKpiT690tYD20N3S7\ns9/rdoWA2bbaWSzLKwSE80ev86DrCnjzhAqKjs/Yc/fYrTO/o4HHMIHEMA4GA1Ud\nDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBQO/UFdav4KK7Orriew\n3tjVtDRNwDArBgNVHSMEJDAigCCxmDQioiXBPziZIvTFt8p1ByaINs49yGIJedkI\nsfbw1zBYBggqAwQFBgcIAQRMeyJhdHRycyI6eyJoZi5BZmZpbGlhdGlvbiI6IiIs\nImhmLkVucm9sbG1lbnRJRCI6InN1cGVyIiwiaGYuVHlwZSI6ImNsaWVudCJ9fTAK\nBggqhkjOPQQDAgNHADBEAiB5I0izNDzxZHFn4HI5T2S8EMQMBSJlIylfaRGr2Wq5\nfgIgV8+FzIiQ7MxwwpTuU3lw2A2/yLGfZvAtjd2bzjv4tEA=\n-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	clientIdentity := &ClientIdentity{}
	clientIdentity.GetMSPIDReturns("BlossomMSP", nil)
	clientIdentity.GetX509CertificateReturns(cert, nil)

	return clientIdentity
}

// A1SystemOwner returns the certificate for the user a1_system_owner:Org2MSP
func A1SystemOwner() *ClientIdentity {
	str := "-----BEGIN CERTIFICATE-----\nMIICDjCCAbWgAwIBAgIUImZU9tOTjP+zavxejU3ZVWeIIU8wCgYIKoZIzj0EAwIw\nEDEOMAwGA1UEAxMFQTEgQ0EwHhcNMjEwNjA4MTYzNjAwWhcNMjIwNjA4MTY0MTAw\nWjArMQ8wDQYDVQQLEwZjbGllbnQxGDAWBgNVBAMMD2ExX3N5c3RlbV9vd25lcjBZ\nMBMGByqGSM49AgEGCCqGSM49AwEHA0IABLPr4DyEXa0lYQ7a27Wgg3q8clYa3lXQ\nB7JLHEOS9egDFZvKTvm0GNMrCzSYO/nHS6f92W9olZmWeRMwY+VmrkajgdEwgc4w\nDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFJozlajikYWp\nmEE5lhlC9i8FkOg6MCsGA1UdIwQkMCKAILZfTTDhOMj74tq6jntkZrAL7BKvZ3HV\nWBBVnvIJv/ovMGIGCCoDBAUGBwgBBFZ7ImF0dHJzIjp7ImhmLkFmZmlsaWF0aW9u\nIjoiIiwiaGYuRW5yb2xsbWVudElEIjoiYTFfc3lzdGVtX293bmVyIiwiaGYuVHlw\nZSI6ImNsaWVudCJ9fTAKBggqhkjOPQQDAgNHADBEAiAj7DNlUbXSJYAaMjC1AmeK\nkOJtU4rwfNpSj9nGXTNBhAIgcxpW1zVlMvIAoqipThlyi9roWySNkDYwhOffIE5E\nZLc=\n-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	clientIdentity := &ClientIdentity{}
	clientIdentity.GetMSPIDReturns("A1MSP", nil)
	clientIdentity.GetX509CertificateReturns(cert, nil)

	return clientIdentity
}

// A1SystemAdmin returns the certificate for the user a1_system_admin:Org2MSP
func A1SystemAdmin() *ClientIdentity {
	str := "-----BEGIN CERTIFICATE-----\nMIICDzCCAbWgAwIBAgIUOzq3ysq1CDdaLgAxGG+FHjihZMUwCgYIKoZIzj0EAwIw\nEDEOMAwGA1UEAxMFQTEgQ0EwHhcNMjEwNjA4MTYzNjAwWhcNMjIwNjA4MTY0MTAw\nWjArMQ8wDQYDVQQLEwZjbGllbnQxGDAWBgNVBAMMD2ExX3N5c3RlbV9hZG1pbjBZ\nMBMGByqGSM49AgEGCCqGSM49AwEHA0IABKLc9kUhM+/W/8ORsox5AodW7lIGWZIk\n/DAFZqmXZCiCwPvR0FkTsl2I01mmfiejBzpeph32BRxg/y+x2BIyerSjgdEwgc4w\nDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFIgLWa1ChlcQ\nhMTCwtjQlCrXsDaAMCsGA1UdIwQkMCKAILZfTTDhOMj74tq6jntkZrAL7BKvZ3HV\nWBBVnvIJv/ovMGIGCCoDBAUGBwgBBFZ7ImF0dHJzIjp7ImhmLkFmZmlsaWF0aW9u\nIjoiIiwiaGYuRW5yb2xsbWVudElEIjoiYTFfc3lzdGVtX2FkbWluIiwiaGYuVHlw\nZSI6ImNsaWVudCJ9fTAKBggqhkjOPQQDAgNIADBFAiEA5tH4NcoV8S3kH6zzLr5Y\nxHct1q6TLMxgTpIO+1l3/eYCIE22P1If8IkLQALevFb6a9riIdAWugomIRUQ+pl6\ndP1D\n-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	clientIdentity := &ClientIdentity{}
	clientIdentity.GetMSPIDReturns("A1MSP", nil)
	clientIdentity.GetX509CertificateReturns(cert, nil)

	return clientIdentity
}

// A1AcqSpec returns the certificate for the user a1_acq_spec:Org2MSP
func A1AcqSpec() *ClientIdentity {
	str := "-----BEGIN CERTIFICATE-----\nMIICBjCCAa2gAwIBAgIUDbWFGbDK6iWkfKhsLzT3thzxCZcwCgYIKoZIzj0EAwIw\nEDEOMAwGA1UEAxMFQTEgQ0EwHhcNMjEwNjA4MTYzNjAwWhcNMjIwNjA4MTY0MTAw\nWjAnMQ8wDQYDVQQLEwZjbGllbnQxFDASBgNVBAMMC2ExX2FjcV9zcGVjMFkwEwYH\nKoZIzj0CAQYIKoZIzj0DAQcDQgAEbay46lcreU7s4+OvI2I6c+7X61Xz12DiMCkx\nQY2IZLNDEJN+DfwTMP/kPgRibAf7kVT9FwvdUEUmr4Wf2p2NcqOBzTCByjAOBgNV\nHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUOnpUNlvfq7tTOr4k\nPG6w245UE10wKwYDVR0jBCQwIoAgtl9NMOE4yPvi2rqOe2RmsAvsEq9ncdVYEFWe\n8gm/+i8wXgYIKgMEBQYHCAEEUnsiYXR0cnMiOnsiaGYuQWZmaWxpYXRpb24iOiIi\nLCJoZi5FbnJvbGxtZW50SUQiOiJhMV9hY3Ffc3BlYyIsImhmLlR5cGUiOiJjbGll\nbnQifX0wCgYIKoZIzj0EAwIDRwAwRAIgNWFbO85KyPtm27q2jzoEshf7qAKOA8Yk\nquH/MX2zzYgCIFRsLYsmRl90p+vgXd3SPFC7DIkVuZPeZE2YvEs0LLcF\n-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	clientIdentity := &ClientIdentity{}
	clientIdentity.GetMSPIDReturns("A1MSP", nil)
	clientIdentity.GetX509CertificateReturns(cert, nil)

	return clientIdentity
}

// A2SystemOwner returns the certificate for the user a1_system_owner:Org2MSP
func A2SystemOwner() *ClientIdentity {
	str := "-----BEGIN CERTIFICATE-----\nMIICDjCCAbWgAwIBAgIUbe6WxqOMPQERISaz5eetFUt/xvIwCgYIKoZIzj0EAwIw\nEDEOMAwGA1UEAxMFQTIgQ0EwHhcNMjEwNjA4MTYzNzAwWhcNMjIwNjA4MTY0MjAw\nWjArMQ8wDQYDVQQLEwZjbGllbnQxGDAWBgNVBAMMD2EyX3N5c3RlbV9vd25lcjBZ\nMBMGByqGSM49AgEGCCqGSM49AwEHA0IABGsxwNjB3Rkdcml6NW1ysRTiivA4lQ+8\nje0PPP9iIwacRD5Oe+Bk80SH27AgSFo9tlzmWkNIf9Q5VVCTwBw8ZvOjgdEwgc4w\nDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFOicMHpU/QBQ\ntAV+g9aGNPMjSfEBMCsGA1UdIwQkMCKAIL97QI5L46TlIQuBdm/Nd/AIG7RR8w5R\ngnOC1D8DBWTkMGIGCCoDBAUGBwgBBFZ7ImF0dHJzIjp7ImhmLkFmZmlsaWF0aW9u\nIjoiIiwiaGYuRW5yb2xsbWVudElEIjoiYTJfc3lzdGVtX293bmVyIiwiaGYuVHlw\nZSI6ImNsaWVudCJ9fTAKBggqhkjOPQQDAgNHADBEAiA2cZAF0k0fCEQoaZXOBRXN\noLt7wPZOUNcaVDRNPzxE7QIgD8PNKI2ZwehMYpcRFxb6FwUypVq2MLolJAkXxuv0\nOQ4=\n-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	clientIdentity := &ClientIdentity{}
	clientIdentity.GetMSPIDReturns("A2MSP", nil)
	clientIdentity.GetX509CertificateReturns(cert, nil)

	return clientIdentity
}

// A2SystemAdmin returns the certificate for the user a1_system_admin:Org2MSP
func A2SystemAdmin() *ClientIdentity {
	str := "-----BEGIN CERTIFICATE-----\nMIICDjCCAbWgAwIBAgIUEVSqLqryBs1l/Alsf52VUjoQlLgwCgYIKoZIzj0EAwIw\nEDEOMAwGA1UEAxMFQTIgQ0EwHhcNMjEwNjA4MTYzNzAwWhcNMjIwNjA4MTY0MjAw\nWjArMQ8wDQYDVQQLEwZjbGllbnQxGDAWBgNVBAMMD2EyX3N5c3RlbV9hZG1pbjBZ\nMBMGByqGSM49AgEGCCqGSM49AwEHA0IABOYRa349Cr+vjH1anwCYEPTVi8iNniAM\nOWZlcXLtT5Fzr46VEwk9L5cowKe2EJyOE0TYeMY8B29/sQxieXyxRA6jgdEwgc4w\nDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFMlclyWudW/S\nyd7Q0eNFhmPVB7u+MCsGA1UdIwQkMCKAIL97QI5L46TlIQuBdm/Nd/AIG7RR8w5R\ngnOC1D8DBWTkMGIGCCoDBAUGBwgBBFZ7ImF0dHJzIjp7ImhmLkFmZmlsaWF0aW9u\nIjoiIiwiaGYuRW5yb2xsbWVudElEIjoiYTJfc3lzdGVtX2FkbWluIiwiaGYuVHlw\nZSI6ImNsaWVudCJ9fTAKBggqhkjOPQQDAgNHADBEAiBJ0wWRJHKYmXLQSkadzW9N\nInBzec1pJYJKCDF3HOKLkQIgDr9Tjp4b+sjcPzeB8PZQye+oKeCunGY9VpJEd7Af\n7Z8=\n-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	clientIdentity := &ClientIdentity{}
	clientIdentity.GetMSPIDReturns("A2MSP", nil)
	clientIdentity.GetX509CertificateReturns(cert, nil)

	return clientIdentity
}

// A2AcqSpec returns the certificate for the user a1_acq_spec:Org2MSP
func A2AcqSpec() *ClientIdentity {
	str := "-----BEGIN CERTIFICATE-----\nMIICBzCCAa2gAwIBAgIUZ3ebIFQkSTANq88XLbWsSfCcSyMwCgYIKoZIzj0EAwIw\nEDEOMAwGA1UEAxMFQTIgQ0EwHhcNMjEwNjA4MTYzNzAwWhcNMjIwNjA4MTY0MjAw\nWjAnMQ8wDQYDVQQLEwZjbGllbnQxFDASBgNVBAMMC2EyX2FjcV9zcGVjMFkwEwYH\nKoZIzj0CAQYIKoZIzj0DAQcDQgAEcr3+DQjStE3crMULpQ6o1ceU/YpJqHfKENdA\nFOW9dbVkQBn8c/F+fWxnA99rpVHO040S//gqsjlv2fuyhIwD4aOBzTCByjAOBgNV\nHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUVPfcn69S91oSML1Z\n6IyOja+Df68wKwYDVR0jBCQwIoAgv3tAjkvjpOUhC4F2b8138AgbtFHzDlGCc4LU\nPwMFZOQwXgYIKgMEBQYHCAEEUnsiYXR0cnMiOnsiaGYuQWZmaWxpYXRpb24iOiIi\nLCJoZi5FbnJvbGxtZW50SUQiOiJhMl9hY3Ffc3BlYyIsImhmLlR5cGUiOiJjbGll\nbnQifX0wCgYIKoZIzj0EAwIDSAAwRQIhAMl2cmrf3BULNjr+E6czXqD5Md+j6UNC\npQ3kc5K5wo/cAiBZ9yVWGCIE53WECDMMqFZvkfsyL4ChaFZvSb0PQzZC2A==\n-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	clientIdentity := &ClientIdentity{}
	clientIdentity.GetMSPIDReturns("A2MSP", nil)
	clientIdentity.GetX509CertificateReturns(cert, nil)

	return clientIdentity
}

func (c *Mock) SetGraphState(graph pip.Graph) {
	graphBytes, err := json.Marshal(graph)
	if err != nil {
		panic(err)
	}

	c.Stub.GetStateReturns(graphBytes, nil)
}
