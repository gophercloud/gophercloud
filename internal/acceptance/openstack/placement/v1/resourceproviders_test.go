//go:build acceptance || placement || resourceproviders

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const InventoryResourceClass = "VCPU"
const MissingInventoryResourceClass = "NO_SUCH_CLASS"
const NonExistentRPUUID = "00000000-0000-0000-0000-000000000000"

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
		Inventories: map[string]resourceproviders.Inventory{
			InventoryResourceClass: {
				AllocationRatio: 1.0,
				MaxUnit:         4,
				MinUnit:         1,
				Reserved:        0,
				StepSize:        1,
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
		Inventories: map[string]resourceproviders.Inventory{
			InventoryResourceClass: {
				AllocationRatio: 1.0,
				MaxUnit:         4,
				MinUnit:         1,
				Reserved:        0,
				StepSize:        1,
				Total:           4,
			},
		},
	}).Extract()
	th.AssertNoErr(t, err)

	_, err = resourceproviders.GetInventory(context.TODO(), client, resourceProvider.UUID, MissingInventoryResourceClass).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
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
		Inventories: map[string]resourceproviders.Inventory{
			InventoryResourceClass: {
				AllocationRatio: 1.0,
				MaxUnit:         4,
				MinUnit:         1,
				Reserved:        0,
				StepSize:        1,
				Total:           4,
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
		Inventory:                  expectedInventory,
	}

	_, err = resourceproviders.UpdateInventory(context.TODO(), client, resourceProvider.UUID, InventoryResourceClass, updateOpts).Extract()
	th.AssertNoErr(t, err)

	updatedInventories, err := resourceproviders.GetInventories(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	actualInventory, ok := updatedInventories.Inventories[InventoryResourceClass]
	th.AssertEquals(t, true, ok)
	th.AssertDeepEquals(t, expectedInventory, actualInventory)
}

func TestResourceProviderUpdateInventoryNotFound(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	updateOpts := resourceproviders.UpdateInventoryOpts{
		ResourceProviderGeneration: 0,
		Inventory: resourceproviders.Inventory{
			AllocationRatio: 1.0,
			MaxUnit:         1,
			MinUnit:         1,
			Reserved:        0,
			StepSize:        1,
			Total:           1,
		},
	}

	_, err = resourceproviders.UpdateInventory(context.TODO(), client, NonExistentRPUUID, InventoryResourceClass, updateOpts).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
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
		Inventories: map[string]resourceproviders.Inventory{
			InventoryResourceClass: {
				AllocationRatio: 1.0,
				MaxUnit:         4,
				MinUnit:         1,
				Reserved:        0,
				StepSize:        1,
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
	th.AssertEquals(t, false, found)
}

func TestResourceProviderDeleteInventoryNotFound(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	err = resourceproviders.DeleteInventory(context.TODO(), client, resourceProvider.UUID, MissingInventoryResourceClass).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestResourceProviderDeleteInventoriesSuccess(t *testing.T) {
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
		Inventories: map[string]resourceproviders.Inventory{
			InventoryResourceClass: {
				AllocationRatio: 1.0,
				MaxUnit:         4,
				MinUnit:         1,
				Reserved:        0,
				StepSize:        1,
				Total:           4,
			},
			"MEMORY_MB": {
				AllocationRatio: 1.0,
				MaxUnit:         1024,
				MinUnit:         1,
				Reserved:        0,
				StepSize:        1,
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
		Inventories:                expectedInventories,
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
		Inventories: map[string]resourceproviders.Inventory{
			InventoryResourceClass: {
				AllocationRatio: 1.0,
				MaxUnit:         1,
				MinUnit:         1,
				Reserved:        0,
				StepSize:        1,
				Total:           1,
			},
		},
	}

	_, err = resourceproviders.UpdateInventories(context.TODO(), client, NonExistentRPUUID, updateOpts).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestResourceProviderTraits(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

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
