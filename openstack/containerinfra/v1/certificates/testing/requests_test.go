package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/certificates"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetCertificates(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleGetCertificateSuccessfully(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	actual, err := certificates.Get(context.TODO(), sc, "d564b18a-2890-4152-be3d-e05d784ff72").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedCertificate, *actual)
}

func TestCreateCertificates(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateCertificateSuccessfully(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	opts := certificates.CreateOpts{
		BayUUID: "d564b18a-2890-4152-be3d-e05d784ff727",
		CSR:     "FAKE_CERTIFICATE_CSR",
	}

	actual, err := certificates.Create(context.TODO(), sc, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedCreateCertificateResponse, *actual)
}

func TestUpdateCertificates(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleUpdateCertificateSuccessfully(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	err := certificates.Update(context.TODO(), sc, "d564b18a-2890-4152-be3d-e05d784ff72").ExtractErr()
	th.AssertNoErr(t, err)
}
