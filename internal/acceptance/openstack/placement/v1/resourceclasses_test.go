//go:build acceptance || placement || resourceclasses

package v1

import (
	"context"
	"net/http"
	"slices"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceclasses"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestResourceClassesList(t *testing.T) {
	// Resource classes were introduced in 1.2
	clients.SkipReleasesBelow(t, "stable/ocata")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.2"

	allPages, err := resourceclasses.List(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allResourceClasses, err := resourceclasses.ExtractResourceClasses(allPages)
	th.AssertNoErr(t, err)

	// Ensure VCPU is in the list
	th.AssertEquals(t, true, slices.ContainsFunc(allResourceClasses, func(rc resourceclasses.ResourceClass) bool {
		return rc.Name == "VCPU"
	}))
}

func TestResourceClassGetSuccess(t *testing.T) {
	// Resource classes were introduced in 1.2
	clients.SkipReleasesBelow(t, "stable/ocata")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.2"

	// VCPU is a standard resource class
	rc, err := resourceclasses.Get(context.TODO(), client, "VCPU").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "VCPU", rc.Name)
}

func TestResourceClassGetNegative(t *testing.T) {
	// Resource classes were introduced in 1.2
	clients.SkipReleasesBelow(t, "stable/ocata")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.2"

	_, err = resourceclasses.Get(context.TODO(), client, "NON_EXISTENT_RC").Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}
