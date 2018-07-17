// +build acceptance containerinfra

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/certificates"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCertificatesCRUD(t *testing.T) {
	client, err := clients.NewContainerInfraV1Client()
	th.AssertNoErr(t, err)

	clusterID := "8934d2d1-6bce-4ffa-a017-fb437777269d"

	certificate, err := certificates.Get(client, clusterID).Extract()
	th.AssertNoErr(t, err)
	t.Log(certificate.Pem)
}
