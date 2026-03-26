//go:build acceptance || placement || traits

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
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/traits"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTraitsList(t *testing.T) {
	// The Traits API requires microversion 1.6 or later
	clients.SkipReleasesBelow(t, "stable/pike")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.6"

	allPages, err := traits.List(client, traits.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allTraits, err := traits.ExtractTraits(allPages)
	th.AssertNoErr(t, err)

	// Ensure COMPUTE_NODE is in the list
	// os-traits never removes traits, so this should always pass
	th.AssertEquals(t, true, slices.Contains(allTraits, "COMPUTE_NODE"))
}

func TestTraitGet(t *testing.T) {
	// The Traits API requires microversion 1.6 or later
	clients.SkipReleasesBelow(t, "stable/pike")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.6"

	// Verify that Get confirms the existence of the COMPUTE_NODE trait
	// os-traits never removes traits, so this should always pass
	err = traits.Get(context.TODO(), client, "COMPUTE_NODE").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestTraitGetNegative(t *testing.T) {
	// The Traits API requires microversion 1.6 or later
	clients.SkipReleasesBelow(t, "stable/pike")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.6"

	// Verify that Get returns an error for a non-existent trait
	err = traits.Get(context.TODO(), client, "CUSTOM_NON_EXISTENT_TRAIT").ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestTraitsListFiltering(t *testing.T) {
	// The Traits API requires microversion 1.6 or later
	clients.SkipReleasesBelow(t, "stable/pike")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.6"

	// os-traits never removes traits, so this should always pass
	listOpts := traits.ListOpts{
		Name: "startswith:HW_",
	}

	allPages, err := traits.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	filteredTraits, err := traits.ExtractTraits(allPages)
	th.AssertNoErr(t, err)

	for _, trait := range filteredTraits {
		th.AssertEquals(t, true, strings.HasPrefix(trait, "HW_"))
	}
}

func TestTraitsCreateSuccess(t *testing.T) {
	// The Traits API requires microversion 1.6 or later
	clients.SkipReleasesBelow(t, "stable/pike")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.6"

	traitName := strings.ToUpper(tools.RandomString("CUSTOM_", 8))
	createOpts := traits.CreateOpts{}

	err = traits.Create(context.TODO(), client, traitName, createOpts).ExtractErr()
	th.AssertNoErr(t, err)

	// Assert that the trait now exists
	err = traits.Get(context.TODO(), client, traitName).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestTraitsCreateDuplicate(t *testing.T) {
	// The Traits API requires microversion 1.6 or later
	clients.SkipReleasesBelow(t, "stable/pike")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.6"

	traitName := strings.ToUpper(tools.RandomString("CUSTOM_", 8))
	createOpts := traits.CreateOpts{}

	// Create the trait for the first time
	err = traits.Create(context.TODO(), client, traitName, createOpts).ExtractErr()
	th.AssertNoErr(t, err)

	// Creating the same trait again results in 204 (no error)
	err = traits.Create(context.TODO(), client, traitName, createOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

// Test of creating a trait name that cannot be created in an API
func TestTraitsCreateInvalidName(t *testing.T) {
	// The Traits API requires microversion 1.6 or later
	clients.SkipReleasesBelow(t, "stable/pike")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.6"

	traitName := "HW_WE_CANNOT_CREATE_THIS_TRAIT"
	createOpts := traits.CreateOpts{}

	err = traits.Create(context.TODO(), client, traitName, createOpts).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusBadRequest))
}
