package mocks

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/pkg/attrmgr"
	"github.com/usnistgov/blossom/chaincode/model"
)

func buildIdentity(certStr, mspid string) (*ClientIdentity, error) {
	c := &ClientIdentity{}
	c.GetMSPIDReturns(mspid, nil)

	block, _ := pem.Decode([]byte(certStr))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return &ClientIdentity{}, fmt.Errorf("failed to parse certificate: %s", err)
	}

	c.GetX509CertificateReturns(cert, nil)

	attrs, err := attrmgr.New().GetAttributesFromCert(cert)
	if err != nil {
		return nil, fmt.Errorf("error getting attributes from cert: %s", err)
	}

	c.GetAttributeValueReturns(attrs.Value(model.RoleAttribute))
	c.AssertAttributeValueCalls(func(attr string, value string) error {
		retrievedValue, ok, err := attrs.Value(attr)
		if err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("")
		}

		if retrievedValue != value {
			return fmt.Errorf("")
		}

		return nil
	})

	return c, nil
}

// Super returns the certificate for the user Org1 Admin:Org1MSP
func Super() (*ClientIdentity, error) {
	str := "-----BEGIN CERTIFICATE-----\nMIICbTCCAhSgAwIBAgIUCjoPBiSPEj8sK7jUC05kpQRGtOUwCgYIKoZIzj0EAwIw\ncDELMAkGA1UEBhMCVVMxFzAVBgNVBAgTDk5vcnRoIENhcm9saW5hMQ8wDQYDVQQH\nEwZEdXJoYW0xGTAXBgNVBAoTEG9yZzEuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh\nLm9yZzEuZXhhbXBsZS5jb20wHhcNMjIwNTE2MjMzMjAwWhcNMjMwNTE3MDAyMTAw\nWjAlMQ8wDQYDVQQLEwZjbGllbnQxEjAQBgNVBAMTCWFkbWludXNlcjBZMBMGByqG\nSM49AgEGCCqGSM49AwEHA0IABOA2FSuCUy+ZRPqMCsHC92xpTnI4BoD2Hm7XMs2Y\nrNhjbzexWza9NZTiwkgpd8+G5vq3c2XLlqouhgKa84SHzGmjgdYwgdMwDgYDVR0P\nAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFLuP4UjYW5SWZWOCIUGq\njnfrrXJ2MB8GA1UdIwQYMBaAFIk6JIULS50ZiEd73yMeaoOSjwVGMHMGCCoDBAUG\nBwgBBGd7ImF0dHJzIjp7ImJsb3Nzb20ucm9sZSI6ImFkbWluIiwiaGYuQWZmaWxp\nYXRpb24iOiIiLCJoZi5FbnJvbGxtZW50SUQiOiJhZG1pbnVzZXIiLCJoZi5UeXBl\nIjoiY2xpZW50In19MAoGCCqGSM49BAMCA0cAMEQCIE0Y3ISCTG1w+6ZLg9zCH+5l\nHSV3V+2CR+qnm52axzU5AiBmj6H5aNxa0TEX5ZvP1HIkQcqvx4U1uGke8DZAZXqr\nMA==\n-----END CERTIFICATE-----\n"
	return buildIdentity(str, "Org1MSP")
}

func UserInAdminMSPWithoutAdminRole() (*ClientIdentity, error) {
	str := "-----BEGIN CERTIFICATE-----\nMIIB8zCCAZmgAwIBAgIUYa0rLJ/xI/dBo096OiCQ5fYv/iQwCgYIKoZIzj0EAwIw\ncDELMAkGA1UEBhMCVVMxFzAVBgNVBAgTDk5vcnRoIENhcm9saW5hMQ8wDQYDVQQH\nEwZEdXJoYW0xGTAXBgNVBAoTEG9yZzEuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh\nLm9yZzEuZXhhbXBsZS5jb20wHhcNMjIwNTE2MjMzMjAwWhcNMjMwNTE3MDAxOTAw\nWjAhMQ8wDQYDVQQLEwZjbGllbnQxDjAMBgNVBAMTBWFkbWluMFkwEwYHKoZIzj0C\nAQYIKoZIzj0DAQcDQgAEhI8kekn7wbLqDfJG4tw/ZFUejpZje4i7oBcESan3dm17\nmJq5ZQs1Al9p/M6JO3DppdH4ELkwd4PxTnSSvHLP66NgMF4wDgYDVR0PAQH/BAQD\nAgeAMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFH19EN0tx95LX4Z30sKbXkZDgJ9g\nMB8GA1UdIwQYMBaAFIk6JIULS50ZiEd73yMeaoOSjwVGMAoGCCqGSM49BAMCA0gA\nMEUCIQCtazR/ulaekKNnA+Nu/N9o86Afy7M2ZRYy2Ro838iLrQIgFOFvvRHtxzkw\n+ZeVIO04xiZ3bLJmXO6aqqD4fjMwUNk=\n-----END CERTIFICATE-----\n"
	return buildIdentity(str, "Org1MSP")
}

// Org2SystemOwner returns the certificate for the user a1_system_owner:Org2MSP
func Org2SystemOwner() (*ClientIdentity, error) {
	str := "-----BEGIN CERTIFICATE-----\nMIICcDCCAhagAwIBAgIUA2Z426Yj8z1EnvQWEFOlebASO6QwCgYIKoZIzj0EAwIw\nbDELMAkGA1UEBhMCVUsxEjAQBgNVBAgTCUhhbXBzaGlyZTEQMA4GA1UEBxMHSHVy\nc2xleTEZMBcGA1UEChMQb3JnMi5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eub3Jn\nMi5leGFtcGxlLmNvbTAeFw0yMjA1MTYxNzUzMDBaFw0yMzA1MTYyMzA0MDBaMCUx\nDzANBgNVBAsTBmNsaWVudDESMBAGA1UEAxMJb3JnMnVzZXIxMFkwEwYHKoZIzj0C\nAQYIKoZIzj0DAQcDQgAEXerh6UZsMbAIvsstpBKpTzUtpGy4dErjN6w99biNDsVT\nL0Je6vEiKZRPZrqBUfqEt92vQEt2UX3V/Oxzvp6/2aOB3DCB2TAOBgNVHQ8BAf8E\nBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUzgJrJl3EtiCmF/nITaWnfrNU\nF7AwHwYDVR0jBBgwFoAUHMCz9Au/QVC4e68lU2xY3PBiGGgweQYIKgMEBQYHCAEE\nbXsiYXR0cnMiOnsiYmxvc3NvbS5yb2xlIjoiU3lzdGVtT3duZXIiLCJoZi5BZmZp\nbGlhdGlvbiI6IiIsImhmLkVucm9sbG1lbnRJRCI6Im9yZzJ1c2VyMSIsImhmLlR5\ncGUiOiJjbGllbnQifX0wCgYIKoZIzj0EAwIDSAAwRQIhAKdg35h+vEQJWEAiaeWJ\n4IAaXtZ0c8/NbAFON0RukdFXAiBo6GP7gFu9R+NYKq2G4EG/hpY1i4NA1+m+z9fP\nysrXAg==\n-----END CERTIFICATE-----\n"
	return buildIdentity(str, "Org2MSP")
}

// Org2SystemAdmin returns the certificate for the user a1_system_admin:Org2MSP
func Org2SystemAdmin() (*ClientIdentity, error) {
	str := "-----BEGIN CERTIFICATE-----\nMIICeDCCAh+gAwIBAgIUAs02aLwvcW3v1LVI85F+1SZM+mIwCgYIKoZIzj0EAwIw\nbDELMAkGA1UEBhMCVUsxEjAQBgNVBAgTCUhhbXBzaGlyZTEQMA4GA1UEBxMHSHVy\nc2xleTEZMBcGA1UEChMQb3JnMi5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eub3Jn\nMi5leGFtcGxlLmNvbTAeFw0yMjA1MTYxNzUzMDBaFw0yMzA1MTYyMzA0MDBaMCUx\nDzANBgNVBAsTBmNsaWVudDESMBAGA1UEAxMJb3JnMnVzZXIyMFkwEwYHKoZIzj0C\nAQYIKoZIzj0DAQcDQgAEYjzbsttU0OCDQxe0FEG+DlhjZ/U1DL7FqnExtBWfvh4G\nUi+rQn3Tfq/OfO6mdH54aozdl37K3g0btDBg3jQx2qOB5TCB4jAOBgNVHQ8BAf8E\nBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUdPhoYoEgkL9yWEdMokbC35aa\nU08wHwYDVR0jBBgwFoAUHMCz9Au/QVC4e68lU2xY3PBiGGgwgYEGCCoDBAUGBwgB\nBHV7ImF0dHJzIjp7ImJsb3Nzb20ucm9sZSI6IlN5c3RlbUFkbWluaXN0cmF0b3Ii\nLCJoZi5BZmZpbGlhdGlvbiI6IiIsImhmLkVucm9sbG1lbnRJRCI6Im9yZzJ1c2Vy\nMiIsImhmLlR5cGUiOiJjbGllbnQifX0wCgYIKoZIzj0EAwIDRwAwRAIgAgrS6GBT\nJ0p9C12nbun2stn0GO9C9aVWMR4FgtoBJq4CICAoxGLLik6HIiMgwX20iQ86ufWb\ne+meWU73BdL4OkdS\n-----END CERTIFICATE-----\n"
	return buildIdentity(str, "Org2MSP")
}

// Org2AcqSpec returns the certificate for the user a1_acq_spec:Org2MSP
func Org2AcqSpec() (*ClientIdentity, error) {
	str := "-----BEGIN CERTIFICATE-----\nMIICezCCAiGgAwIBAgIUa48kTQb2XI35EY0ilJs/mwqhqC0wCgYIKoZIzj0EAwIw\nbDELMAkGA1UEBhMCVUsxEjAQBgNVBAgTCUhhbXBzaGlyZTEQMA4GA1UEBxMHSHVy\nc2xleTEZMBcGA1UEChMQb3JnMi5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eub3Jn\nMi5leGFtcGxlLmNvbTAeFw0yMjA1MTYxNzUzMDBaFw0yMzA1MTYyMzA0MDBaMCUx\nDzANBgNVBAsTBmNsaWVudDESMBAGA1UEAxMJb3JnMnVzZXIzMFkwEwYHKoZIzj0C\nAQYIKoZIzj0DAQcDQgAE9oql2GjPZ7o56H/JbLtVIWmW8E+zgY5VOeb0U2N4PtiF\neebn0DKExZ6eeaDfR7onXe1vrKTZdV1nSDtbyV4A8KOB5zCB5DAOBgNVHQ8BAf8E\nBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUWSf1VvVZAYWix0qBdqQ4RZXb\nWJkwHwYDVR0jBBgwFoAUHMCz9Au/QVC4e68lU2xY3PBiGGgwgYMGCCoDBAUGBwgB\nBHd7ImF0dHJzIjp7ImJsb3Nzb20ucm9sZSI6IkFjcXVpc2l0aW9uU3BlY2lhbGlz\ndCIsImhmLkFmZmlsaWF0aW9uIjoiIiwiaGYuRW5yb2xsbWVudElEIjoib3JnMnVz\nZXIzIiwiaGYuVHlwZSI6ImNsaWVudCJ9fTAKBggqhkjOPQQDAgNIADBFAiEAnF7B\nZU/AoetM/DYM9brhFxgwo39NgIRv7pRIEnLUKYgCIDORyIDSLkhoKnb952W6Ez69\nL7kZOaCwbzGXSWG6+22M\n-----END CERTIFICATE-----\n"
	return buildIdentity(str, "Org2MSP")
}

// Org3SystemOwner returns the certificate for the user a1_system_owner:Org2MSP
func Org3SystemOwner() (*ClientIdentity, error) {
	str := "-----BEGIN CERTIFICATE-----\nMIICdDCCAhugAwIBAgIUMqT4PW5CZ5XlDkJBpVtQivcn2Q0wCgYIKoZIzj0EAwIw\ncTELMAkGA1UEBhMCVVMxFzAVBgNVBAgTDk5vcnRoIENhcm9saW5hMRAwDgYDVQQH\nEwdSYWxlaWdoMRkwFwYDVQQKExBvcmczLmV4YW1wbGUuY29tMRwwGgYDVQQDExNj\nYS5vcmczLmV4YW1wbGUuY29tMB4XDTIyMDUxNjE3NTQwMFoXDTIzMDUxNjIzMDQw\nMFowJTEPMA0GA1UECxMGY2xpZW50MRIwEAYDVQQDEwlvcmczdXNlcjEwWTATBgcq\nhkjOPQIBBggqhkjOPQMBBwNCAARHhDyGVgM5BGc0L2f9uFNguM1F/Lo0f5M+jV2M\nY0TEMLiHQfZG6TinZVsGU4rref9dkMUq0GXHEyufGVuQPYRBo4HcMIHZMA4GA1Ud\nDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBTGyYxgEXdDudXAFqF4\nWZJkiuZGYjAfBgNVHSMEGDAWgBQrAFwJkPnMMSLTvd13Ix+eWVgg2jB5BggqAwQF\nBgcIAQRteyJhdHRycyI6eyJibG9zc29tLnJvbGUiOiJTeXN0ZW1Pd25lciIsImhm\nLkFmZmlsaWF0aW9uIjoiIiwiaGYuRW5yb2xsbWVudElEIjoib3JnM3VzZXIxIiwi\naGYuVHlwZSI6ImNsaWVudCJ9fTAKBggqhkjOPQQDAgNHADBEAiAA7Mbf19lRVjvt\nH/7ylHc01qPOVNGqqZ0Fr/OvBCVwdwIgO3Pu5wUG0NyIZWMGHLk17Tt0F6j5OEL+\n6cw47hYqu9s=\n-----END CERTIFICATE-----\n"
	return buildIdentity(str, "Org3MSP")
}

// Org3SystemAdmin returns the certificate for the user a1_system_admin:Org2MSP
func Org3SystemAdmin() (*ClientIdentity, error) {
	str := "-----BEGIN CERTIFICATE-----\nMIICfjCCAiSgAwIBAgIUPFTE1g93av4t7hYFt5z4BSc8HYowCgYIKoZIzj0EAwIw\ncTELMAkGA1UEBhMCVVMxFzAVBgNVBAgTDk5vcnRoIENhcm9saW5hMRAwDgYDVQQH\nEwdSYWxlaWdoMRkwFwYDVQQKExBvcmczLmV4YW1wbGUuY29tMRwwGgYDVQQDExNj\nYS5vcmczLmV4YW1wbGUuY29tMB4XDTIyMDUxNjE3NTQwMFoXDTIzMDUxNjIzMDQw\nMFowJTEPMA0GA1UECxMGY2xpZW50MRIwEAYDVQQDEwlvcmczdXNlcjIwWTATBgcq\nhkjOPQIBBggqhkjOPQMBBwNCAASmmx/CkwodIuoqlWWcL5mgKP1ZqXg5CNXWlWJg\nERICfv3+oru6APdn7TjT/udrGA3jFQfNA9tlTMampB8+FAzuo4HlMIHiMA4GA1Ud\nDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBTc7pQwiQI8CUveZvpb\nm5laZc6eoTAfBgNVHSMEGDAWgBQrAFwJkPnMMSLTvd13Ix+eWVgg2jCBgQYIKgME\nBQYHCAEEdXsiYXR0cnMiOnsiYmxvc3NvbS5yb2xlIjoiU3lzdGVtQWRtaW5pc3Ry\nYXRvciIsImhmLkFmZmlsaWF0aW9uIjoiIiwiaGYuRW5yb2xsbWVudElEIjoib3Jn\nM3VzZXIyIiwiaGYuVHlwZSI6ImNsaWVudCJ9fTAKBggqhkjOPQQDAgNIADBFAiEA\n2B8ozbFhdqDGumcILHbtOVEv2imW0MqUAGx/+NOkdHcCIGrYxC3LO1usFfNH51/c\njijm1W2D0cJ78EerfYTQXzPu\n-----END CERTIFICATE-----\n"
	return buildIdentity(str, "Org3MSP")
}

// Org3AcqSpec returns the certificate for the user a1_acq_spec:Org2MSP
func Org3AcqSpec() (*ClientIdentity, error) {
	str := "-----BEGIN CERTIFICATE-----\nMIICgDCCAiagAwIBAgIUBPM4p2Bi4bHMlvk25PCYydsVlkIwCgYIKoZIzj0EAwIw\ncTELMAkGA1UEBhMCVVMxFzAVBgNVBAgTDk5vcnRoIENhcm9saW5hMRAwDgYDVQQH\nEwdSYWxlaWdoMRkwFwYDVQQKExBvcmczLmV4YW1wbGUuY29tMRwwGgYDVQQDExNj\nYS5vcmczLmV4YW1wbGUuY29tMB4XDTIyMDUxNjE3NTQwMFoXDTIzMDUxNjIzMDQw\nMFowJTEPMA0GA1UECxMGY2xpZW50MRIwEAYDVQQDEwlvcmczdXNlcjMwWTATBgcq\nhkjOPQIBBggqhkjOPQMBBwNCAAQT7F5IOHkQeGOtrNU1iGzG2e3B/6erIIZmcwbL\nmnTwGDYjnvsOMyf3XnyO0on/v+DbbXtKhUj8ONC+EK0hyIWho4HnMIHkMA4GA1Ud\nDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBRphgVyZ/k4dZWOfswO\ngu3wL/pHTzAfBgNVHSMEGDAWgBQrAFwJkPnMMSLTvd13Ix+eWVgg2jCBgwYIKgME\nBQYHCAEEd3siYXR0cnMiOnsiYmxvc3NvbS5yb2xlIjoiQWNxdWlzaXRpb25TcGVj\naWFsaXN0IiwiaGYuQWZmaWxpYXRpb24iOiIiLCJoZi5FbnJvbGxtZW50SUQiOiJv\ncmczdXNlcjMiLCJoZi5UeXBlIjoiY2xpZW50In19MAoGCCqGSM49BAMCA0gAMEUC\nIQCnpR2alnl3loAJ7Lm9ETKFJpaZfj6u3gcq4joSjsfamAIgU8jYF1xmWB7Yaw+9\neZPogDVuoVVwofsyTWy1G9bXcYE=\n-----END CERTIFICATE-----\n"
	return buildIdentity(str, "Org3MSP")
}
