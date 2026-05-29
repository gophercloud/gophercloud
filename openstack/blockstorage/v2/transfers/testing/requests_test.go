package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/transfers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateTransfer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateTransfer(t, fakeServer)

	actual, err := transfers.Create(context.TODO(), client.ServiceClient(fakeServer), TransferRequest).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, TransferResponse, *actual)
}

func TestAcceptTransfer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAcceptTransfer(t, fakeServer)

	actual, err := transfers.Accept(context.TODO(), client.ServiceClient(fakeServer), TransferResponse.ID, AcceptRequest).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, AcceptResponse, *actual)
}

func TestDeleteTransfer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteTransfer(t, fakeServer)

	err := transfers.Delete(context.TODO(), client.ServiceClient(fakeServer), TransferResponse.ID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListTransfers(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListTransfers(t, fakeServer)

	expectedResponse := TransferListResponse
	expectedResponse[0].AuthKey = ""

	count := 0
	err := transfers.List(client.ServiceClient(fakeServer), &transfers.ListOpts{AllTenants: true}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := transfers.ExtractTransfers(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, expectedResponse, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListTransfersAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListTransfers(t, fakeServer)

	expectedResponse := TransferListResponse
	expectedResponse[0].AuthKey = ""

	allPages, err := transfers.List(client.ServiceClient(fakeServer), &transfers.ListOpts{AllTenants: true}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := transfers.ExtractTransfers(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedResponse, actual)
}

func TestGetTransfer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetTransfer(t, fakeServer)

	expectedResponse := TransferResponse
	expectedResponse.AuthKey = ""

	actual, err := transfers.Get(context.TODO(), client.ServiceClient(fakeServer), TransferResponse.ID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedResponse, *actual)
}
