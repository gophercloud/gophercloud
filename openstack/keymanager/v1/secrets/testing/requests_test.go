package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListSecrets(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSecretsSuccessfully(t)

	count := 0
	err := secrets.List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := secrets.ExtractSecrets(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedSecretsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListSecretsAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSecretsSuccessfully(t)

	allPages, err := secrets.List(client.ServiceClient(), nil).AllPages()
	th.AssertNoErr(t, err)
	actual, err := secrets.ExtractSecrets(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedSecretsSlice, actual)
}

func TestGetSecret(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSecretSuccessfully(t)

	actual, err := secrets.Get(client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, FirstSecret, *actual)
}

func TestCreateSecret(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSecretSuccessfully(t)

	createOpts := secrets.CreateOpts{
		Algorithm:          "aes",
		BitLength:          256,
		Mode:               "cbc",
		Name:               "mysecret",
		Payload:            "foobar",
		PayloadContentType: "text/plain",
		SecretType:         secrets.OpaqueSecret,
	}

	actual, err := secrets.Create(client.ServiceClient(), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedCreateResult, *actual)
}

func TestDeleteSecret(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSecretSuccessfully(t)

	res := secrets.Delete(client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateSecret(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateSecretSuccessfully(t)

	updateOpts := secrets.UpdateOpts{
		Payload: "foobar",
	}

	err := secrets.Update(client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", updateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetPayloadSecret(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetPayloadSuccessfully(t)

	res := secrets.GetPayload(client.ServiceClient(), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c")
	th.AssertNoErr(t, res.Err)
	payload, err := res.Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, GetPayloadResult, string(payload))
}
