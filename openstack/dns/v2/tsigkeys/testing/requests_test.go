package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/tsigkeys"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	count := 0
	err := tsigkeys.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := tsigkeys.ExtractTSIGKeys(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedTSIGKeysSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestListAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	allPages, err := tsigkeys.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allTSIGKeys, err := tsigkeys.ExtractTSIGKeys(allPages)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 2, len(allTSIGKeys))
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer)

	actual, err := tsigkeys.Get(context.TODO(), client.ServiceClient(fakeServer), "8add45a3-0f29-489f-854e-7609baf8d7a1").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstTSIGKey, actual)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSuccessfully(t, fakeServer)

	createOpts := tsigkeys.CreateOpts{
		Name:       "poolsecondarykey",
		Algorithm:  "hmac-sha256",
		Secret:     "my-base64-secret-example==",
		Scope:      "POOL",
		ResourceID: "adcc2fb6-7984-4453-a6f9-2cc2a24a38bb",
	}

	actual, err := tsigkeys.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &CreatedTSIGKey, actual)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateSuccessfully(t, fakeServer)

	updateOpts := tsigkeys.UpdateOpts{
		Name:   "updatedsecondarykey",
		Secret: "new-base64-secret-example==",
	}

	UpdatedTSIGKey := CreatedTSIGKey
	UpdatedTSIGKey.Name = "updatedsecondarykey"
	UpdatedTSIGKey.Secret = "new-base64-secret-example=="

	actual, err := tsigkeys.Update(context.TODO(), client.ServiceClient(fakeServer), UpdatedTSIGKey.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, UpdatedTSIGKey.Name, actual.Name)
	th.CheckEquals(t, UpdatedTSIGKey.Secret, actual.Secret)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSuccessfully(t, fakeServer)

	err := tsigkeys.Delete(context.TODO(), client.ServiceClient(fakeServer), "8add45a3-0f29-489f-854e-7609baf8d7a1").ExtractErr()
	th.AssertNoErr(t, err)
}
