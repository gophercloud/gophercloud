package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumetransfers"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreateTransfer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateTransfer(t)

	actual, err := volumetransfers.Create(client.ServiceClient(), TransferRequest).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, TransferResponse, *actual)
}

func TestAcceptTransfer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAcceptTransfer(t)

	actual, err := volumetransfers.Accept(client.ServiceClient(), TransferResponse.ID, AcceptRequest).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, AcceptResponse, *actual)
}

func TestDeleteTransfer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteTransfer(t)

	err := volumetransfers.Delete(client.ServiceClient(), TransferResponse.ID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListTransfers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTransfers(t)

	expectedResponse := TransferListResponse
	expectedResponse[0].AuthKey = ""

	count := 0
	err := volumetransfers.List(client.ServiceClient(), &volumetransfers.ListOpts{AllTenants: true}).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := volumetransfers.ExtractTransfers(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, expectedResponse, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListTransfersAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTransfers(t)

	expectedResponse := TransferListResponse
	expectedResponse[0].AuthKey = ""

	allPages, err := volumetransfers.List(client.ServiceClient(), &volumetransfers.ListOpts{AllTenants: true}).AllPages()
	th.AssertNoErr(t, err)
	actual, err := volumetransfers.ExtractTransfers(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedResponse, actual)
}

func TestGetTransfer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetTransfer(t)

	expectedResponse := TransferResponse
	expectedResponse.AuthKey = ""

	actual, err := volumetransfers.Get(client.ServiceClient(), TransferResponse.ID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedResponse, *actual)
}
