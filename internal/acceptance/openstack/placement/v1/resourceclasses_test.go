//go:build acceptance || placement || resourceclasses

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

func TestResourceClassCreateByPostSuccess(t *testing.T) {
	// Resource classes were introduced in 1.2
	clients.SkipReleasesBelow(t, "stable/ocata")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.2"

	name := strings.ToUpper(tools.RandomString("CUSTOM_", 8))
	createOpts := resourceclasses.CreateOpts{
		Name: name,
	}

	// Act: Create a resource class using POST
	err = resourceclasses.Create(context.TODO(), client, createOpts).ExtractErr()
	th.AssertNoErr(t, err)

	// Assert: The resource class exists
	rc, err := resourceclasses.Get(context.TODO(), client, name).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, name, rc.Name)
}

func TestResourceClassCreateByPostDuplicate(t *testing.T) {
	// Resource classes were introduced in 1.2
	clients.SkipReleasesBelow(t, "stable/ocata")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.2"

	name := strings.ToUpper(tools.RandomString("CUSTOM_", 8))
	createOpts := resourceclasses.CreateOpts{
		Name: name,
	}

	// Act: Create a resource class using POST
	err = resourceclasses.Create(context.TODO(), client, createOpts).ExtractErr()
	th.AssertNoErr(t, err)

	// Act: Try to create the same resource class again
	err = resourceclasses.Create(context.TODO(), client, createOpts).ExtractErr()
	// Assert: The error is a conflict
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}

func TestResourceClassCreateByUpdateSuccess(t *testing.T) {
	// Creating by Update (PUT) requires microversion 1.7 or later
	clients.SkipReleasesBelow(t, "stable/pike")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.7"

	name := strings.ToUpper(tools.RandomString("CUSTOM_", 8))

	// Act: Create a resource class using PUT (Update)
	err = resourceclasses.Update(context.TODO(), client, name).ExtractErr()
	// No error, with 201 returned
	th.AssertNoErr(t, err)

	// Act: Try to create the same resource class again
	err = resourceclasses.Update(context.TODO(), client, name).ExtractErr()
	// No error, with 204 returned
	th.AssertNoErr(t, err)

	// Assert: The resource class exists
	rc, err := resourceclasses.Get(context.TODO(), client, name).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, name, rc.Name)
}

func TestResourceClassCreateByUpdateNonCustomName(t *testing.T) {
	// Creating by Update (PUT) requires microversion 1.7 or later
	clients.SkipReleasesBelow(t, "stable/pike")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.7"

	name := "CANNOT_CREATE_THIS"

	// Act: Try to create a resource class with a non-custom name using PUT (Update)
	err = resourceclasses.Update(context.TODO(), client, name).ExtractErr()
	// Assert: We get 400
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusBadRequest))
}
