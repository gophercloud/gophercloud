//go:build acceptance || placement || allocationcandidates

package v1

import (
	"context"
	"slices"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/allocationcandidates"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// createRPWithVCPUInventory creates a resource provider and seeds it with
// VCPU inventory, returning the provider UUID. The caller is responsible for
// deferring deletion of the provider.
func createRPWithVCPUInventory(t *testing.T, microversion string) (string, func()) {
	t.Helper()

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = microversion

	rp, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)

	inventories, err := resourceproviders.GetInventories(context.TODO(), client, rp.UUID).Extract()
	th.AssertNoErr(t, err)

	inventories, err = resourceproviders.UpdateInventories(context.TODO(), client, rp.UUID, resourceproviders.UpdateInventoriesOpts{
		ResourceProviderGeneration: inventories.ResourceProviderGeneration,
		Inventories: map[string]resourceproviders.Inventory{
			"VCPU": {
				AllocationRatio: 1.0,
				MaxUnit:         8,
				MinUnit:         1,
				Reserved:        0,
				StepSize:        1,
				Total:           8,
			},
		},
	}).Extract()
	th.AssertNoErr(t, err)

	_, err = resourceproviders.UpdateTraits(context.TODO(), client, rp.UUID, resourceproviders.UpdateTraitsOpts{
		ResourceProviderGeneration: inventories.ResourceProviderGeneration,
		Traits:                     []string{"COMPUTE_NODE"},
	}).Extract()
	th.AssertNoErr(t, err)

	cleanup := func() { DeleteResourceProvider(t, client, rp.UUID) }
	return rp.UUID, cleanup
}

func TestAllocationCandidatesList(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/train")
	clients.RequireAdmin(t)

	microversion := "1.34"
	rpUUID, cleanup := createRPWithVCPUInventory(t, microversion)
	defer cleanup()

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = microversion

	page, err := allocationcandidates.List(client, allocationcandidates.ListOpts{
		Resources: "VCPU:1",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	result, err := allocationcandidates.ExtractAllocationCandidates(page)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, true, len(result.AllocationRequests) > 0)

	// Assert: The provider's summary contains the exact inventory we seeded:
	// VCPU total=8, reserved=0 → capacity=8, used=0.
	summary, present := result.ProviderSummaries[rpUUID]
	th.AssertEquals(t, true, present)
	vcpuSummary, present := summary.Resources["VCPU"]
	th.AssertEquals(t, true, present)
	th.AssertEquals(t, 8, vcpuSummary.Capacity)
	th.AssertEquals(t, 0, vcpuSummary.Used)
	th.AssertEquals(t, true, slices.Contains(*summary.Traits, "COMPUTE_NODE"))
	// It is a root provider: root UUID equals its own UUID, parent is absent.
	th.AssertEquals(t, rpUUID, *summary.RootProviderUUID)
	th.AssertEquals(t, (*string)(nil), summary.ParentProviderUUID)

	// Assert: The allocation request contains the exact resource amount
	// and the unsuffixed group maps to the newly created RP.
	var req allocationcandidates.AllocationRequest
	for _, r := range result.AllocationRequests {
		if _, present := r.Allocations[rpUUID]; present {
			req = r
			break
		}
	}
	th.AssertEquals(t, 1, req.Allocations[rpUUID].Resources["VCPU"])
	th.AssertDeepEquals(t, []string{rpUUID}, (*req.Mappings)[""])

	tools.PrintResource(t, result)
}

func TestAllocationCandidatesListPre129(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/queens")
	clients.RequireAdmin(t)

	microversion := "1.17"
	rpUUID, cleanup := createRPWithVCPUInventory(t, microversion)
	defer cleanup()

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = microversion

	page, err := allocationcandidates.List(client, allocationcandidates.ListOpts{
		Resources: "VCPU:1",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	result, err := allocationcandidates.ExtractAllocationCandidates(page)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, true, len(result.AllocationRequests) > 0)

	summary, present := result.ProviderSummaries[rpUUID]
	th.AssertEquals(t, true, present)
	vcpuSummary, present := summary.Resources["VCPU"]
	th.AssertEquals(t, true, present)
	th.AssertEquals(t, 8, vcpuSummary.Capacity)
	th.AssertEquals(t, 0, vcpuSummary.Used)
	th.AssertEquals(t, true, slices.Contains(*summary.Traits, "COMPUTE_NODE"))
	// Root/parent UUIDs are absent below 1.29.
	th.AssertEquals(t, (*string)(nil), summary.RootProviderUUID)
	th.AssertEquals(t, (*string)(nil), summary.ParentProviderUUID)
	// Mappings are absent below 1.34.
	th.AssertEquals(t, (*map[string][]string)(nil), result.AllocationRequests[0].Mappings)
}

func TestAllocationCandidatesList110(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/pike")
	clients.RequireAdmin(t)

	microversion := "1.10"
	rpUUID, cleanup := createRPWithVCPUInventory(t, microversion)
	defer cleanup()

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = microversion

	page, err := allocationcandidates.List(client, allocationcandidates.ListOpts{
		Resources: "VCPU:1",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	// Pre 1.12 uses an incompatible response format; use the separate function
	// from microversions.go.
	result, err := allocationcandidates.ExtractAllocationCandidates110(page)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, true, len(result.AllocationRequests) > 0)
	th.AssertEquals(t, true, len(result.ProviderSummaries) > 0)

	// Assert: UUID of the created RP present and resource amount correct.
	var foundAlloc allocationcandidates.AllocationRequest110Resource
	for _, req := range result.AllocationRequests {
		for _, alloc := range req.Allocations {
			if alloc.ResourceProvider.UUID == rpUUID {
				foundAlloc = alloc
			}
		}
	}
	th.AssertEquals(t, rpUUID, foundAlloc.ResourceProvider.UUID)
	th.AssertEquals(t, 1, foundAlloc.Resources["VCPU"])

	// Assert: The provider summary contains the expected inventory.
	rpSummary, present := result.ProviderSummaries[rpUUID]
	th.AssertEquals(t, true, present)
	vcpuSummary, present := rpSummary.Resources["VCPU"]
	th.AssertEquals(t, true, present)
	th.AssertEquals(t, 8, vcpuSummary.Capacity)
	th.AssertEquals(t, 0, vcpuSummary.Used)
}

func TestAllocationCandidatesIsEmpty110(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/pike")
	clients.RequireAdmin(t)

	microversion := "1.10"
	_, cleanup := createRPWithVCPUInventory(t, microversion)
	defer cleanup()

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = microversion

	page, err := allocationcandidates.List(client, allocationcandidates.ListOpts{
		Resources: "VCPU:1",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	isEmpty, err := page.IsEmpty()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, false, isEmpty)
}

func TestAllocationCandidatesListEmpty(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/pike")
	clients.RequireAdmin(t)

	microversion := "1.10"
	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = microversion

	page, err := allocationcandidates.List(client, allocationcandidates.ListOpts{
		Resources: "MEMORY_MB:999999",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	result, err := allocationcandidates.ExtractAllocationCandidates(page)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(result.AllocationRequests))
	th.AssertEquals(t, 0, len(result.ProviderSummaries))

	isEmpty, err := page.IsEmpty()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, isEmpty)
}
