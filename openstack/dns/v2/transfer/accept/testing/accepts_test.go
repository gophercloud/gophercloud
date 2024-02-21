package testing

import (
	"context"
	"testing"

	transferAccepts "github.com/gophercloud/gophercloud/v2/openstack/dns/v2/transfer/accept"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	count := 0
	err := transferAccepts.List(client.ServiceClient(), nil).EachPage(
		context.TODO(),
		func(_ context.Context, page pagination.Page) (bool, error) {
			count++
			actual, err := transferAccepts.ExtractTransferAccepts(page)
			th.AssertNoErr(t, err)
			th.CheckDeepEquals(t, ExpectedTransferAcceptSlice, actual)
			return true, nil
		})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestListWithOpts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleFilteredListSuccessfully(t)

	listOpts := transferAccepts.ListOpts{
		Status: "ACTIVE",
	}

	allPages, err := transferAccepts.List(client.ServiceClient(), listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allTransferAccepts, err := transferAccepts.ExtractTransferAccepts(allPages)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, len(allTransferAccepts))
}

func TestListAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	allPages, err := transferAccepts.List(client.ServiceClient(), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allTransferAccepts, err := transferAccepts.ExtractTransferAccepts(allPages)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 2, len(allTransferAccepts))
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)

	actual, err := transferAccepts.Get(
		context.TODO(), client.ServiceClient(), "92236f39-0fad-4f8f-bf25-fbdf027de34d").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstTransferAccept, actual)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t)

	createOpts := transferAccepts.CreateOpts{
		Key:                   "M2KA0Y20",
		ZoneTransferRequestID: "fc46bb1f-bdf0-4e67-96e0-f8c04f26261c",
	}

	actual, err := transferAccepts.Create(
		context.TODO(), client.ServiceClient(), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &CreatedTransferAccept, actual)
}
