package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/pagination"

	fake "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/testhelper"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestListFlavors(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleFlavorListSuccessfully(t)

	pages := 0
	err := flavors.List(fake.ServiceClient(), flavors.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := flavors.ExtractFlavors(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 flavors, got %d", len(actual))
		}
		th.CheckDeepEquals(t, FlavorBasic, actual[0])
		th.CheckDeepEquals(t, FlavorAdvance, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllFlavors(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleFlavorListSuccessfully(t)

	allPages, err := flavors.List(fake.ServiceClient(), flavors.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := flavors.ExtractFlavors(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, FlavorBasic, actual[0])
	th.CheckDeepEquals(t, FlavorAdvance, actual[1])
}

func TestCreateFlavor(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleFlavorCreationSuccessfully(t, SingleFlavorBody)

	actual, err := flavors.Create(context.TODO(), fake.ServiceClient(), flavors.CreateOpts{
		Name:            "Basic",
		Description:     "A basic standalone Octavia load balancer.",
		Enabled:         true,
		FlavorProfileId: "9daa2768-74e7-4d13-bf5d-1b8e0dc239e1",
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, FlavorDb, *actual)
}

func TestRequiredCreateOpts(t *testing.T) {
	res := flavors.Create(context.TODO(), fake.ServiceClient(), flavors.CreateOpts{})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestGetFlavor(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleFlavorGetSuccessfully(t)

	client := fake.ServiceClient()
	actual, err := flavors.Get(context.TODO(), client, "5548c807-e6e8-43d7-9ea4-b38d34dd74a0").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, FlavorDb, *actual)
}

func TestDeleteFlavor(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleFlavorDeletionSuccessfully(t)

	res := flavors.Delete(context.TODO(), fake.ServiceClient(), "5548c807-e6e8-43d7-9ea4-b38d34dd74a0")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateFlavor(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleFlavorUpdateSuccessfully(t)

	client := fake.ServiceClient()
	actual, err := flavors.Update(context.TODO(), client, "5548c807-e6e8-43d7-9ea4-b38d34dd74a0", flavors.UpdateOpts{
		Name:        "Basic v2",
		Description: "Rename flavor",
		Enabled:     true,
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, FlavorUpdated, *actual)
}
