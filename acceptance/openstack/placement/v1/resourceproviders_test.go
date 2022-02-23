package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/placement/v1/resourceproviders"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestResourceProviderList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	allPages, err := resourceproviders.List(client, resourceproviders.ListOpts{}).AllPages()
	th.AssertNoErr(t, err)

	allResourceProviders, err := resourceproviders.ExtractResourceProviders(allPages)
	th.AssertNoErr(t, err)

	for _, v := range allResourceProviders {
		tools.PrintResource(t, v)
	}
}

func TestResourceProviderCreate(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	resourceProvider, err = CreateResourceProviderWithParent(t, client, resourceProvider.UUID)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)
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
	usage, err := resourceproviders.GetUsages(client, resourceProvider.UUID).Extract()
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
	usage, err := resourceproviders.GetInventories(client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, usage)
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
	usage, err := resourceproviders.GetTraits(client, resourceProvider.UUID).Extract()
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
	usage, err := resourceproviders.GetAllocations(client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, usage)
}
