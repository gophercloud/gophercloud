package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/secrets"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListSecrets(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSecretsSuccessfully(t, fakeServer)

	count := 0
	err := secrets.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := secrets.ExtractSecrets(page)
		th.AssertNoErr(t, err)

		th.AssertDeepEquals(t, ExpectedSecretsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, count, 1)
}

func TestListSecretsAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSecretsSuccessfully(t, fakeServer)

	allPages, err := secrets.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := secrets.ExtractSecrets(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedSecretsSlice, actual)
}

func TestGetSecret(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSecretSuccessfully(t, fakeServer)

	actual, err := secrets.Get(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, FirstSecret, *actual)
}

func TestCreateSecret(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSecretSuccessfully(t, fakeServer)

	expiration := time.Date(2028, 6, 21, 2, 49, 48, 0, time.UTC)
	createOpts := secrets.CreateOpts{
		Algorithm:          "aes",
		BitLength:          256,
		Mode:               "cbc",
		Name:               "mysecret",
		Payload:            "foobar",
		PayloadContentType: "text/plain",
		SecretType:         secrets.OpaqueSecret,
		Expiration:         &expiration,
	}

	actual, err := secrets.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedCreateResult, *actual)
}

func TestDeleteSecret(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSecretSuccessfully(t, fakeServer)

	res := secrets.Delete(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateSecret(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateSecretSuccessfully(t, fakeServer)

	updateOpts := secrets.UpdateOpts{
		Payload: "foobar",
	}

	err := secrets.Update(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", updateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetPayloadSecret(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetPayloadSuccessfully(t, fakeServer)

	res := secrets.GetPayload(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", nil)
	th.AssertNoErr(t, res.Err)
	payload, err := res.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, GetPayloadResponse, string(payload))
}

func TestGetMetadataSuccessfully(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetMetadataSuccessfully(t, fakeServer)

	actual, err := secrets.GetMetadata(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedMetadata, actual)
}

func TestCreateMetadataSuccessfully(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateMetadataSuccessfully(t, fakeServer)

	createOpts := secrets.MetadataOpts{
		"foo":       "bar",
		"something": "something else",
	}

	actual, err := secrets.CreateMetadata(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedCreateMetadataResult, actual)
}

func TestGetMetadatumSuccessfully(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetMetadatumSuccessfully(t, fakeServer)

	actual, err := secrets.GetMetadatum(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", "foo").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedMetadatum, *actual)
}

func TestCreateMetadatumSuccessfully(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateMetadatumSuccessfully(t, fakeServer)

	createOpts := secrets.MetadatumOpts{
		Key:   "foo",
		Value: "bar",
	}

	err := secrets.CreateMetadatum(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", createOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUpdateMetadatumSuccessfully(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateMetadatumSuccessfully(t, fakeServer)

	updateOpts := secrets.MetadatumOpts{
		Key:   "foo",
		Value: "bar",
	}

	actual, err := secrets.UpdateMetadatum(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedMetadatum, *actual)
}

func TestDeleteMetadatumSuccessfully(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteMetadatumSuccessfully(t, fakeServer)

	err := secrets.DeleteMetadatum(context.TODO(), client.ServiceClient(fakeServer), "1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", "foo").ExtractErr()
	th.AssertNoErr(t, err)
}
