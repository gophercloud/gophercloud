//go:build acceptance || baremetal || allocations

package v1

import (
	"context"
	"strings"
	"testing"
	"time"

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

func TestAllocationsGetDeleteByNode(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.52"

	allocation, err := CreateAllocation(t, client)
	th.AssertNoErr(t, err)

	allocationDeleted := false
	defer func() {
		if !allocationDeleted {
			DeleteAllocation(t, client, allocation)
		}
	}()

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	err = WaitForAllocationState(ctx, client, allocation.UUID, allocations.Active)
	if err != nil {
		if strings.Contains(err.Error(), "no available nodes match the resource class") {
			t.Skipf("skipping node allocation test because no matching nodes are available: %s", err)
		}
		th.AssertNoErr(t, err)
	}

	allocation, err = allocations.Get(context.TODO(), client, allocation.UUID).Extract()
	th.AssertNoErr(t, err)
	if allocation.NodeUUID == "" {
		t.Skipf("allocation %s did not get assigned to a node", allocation.UUID)
	}
	nodeUUID := allocation.NodeUUID

	allocation, err = allocations.GetByNode(context.TODO(), client, nodeUUID, nil).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, nodeUUID, allocation.NodeUUID)

	allocation, err = allocations.GetByNode(context.TODO(), client, nodeUUID, allocations.GetByNodeOpts{
		Fields: []string{"uuid", "node_uuid"},
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, nodeUUID, allocation.NodeUUID)

	err = allocations.DeleteByNode(context.TODO(), client, nodeUUID).ExtractErr()
	th.AssertNoErr(t, err)
	allocationDeleted = true
}
