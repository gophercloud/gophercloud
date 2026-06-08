package testing

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/services"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateSuccessful(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateServiceSuccessfully(t, fakeServer)

	createOpts := services.CreateOpts{
		Name:        "service-two",
		Description: "Service Two",
		Type:        "compute",
		Extra: map[string]any{
			"email": "service@example.com",
		},
	}

	actual, err := services.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondService, *actual)
}

func TestListServices(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListServicesSuccessfully(t, fakeServer)

	count := 0
	err := services.List(client.ServiceClient(fakeServer), services.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := services.ExtractServices(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedServicesSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestListServicesAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListServicesSuccessfully(t, fakeServer)

	allPages, err := services.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedServicesSlice, actual)
	th.AssertEquals(t, "service-one@example.com", ExpectedServicesSlice[0].Extra["email"])
	th.AssertEquals(t, "service@example.com", ExpectedServicesSlice[1].Extra["email"])
}

func TestGetSuccessful(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetServiceSuccessfully(t, fakeServer)

	actual, err := services.Get(context.TODO(), client.ServiceClient(fakeServer), "9876").Extract()

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondService, *actual)
	th.AssertEquals(t, "service@example.com", SecondService.Extra["email"])
}

func TestUpdateSuccessful(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateServiceSuccessfully(t, fakeServer)

	updatedDescription := "Service Two Updated"
	updateOpts := services.UpdateOpts{
		Description: &updatedDescription,
		Type:        "compute2",
	}
	actual, err := services.Update(context.TODO(), client.ServiceClient(fakeServer), "9876", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondServiceUpdated, *actual)
	th.AssertEquals(t, "Service Two Updated", SecondServiceUpdated.Description)
}

func TestDeleteSuccessful(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/services/12345", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := services.Delete(context.TODO(), client.ServiceClient(fakeServer), "12345")
	th.AssertNoErr(t, res.Err)
}
