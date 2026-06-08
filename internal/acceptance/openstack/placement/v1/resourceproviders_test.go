//go:build acceptance || placement || resourceproviders

package v1

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/internal/ptr"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/traits"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const InventoryResourceClass = "VCPU"
const MissingInventoryResourceClass = "NO_SUCH_CLASS"

func TestResourceProviderList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	allPages, err := resourceproviders.List(client, resourceproviders.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allResourceProviders, err := resourceproviders.ExtractResourceProviders(allPages)
	th.AssertNoErr(t, err)

	for _, v := range allResourceProviders {
		tools.PrintResource(t, v)
	}
}

func TestResourceProviderList139(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.39"

	// Arrange: Create a resource provider, traits, and aggregates.
	// Assign them to the created RP.
	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	trait1 := strings.ToUpper(tools.RandomString("CUSTOM_", 8))
	trait2 := strings.ToUpper(tools.RandomString("CUSTOM_", 8))
	err = traits.Create(context.TODO(), client, trait1).ExtractErr()
	th.AssertNoErr(t, err)
	defer traits.Delete(context.TODO(), client, trait1)
	err = traits.Create(context.TODO(), client, trait2).ExtractErr()
	th.AssertNoErr(t, err)
	defer traits.Delete(context.TODO(), client, trait2)

	currentTraits, err := resourceproviders.GetTraits(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	_, err = resourceproviders.UpdateTraits(context.TODO(), client, resourceProvider.UUID, resourceproviders.UpdateTraitsOpts{
		ResourceProviderGeneration: currentTraits.ResourceProviderGeneration,
		Traits:                     []string{trait1, trait2},
	}).Extract()
	th.AssertNoErr(t, err)

	currentAggregates, err := resourceproviders.GetAggregates(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	aggregate1 := tools.RandomUUID()
	aggregate2 := tools.RandomUUID()
	aggregate3 := tools.RandomUUID()
	_, err = resourceproviders.UpdateAggregates(context.TODO(), client, resourceProvider.UUID, resourceproviders.UpdateAggregatesOpts{
		ResourceProviderGeneration: currentAggregates.ResourceProviderGeneration,
		Aggregates:                 []string{aggregate1, aggregate2},
	}).Extract()
	th.AssertNoErr(t, err)

	listOpts := resourceproviders.ListOpts139{
		// Repeating member_of means AND: provider must be in aggregate1 and in any of (aggregate2, aggregate3).
		// We'll expect list to return our provider.
		MemberOf: []string{aggregate1, "in:" + aggregate2 + "," + aggregate3},
		Required: []string{trait1, trait2},
	}

	// Act: List resource providers with the above traits and aggregates as filters.
	allPages, err := resourceproviders.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	// Assert: Our resource provider is in the results and has the traits and aggregates we set.
	allResourceProviders, err := resourceproviders.ExtractResourceProviders(allPages)
	th.AssertNoErr(t, err)
	th.AssertTrue(t, len(allResourceProviders) > 0)

	found := false
	for _, rp := range allResourceProviders {
		if rp.UUID == resourceProvider.UUID {
			found = true
			break
		}
	}
	th.AssertTrue(t, found)
}

func TestResourceProvider(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")
	clients.SkipRelease(t, "stable/ocata")
	clients.SkipRelease(t, "stable/pike")
	clients.SkipRelease(t, "stable/queens")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	resourceProvider2, err := CreateResourceProviderWithParent(t, client, resourceProvider.UUID)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider2.UUID)

	newName := tools.RandomString("TESTACC-", 8)
	updateOpts := resourceproviders.UpdateOpts{
		Name: &newName,
	}
	resourceProviderUpdate, err := resourceproviders.Update(context.TODO(), client, resourceProvider2.UUID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newName, resourceProviderUpdate.Name)

	resourceProviderGet, err := resourceproviders.Get(context.TODO(), client, resourceProvider2.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newName, resourceProviderGet.Name)

}

func TestResourceProviderUsages(t *testing.T) {
	clients.RequireAdmin(t)

	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	// first create new resource provider
	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	// now get the usages for the newly created resource provider
	usage, err := resourceproviders.GetUsages(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, usage)
}

func TestResourceProviderInventories(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	// first create new resource provider
	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	// now get the inventories for the newly created resource provider
	usage, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, usage)
}

func TestResourceProviderInventory(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	inventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	seededInventories, err := resourceproviders.UpdateInventories(context.TODO(), client, resourceProvider.UUID, resourceproviders.UpdateInventoriesOpts{
		ResourceProviderGeneration: inventories.ResourceProviderGeneration,
		Inventories: map[string]resourceproviders.InventoryUpdateBase{
			InventoryResourceClass: {
				AllocationRatio: ptr.To(float32(1.0)),
				MaxUnit:         ptr.To(4),
				MinUnit:         ptr.To(1),
				Reserved:        ptr.To(0),
				StepSize:        ptr.To(1),
				Total:           4,
			},
		},
	}).Extract()
	th.AssertNoErr(t, err)

	inventory, err := resourceproviders.GetInventory(context.TODO(), client, resourceProvider.UUID, InventoryResourceClass).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, seededInventories.ResourceProviderGeneration, inventory.ResourceProviderGeneration)
	th.AssertDeepEquals(t, seededInventories.Inventories[InventoryResourceClass], inventory.Inventory)
}

func TestResourceProviderInventoryNotFound(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	inventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	_, err = resourceproviders.UpdateInventories(context.TODO(), client, resourceProvider.UUID, resourceproviders.UpdateInventoriesOpts{
		ResourceProviderGeneration: inventories.ResourceProviderGeneration,
		Inventories: map[string]resourceproviders.InventoryUpdateBase{
			InventoryResourceClass: {
				AllocationRatio: ptr.To(float32(1.0)),
				MaxUnit:         ptr.To(4),
				MinUnit:         ptr.To(1),
				Reserved:        ptr.To(0),
				StepSize:        ptr.To(1),
				Total:           4,
			},
		},
	}).Extract()
	th.AssertNoErr(t, err)

	_, err = resourceproviders.GetInventory(context.TODO(), client, resourceProvider.UUID, MissingInventoryResourceClass).Extract()
	th.AssertTrue(t, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestResourceProviderUpdateInventory(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	// Arrange: Get the current inventory to retrieve the generation
	inventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	// Arrange: The resource class on this provider must exist first
	seedOpts := resourceproviders.UpdateInventoriesOpts{
		ResourceProviderGeneration: inventories.ResourceProviderGeneration,
		Inventories: map[string]resourceproviders.InventoryUpdateBase{
			InventoryResourceClass: {
				AllocationRatio: ptr.To(float32(1.0)),
				MaxUnit:         ptr.To(4),
				MinUnit:         ptr.To(1),
				// Skipping Reserved on purpose
				StepSize: ptr.To(1),
				Total:    4,
			},
		},
	}

	seededInventories, err := resourceproviders.UpdateInventories(context.TODO(), client, resourceProvider.UUID, seedOpts).Extract()
	th.AssertNoErr(t, err)

	expectedInventory := resourceproviders.Inventory{
		AllocationRatio: 1.0,
		MaxUnit:         8,
		MinUnit:         1,
		Reserved:        0,
		StepSize:        1,
		Total:           8,
	}

	updateOpts := resourceproviders.UpdateInventoryOpts{
		ResourceProviderGeneration: seededInventories.ResourceProviderGeneration,
		InventoryUpdateBase: resourceproviders.InventoryUpdateBase{
			AllocationRatio: ptr.To(expectedInventory.AllocationRatio),
			MaxUnit:         ptr.To(expectedInventory.MaxUnit),
			MinUnit:         ptr.To(expectedInventory.MinUnit),
			StepSize:        ptr.To(expectedInventory.StepSize),
			Total:           expectedInventory.Total,
		},
	}

	_, err = resourceproviders.UpdateInventory(context.TODO(), client, resourceProvider.UUID, InventoryResourceClass, updateOpts).Extract()
	th.AssertNoErr(t, err)

	updatedInventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	actualInventory, ok := updatedInventories.Inventories[InventoryResourceClass]
	th.AssertTrue(t, ok)
	th.AssertDeepEquals(t, expectedInventory, actualInventory)
}

func TestResourceProviderUpdateInventoryNotFound(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	updateOpts := resourceproviders.UpdateInventoryOpts{
		ResourceProviderGeneration: 0,
		InventoryUpdateBase: resourceproviders.InventoryUpdateBase{
			AllocationRatio: ptr.To(float32(1.0)),
			MaxUnit:         ptr.To(1),
			MinUnit:         ptr.To(1),
			Reserved:        ptr.To(0),
			StepSize:        ptr.To(1),
			Total:           1,
		},
	}

	_, err = resourceproviders.UpdateInventory(context.TODO(), client, tools.RandomUUID(), InventoryResourceClass, updateOpts).Extract()
	th.AssertTrue(t, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestResourceProviderDeleteInventorySuccess(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	// Arrange: Get the current inventory to retrieve the generation
	inventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	_, err = resourceproviders.UpdateInventories(context.TODO(), client, resourceProvider.UUID, resourceproviders.UpdateInventoriesOpts{
		ResourceProviderGeneration: inventories.ResourceProviderGeneration,
		Inventories: map[string]resourceproviders.InventoryUpdateBase{
			InventoryResourceClass: {
				AllocationRatio: ptr.To(float32(1.0)),
				MaxUnit:         ptr.To(4),
				MinUnit:         ptr.To(1),
				Reserved:        ptr.To(0),
				StepSize:        ptr.To(1),
				Total:           4,
			},
		},
	}).Extract()
	th.AssertNoErr(t, err)

	err = resourceproviders.DeleteInventory(context.TODO(), client, resourceProvider.UUID, InventoryResourceClass).ExtractErr()
	th.AssertNoErr(t, err)

	// Assert: The inventory should no longer be found
	updatedInventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	_, found := updatedInventories.Inventories[InventoryResourceClass]
	th.AssertFalse(t, found)
}

func TestResourceProviderDeleteInventoryNotFound(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	err = resourceproviders.DeleteInventory(context.TODO(), client, resourceProvider.UUID, MissingInventoryResourceClass).ExtractErr()
	th.AssertTrue(t, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestResourceProviderDeleteInventoriesSuccess(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.20"

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	inventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	_, err = resourceproviders.UpdateInventories(context.TODO(), client, resourceProvider.UUID, resourceproviders.UpdateInventoriesOpts{
		ResourceProviderGeneration: inventories.ResourceProviderGeneration,
		Inventories: map[string]resourceproviders.InventoryUpdateBase{
			InventoryResourceClass: {
				AllocationRatio: ptr.To(float32(1.0)),
				MaxUnit:         ptr.To(4),
				MinUnit:         ptr.To(1),
				Reserved:        ptr.To(0),
				StepSize:        ptr.To(1),
				Total:           4,
			},
			"MEMORY_MB": {
				AllocationRatio: ptr.To(float32(1.0)),
				MaxUnit:         ptr.To(1024),
				MinUnit:         ptr.To(1),
				Reserved:        ptr.To(0),
				StepSize:        ptr.To(1),
				Total:           1024,
			},
		},
	}).Extract()
	th.AssertNoErr(t, err)

	seededInventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, len(seededInventories.Inventories))

	err = resourceproviders.DeleteInventories(context.TODO(), client, resourceProvider.UUID).ExtractErr()
	th.AssertNoErr(t, err)

	updatedInventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(updatedInventories.Inventories))
}

func TestResourceProviderUpdateInventories(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	// Arrange: Get the current inventory to retrieve the generation
	inventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	expectedInventories := map[string]resourceproviders.Inventory{
		"DISK_GB": {
			AllocationRatio: 1.0,
			MaxUnit:         100,
			MinUnit:         1,
			Reserved:        0,
			StepSize:        1,
			Total:           100,
		},
	}

	updateOpts := resourceproviders.UpdateInventoriesOpts{
		ResourceProviderGeneration: inventories.ResourceProviderGeneration,
		Inventories: map[string]resourceproviders.InventoryUpdateBase{
			"DISK_GB": {
				AllocationRatio: ptr.To(float32(1.0)),
				MaxUnit:         ptr.To(100),
				MinUnit:         ptr.To(1),
				Reserved:        ptr.To(0),
				StepSize:        ptr.To(1),
				Total:           100,
			},
		},
	}

	_, err = resourceproviders.UpdateInventories(context.TODO(), client, resourceProvider.UUID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	updatedInventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, expectedInventories, updatedInventories.Inventories)
}

func TestResourceProviderUpdateInventoriesNotFound(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	updateOpts := resourceproviders.UpdateInventoriesOpts{
		ResourceProviderGeneration: 0,
		Inventories: map[string]resourceproviders.InventoryUpdateBase{
			InventoryResourceClass: {
				AllocationRatio: ptr.To(float32(1.0)),
				MaxUnit:         ptr.To(1),
				MinUnit:         ptr.To(1),
				Reserved:        ptr.To(0),
				StepSize:        ptr.To(1),
				Total:           1,
			},
		},
	}

	_, err = resourceproviders.UpdateInventories(context.TODO(), client, tools.RandomUUID(), updateOpts).Extract()
	th.AssertTrue(t, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestResourceProviderTraits(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.20"

	// first create new resource provider
	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	// now get the traits for the newly created resource provider
	usage, err := resourceproviders.GetTraits(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, usage)
}

func TestResourceProviderAllocations(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	// first create new resource provider
	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	// now get the allocations for the newly created resource provider
	usage, err := resourceproviders.GetAllocations(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, usage)
}

func TestResourceProviderAggregates(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/ocata")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	// first create new resource provider
	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	// now get the aggregates for same
	client.Microversion = "1.19"
	aggregates, err := resourceproviders.GetAggregates(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertTrue(t, aggregates.ResourceProviderGeneration != nil)

	// ensure that we handle older microversions where generation is missing
	client.Microversion = "1.1"
	aggregates, err = resourceproviders.GetAggregates(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(aggregates.Aggregates))
	th.AssertDeepEquals(t, (*int)(nil), aggregates.ResourceProviderGeneration)
}

func TestResourceProviderAggregatesNotFound(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/ocata")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.19"
	_, err = resourceproviders.GetAggregates(context.TODO(), client, tools.RandomUUID()).Extract()
	th.AssertTrue(t, gophercloud.ResponseCodeIs(err, http.StatusNotFound))

	client.Microversion = "1.1"
	_, err = resourceproviders.GetAggregates(context.TODO(), client, tools.RandomUUID()).Extract()
	th.AssertTrue(t, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestResourceProviderUpdateAggregates(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/ocata")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	client.Microversion = "1.19"

	before, err := resourceproviders.GetAggregates(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertTrue(t, before.ResourceProviderGeneration != nil)

	updateOpts := resourceproviders.UpdateAggregatesOpts{
		ResourceProviderGeneration: before.ResourceProviderGeneration,
		Aggregates: []string{
			tools.RandomUUID(),
			tools.RandomUUID(),
		},
	}

	_, err = resourceproviders.UpdateAggregates(context.TODO(), client, resourceProvider.UUID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	after, err := resourceproviders.GetAggregates(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(updateOpts.Aggregates), len(after.Aggregates))

	for _, aggregate := range updateOpts.Aggregates {
		th.AssertTrue(t, slices.Contains(after.Aggregates, aggregate))
	}
}

func TestResourceProviderUpdateAggregateMismatch(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/ocata")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	client.Microversion = "1.19"

	current, err := resourceproviders.GetAggregates(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertTrue(t, current.ResourceProviderGeneration != nil)
	wrongGeneration := *current.ResourceProviderGeneration + 100

	updateOpts := resourceproviders.UpdateAggregatesOpts{
		ResourceProviderGeneration: &wrongGeneration,
		Aggregates: []string{
			tools.RandomUUID(),
		},
	}

	_, err = resourceproviders.UpdateAggregates(context.TODO(), client, resourceProvider.UUID, updateOpts).Extract()
	th.AssertTrue(t, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}

func TestResourceProviderUpdateAggregatesPreGeneration(t *testing.T) {
	clients.SkipReleasesBelow(t, "stable/ocata")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	client.Microversion = "1.1"

	updateOpts := resourceproviders.UpdateAggregatesOpts{
		Aggregates: []string{
			tools.RandomUUID(),
			tools.RandomUUID(),
		},
	}

	_, err = resourceproviders.UpdateAggregates(context.TODO(), client, resourceProvider.UUID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	after, err := resourceproviders.GetAggregates(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(updateOpts.Aggregates), len(after.Aggregates))

	for _, aggregate := range updateOpts.Aggregates {
		th.AssertTrue(t, slices.Contains(after.Aggregates, aggregate))
	}
}

func TestResourceProviderUpdateAggregatesPreGenerationWithGenerationSuccess(t *testing.T) {
	// Before microversion 1.19, ResourceProviderGeneration in opts is silently stripped from
	// the request body, so the operation must succeed even when the caller supplies it.
	clients.SkipReleasesBelow(t, "stable/ocata")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	client.Microversion = "1.1"

	gen := 1
	updateOpts := resourceproviders.UpdateAggregatesOpts{
		ResourceProviderGeneration: &gen,
		Aggregates: []string{
			tools.RandomUUID(),
			tools.RandomUUID(),
		},
	}

	_, err = resourceproviders.UpdateAggregates(context.TODO(), client, resourceProvider.UUID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	after, err := resourceproviders.GetAggregates(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(updateOpts.Aggregates), len(after.Aggregates))

	for _, aggregate := range updateOpts.Aggregates {
		th.AssertTrue(t, slices.Contains(after.Aggregates, aggregate))
	}
}

func TestResourceProviderParentDetach(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	// Arrange: Create a parent resource provider
	parentProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, parentProvider.UUID)

	// Arrange: Create a child resource provider with that parent
	childProvider, err := CreateResourceProviderWithParent(t, client, parentProvider.UUID)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, childProvider.UUID)

	// Sanity check: Verify that the child provider has the correct parent before the update
	th.AssertEquals(t, parentProvider.UUID, childProvider.ParentProviderUUID)

	// Act: Update the child resource provider to remove the parent (transform to root)
	client.Microversion = "1.37"
	empty := ""
	updateOpts := resourceproviders.UpdateOpts{
		Name:               &childProvider.Name,
		ParentProviderUUID: &empty,
	}
	updatedChild, err := resourceproviders.Update(context.TODO(), client, childProvider.UUID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	// Assert: Verify that the ParentProviderUUID is now null (empty string in Gophercloud result struct)
	th.AssertEquals(t, "", updatedChild.ParentProviderUUID)

	// Assert: Double check with a Get request
	childGet, err := resourceproviders.Get(context.TODO(), client, childProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "", childGet.ParentProviderUUID)
}
