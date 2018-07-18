package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/certificates"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestGetCertificates(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleGetCertificateSuccessfully(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	actual, err := certificates.Get(sc, "d564b18a-2890-4152-be3d-e05d784ff72").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedCertificate, *actual)
}
