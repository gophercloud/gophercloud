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
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")
	clients.SkipRelease(t, "stable/ocata")
	clients.SkipRelease(t, "stable/pike")
	clients.SkipRelease(t, "stable/queens")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	name := tools.RandomString("TESTACC-", 8)
	t.Logf("Attempting to create resource provider: %s", name)

	createOpts := resourceproviders.CreateOpts{
		Name: name,
	}

	client.Microversion = "1.20"
	resourceProvider, err := resourceproviders.Create(client, createOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, resourceProvider)
}

func TestResourceProviderUsages(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")
	clients.SkipRelease(t, "stable/ocata")
	clients.SkipRelease(t, "stable/pike")
	clients.SkipRelease(t, "stable/queens")
	clients.RequireAdmin(t)

	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	// first create new resource provider
	name := tools.RandomString("TESTACC-", 8)
	t.Logf("Attempting to create resource provider: %s", name)

	createOpts := resourceproviders.CreateOpts{
		Name: name,
	}

	client.Microversion = "1.20"
	resourceProvider, err := resourceproviders.Create(client, createOpts).Extract()
	th.AssertNoErr(t, err)

	// now get the usages for the newly created resource provider
	usage, err := resourceproviders.GetUsages(client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, usage)
}

func TestResourceProviderInventories(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")
	clients.SkipRelease(t, "stable/ocata")
	clients.SkipRelease(t, "stable/pike")
	clients.SkipRelease(t, "stable/queens")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	// first create new resource provider
	name := tools.RandomString("TESTACC-", 8)
	t.Logf("Attempting to create resource provider: %s", name)

	createOpts := resourceproviders.CreateOpts{
		Name: name,
	}

	client.Microversion = "1.20"
	resourceProvider, err := resourceproviders.Create(client, createOpts).Extract()
	th.AssertNoErr(t, err)

	// now get the inventories for the newly created resource provider
	usage, err := resourceproviders.GetInventories(client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, usage)
}

func TestResourceProviderTraits(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")
	clients.SkipRelease(t, "stable/ocata")
	clients.SkipRelease(t, "stable/pike")
	clients.SkipRelease(t, "stable/queens")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	// first create new resource provider
	name := tools.RandomString("TESTACC-", 8)
	t.Logf("Attempting to create resource provider: %s", name)

	createOpts := resourceproviders.CreateOpts{
		Name: name,
	}

	client.Microversion = "1.20"
	resourceProvider, err := resourceproviders.Create(client, createOpts).Extract()
	th.AssertNoErr(t, err)

	// now get the traits for the newly created resource provider
	usage, err := resourceproviders.GetTraits(client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, usage)
}
