package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/allocations"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListAllocations(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAllocationListSuccessfully(t, fakeServer)

	pages := 0
	err := allocations.List(client.ServiceClient(fakeServer), allocations.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := allocations.ExtractAllocations(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 allocations, got %d", len(actual))
		}
		th.AssertEquals(t, "5344a3e2-978a-444e-990a-cbf47c62ef88", actual[0].UUID)
		th.AssertEquals(t, "eff80f47-75f0-4d41-b1aa-cf07c201adac", actual[1].UUID)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestCreateAllocation(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAllocationCreationSuccessfully(t, fakeServer, SingleAllocationBody)

	actual, err := allocations.Create(context.TODO(), client.ServiceClient(fakeServer), allocations.CreateOpts{
		Name:           "allocation-1",
		ResourceClass:  "baremetal",
		CandidateNodes: []string{"344a3e2-978a-444e-990a-cbf47c62ef88"},
		Traits:         []string{"foo"},
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, Allocation1, *actual)
}

func TestDeleteAllocation(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAllocationDeletionSuccessfully(t, fakeServer)

	res := allocations.Delete(context.TODO(), client.ServiceClient(fakeServer), "344a3e2-978a-444e-990a-cbf47c62ef88")
	th.AssertNoErr(t, res.Err)
}

func TestGetAllocation(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAllocationGetSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := allocations.Get(context.TODO(), c, "344a3e2-978a-444e-990a-cbf47c62ef88").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, Allocation1, *actual)
}
