package testing

import (
	"context"
	"testing"

	transferRequests "github.com/gophercloud/gophercloud/v2/openstack/dns/v2/transfer/request"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	count := 0
	err := transferRequests.List(client.ServiceClient(fakeServer), nil).EachPage(
		context.TODO(),
		func(_ context.Context, page pagination.Page) (bool, error) {
			count++
			actual, err := transferRequests.ExtractTransferRequests(page)
			th.AssertNoErr(t, err)
			th.CheckDeepEquals(t, ExpectedTransferRequestsSlice, actual)
			return true, nil
		})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestListWithOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	listOpts := transferRequests.ListOpts{
		Status: "ACTIVE",
	}

	allPages, err := transferRequests.List(client.ServiceClient(fakeServer), listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allTransferRequests, err := transferRequests.ExtractTransferRequests(allPages)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 2, len(allTransferRequests))
}

func TestListAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	allPages, err := transferRequests.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allTransferRequests, err := transferRequests.ExtractTransferRequests(allPages)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 2, len(allTransferRequests))
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer)

	actual, err := transferRequests.Get(
		context.TODO(), client.ServiceClient(fakeServer), "a86dba58-0043-4cc6-a1bb-69d5e86f3ca3").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstTransferRequest, actual)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSuccessfully(t, fakeServer)

	createOpts := transferRequests.CreateOpts{
		TargetProjectID: "05d98711-b3a1-4264-a395-f46383671ee6",
		Description:     "This is a first example zone transfer request.",
	}

	actual, err := transferRequests.Create(
		context.TODO(), client.ServiceClient(fakeServer), FirstTransferRequest.ZoneID, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &CreatedTransferRequest, actual)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateSuccessfully(t, fakeServer)

	var description = "Updated Description"
	updateOpts := transferRequests.UpdateOpts{
		Description: description,
	}

	UpdatedTransferRequest := CreatedTransferRequest
	UpdatedTransferRequest.Description = "Updated Description"

	actual, err := transferRequests.Update(
		context.TODO(), client.ServiceClient(fakeServer), UpdatedTransferRequest.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &UpdatedTransferRequest, actual)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSuccessfully(t, fakeServer)

	DeletedZone := CreatedTransferRequest

	err := transferRequests.Delete(context.TODO(), client.ServiceClient(fakeServer), DeletedZone.ID).ExtractErr()
	th.AssertNoErr(t, err)
}
