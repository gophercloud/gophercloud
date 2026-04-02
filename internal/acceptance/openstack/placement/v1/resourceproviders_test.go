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

func TestResourceProviderAggregatesGetSuccess(t *testing.T) {
	// Resource_provider_generation in the aggregates response was introduced in microversion 1.19.
	clients.SkipReleasesBelow(t, "stable/ocata")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	client.Microversion = "1.19"

	aggregates, err := resourceproviders.GetAggregates(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, aggregates.ResourceProviderGeneration != nil)
}

func TestResourceProviderAggregatesGetPreGenerationSuccess(t *testing.T) {
	// Resource provider aggregates operations were introduced in microversion 1.1.
	clients.SkipReleasesBelow(t, "stable/ocata")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	client.Microversion = "1.1"

	aggregates, err := resourceproviders.GetAggregates(context.TODO(), client, resourceProvider.UUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(aggregates.Aggregates))
	th.AssertDeepEquals(t, (*int)(nil), aggregates.ResourceProviderGeneration)
}

func TestResourceProviderAggregatesGetNegative(t *testing.T) {
	// Resource_provider_generation in the aggregates response was introduced in microversion 1.19.
	clients.SkipReleasesBelow(t, "stable/ocata")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.19"

	_, err = resourceproviders.GetAggregates(context.TODO(), client, "00000000-0000-0000-0000-000000000000").Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestResourceProviderAggregatesGetPreGenerationNegative(t *testing.T) {
	// Resource provider aggregates operations were introduced in microversion 1.1.
	clients.SkipReleasesBelow(t, "stable/ocata")
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.1"

	_, err = resourceproviders.GetAggregates(context.TODO(), client, "00000000-0000-0000-0000-000000000000").Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}
