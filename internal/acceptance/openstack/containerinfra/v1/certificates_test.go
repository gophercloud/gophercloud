//go:build acceptance || containerinfra || certificates

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/certificates"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCertificatesCRUD(t *testing.T) {
	t.Skip("Test must be rewritten to drop hardcoded cluster ID")

	client, err := clients.NewContainerInfraV1Client()
	th.AssertNoErr(t, err)

	clusterUUID := "8934d2d1-6bce-4ffa-a017-fb437777269d"

	opts := certificates.CreateOpts{
		BayUUID: clusterUUID,
		CSR: "-----BEGIN CERTIFICATE REQUEST-----\n" +
			"MIIByjCCATMCAQAwgYkxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlh" +
			"MRYwFAYDVQQHEw1Nb3VudGFpbiBWaWV3MRMwEQYDVQQKEwpHb29nbGUgSW5jMR8w" +
			"HQYDVQQLExZJbmZvcm1hdGlvbiBUZWNobm9sb2d5MRcwFQYDVQQDEw53d3cuZ29v" +
			"Z2xlLmNvbTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEApZtYJCHJ4VpVXHfV" +
			"IlstQTlO4qC03hjX+ZkPyvdYd1Q4+qbAeTwXmCUKYHThVRd5aXSqlPzyIBwieMZr" +
			"WFlRQddZ1IzXAlVRDWwAo60KecqeAXnnUK+5fXoTI/UgWshre8tJ+x/TMHaQKR/J" +
			"cIWPhqaQhsJuzZbvAdGA80BLxdMCAwEAAaAAMA0GCSqGSIb3DQEBBQUAA4GBAIhl" +
			"4PvFq+e7ipARgI5ZM+GZx6mpCz44DTo0JkwfRDf+BtrsaC0q68eTf2XhYOsq4fkH" +
			"Q0uA0aVog3f5iJxCa3Hp5gxbJQ6zV6kJ0TEsuaaOhEko9sdpCoPOnRBm2i/XRD2D" +
			"6iNh8f8z0ShGsFqjDgFHyF3o+lUyj+UC6H1QW7bn\n" +
			"-----END CERTIFICATE REQUEST-----",
	}

	createResponse, err := certificates.Create(context.TODO(), client, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, opts.CSR, createResponse.CSR)

	certificate, err := certificates.Get(context.TODO(), client, clusterUUID).Extract()
	th.AssertNoErr(t, err)
	t.Log(certificate.PEM)

	err = certificates.Update(context.TODO(), client, clusterUUID).ExtractErr()
	th.AssertNoErr(t, err)
}
