package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetransfers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateTransfer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateTransfer(t)

	actual, err := sharetransfers.Create(context.TODO(), client.ServiceClient(), TransferRequest).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, TransferResponse, *actual)
}

func TestAcceptTransfer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAcceptTransfer(t)

	err := sharetransfers.Accept(context.TODO(), client.ServiceClient(), TransferResponse.ID, AcceptRequest).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDeleteTransfer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteTransfer(t)

	err := sharetransfers.Delete(context.TODO(), client.ServiceClient(), TransferResponse.ID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListTransfers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTransfers(t)

	expectedResponse := TransferListResponse
	expectedResponse[0].AuthKey = ""

	count := 0
	err := sharetransfers.List(client.ServiceClient(), &sharetransfers.ListOpts{AllTenants: true}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := sharetransfers.ExtractTransfers(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, expectedResponse, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListTransfersDetail(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTransfersDetail(t)

	expectedResponse := TransferListResponse
	expectedResponse[0].AuthKey = ""

	count := 0
	err := sharetransfers.ListDetail(client.ServiceClient(), &sharetransfers.ListOpts{AllTenants: true}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := sharetransfers.ExtractTransfers(page)
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

	allPages, err := sharetransfers.List(client.ServiceClient(), &sharetransfers.ListOpts{AllTenants: true}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := sharetransfers.ExtractTransfers(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedResponse, actual)
}

func TestGetTransfer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetTransfer(t)

	expectedResponse := TransferResponse
	expectedResponse.AuthKey = ""

	actual, err := sharetransfers.Get(context.TODO(), client.ServiceClient(), TransferResponse.ID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedResponse, *actual)
}
