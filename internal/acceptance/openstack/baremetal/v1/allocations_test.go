//go:build acceptance || baremetal || allocations

package v1

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/allocations"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
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
	th.AssertEquals(t, found, true)
}

func TestAllocationsGetDeleteByNode(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBareMetalV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.52"

	node, err := CreateFakeNode(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node)

	resourceClass := "baremetal"
	node, err = nodes.Update(context.TODO(), client, node.UUID, nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:    nodes.ReplaceOp,
			Path:  "/resource_class",
			Value: resourceClass,
		},
	}).Extract()
	th.AssertNoErr(t, err)

	node, err = ChangeProvisionStateAndWait(context.TODO(), client, node, nodes.ProvisionStateOpts{
		Target: nodes.TargetManage,
	}, nodes.Manageable)
	th.AssertNoErr(t, err)

	node, err = ChangeProvisionStateAndWait(context.TODO(), client, node, nodes.ProvisionStateOpts{
		Target: nodes.TargetProvide,
	}, nodes.Available)
	th.AssertNoErr(t, err)

	allocation, err := allocations.Create(context.TODO(), client, allocations.CreateOpts{
		ResourceClass:  resourceClass,
		CandidateNodes: []string{node.UUID},
	}).Extract()
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
	th.AssertNoErr(t, err)

	allocation, err = allocations.GetByNode(context.TODO(), client, node.UUID, nil).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, node.UUID, allocation.NodeUUID)

	allocation, err = allocations.GetByNode(context.TODO(), client, node.UUID, allocations.GetByNodeOpts{
		Fields: []string{"uuid", "node_uuid"},
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, node.UUID, allocation.NodeUUID)

	err = allocations.DeleteByNode(context.TODO(), client, node.UUID).ExtractErr()
	th.AssertNoErr(t, err)
	allocationDeleted = true
}
