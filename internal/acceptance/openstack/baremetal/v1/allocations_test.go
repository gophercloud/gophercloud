//go:build acceptance || baremetal || allocations

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/allocations"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAllocationsCreateDestroy(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.52"

	allocation, err := CreateAllocation(t, client)
	th.AssertNoErr(t, err)
	defer DeleteAllocation(t, client, allocation)

	found := false
	err = allocations.List(client, allocations.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		allocationList, err := allocations.ExtractAllocations(page)
		if err != nil {
			return false, err
		}

		for _, a := range allocationList {
			if a.UUID == allocation.UUID {
				found = true
				return true, nil
			}
		}

		return false, nil
	})
	th.AssertNoErr(t, err)
	th.AssertTrue(t, found)
}

func TestAllocationsPagination(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.52"

	var created []string
	for range 3 {
		allocation, err := CreateAllocation(t, client)
		th.AssertNoErr(t, err)
		t.Cleanup(func() { DeleteAllocation(t, client, allocation) })
		created = append(created, allocation.UUID)
	}

	CheckPagination(t, func(limit int) pagination.Pager {
		return allocations.List(client, allocations.ListOpts{Limit: limit})
	}, func(page pagination.Page) ([]string, error) {
		list, err := allocations.ExtractAllocations(page)
		var ids []string
		for _, a := range list {
			ids = append(ids, a.UUID)
		}
		return ids, err
	}, created)
}
