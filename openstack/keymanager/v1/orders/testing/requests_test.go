package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/orders"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListOrders(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListOrdersSuccessfully(t, fakeServer)

	count := 0
	err := orders.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := orders.ExtractOrders(page)
		th.AssertNoErr(t, err)

		th.AssertDeepEquals(t, ExpectedOrdersSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, count, 1)
}

func TestListOrdersAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListOrdersSuccessfully(t, fakeServer)

	allPages, err := orders.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := orders.ExtractOrders(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedOrdersSlice, actual)
}

func TestGetOrder(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetOrderSuccessfully(t, fakeServer)

	actual, err := orders.Get(context.TODO(), client.ServiceClient(fakeServer), "46f73695-82bb-447a-bf96-6635f0fb0ce7").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, SecondOrder, *actual)
}

func TestCreateOrder(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateOrderSuccessfully(t, fakeServer)

	createOpts := orders.CreateOpts{
		Type: orders.KeyOrder,
		Meta: orders.MetaOpts{
			Algorithm:          "aes",
			BitLength:          256,
			Mode:               "cbc",
			PayloadContentType: "application/octet-stream",
		},
	}

	actual, err := orders.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, SecondOrder, *actual)
}

func TestDeleteOrder(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteOrderSuccessfully(t, fakeServer)

	res := orders.Delete(context.TODO(), client.ServiceClient(fakeServer), "46f73695-82bb-447a-bf96-6635f0fb0ce7")
	th.AssertNoErr(t, res.Err)
}
