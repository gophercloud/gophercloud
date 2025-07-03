package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/ptr"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/pagination"

	fake "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/testhelper"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListFlavors(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorListSuccessfully(t, fakeServer)

	pages := 0
	err := flavors.List(fake.ServiceClient(fakeServer), flavors.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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

func TestListFlavorsEnabled(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	func() {
		testCases := []string{
			"true",
			"false",
			"",
		}

		cases := 0
		fakeServer.Mux.HandleFunc("/v2.0/lbaas/flavors", func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			if err := r.ParseForm(); err != nil {
				t.Errorf("Failed to parse request form %v", err)
			}
			enabled := r.Form.Get("enabled")
			if enabled != testCases[cases] {
				t.Errorf("Expected enabled=%s got %q", testCases[cases], enabled)
			}
			cases++
			fmt.Fprint(w, `{"flavorprofiles":[]}`)
		})
	}()

	var nilBool *bool
	enabled := true
	filters := []*bool{
		&enabled,
		new(bool),
		nilBool,
	}
	for _, filter := range filters {
		allPages, err := flavors.List(fake.ServiceClient(fakeServer), flavors.ListOpts{Enabled: filter}).AllPages(context.TODO())
		th.AssertNoErr(t, err)
		_, err = flavors.ExtractFlavors(allPages)
		th.AssertNoErr(t, err)
	}
}

func TestListAllFlavors(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorListSuccessfully(t, fakeServer)

	allPages, err := flavors.List(fake.ServiceClient(fakeServer), flavors.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := flavors.ExtractFlavors(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, FlavorBasic, actual[0])
	th.CheckDeepEquals(t, FlavorAdvance, actual[1])
}

func TestCreateFlavor(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorCreationSuccessfully(t, fakeServer, SingleFlavorBody)

	actual, err := flavors.Create(context.TODO(), fake.ServiceClient(fakeServer), flavors.CreateOpts{
		Name:            "Basic",
		Description:     "A basic standalone Octavia load balancer.",
		Enabled:         ptr.To(true),
		FlavorProfileId: "9daa2768-74e7-4d13-bf5d-1b8e0dc239e1",
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, FlavorDb, *actual)
}

func TestCreateFlavorDisabled(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorCreationSuccessfullyDisabled(t, fakeServer, SingleFlavorDisabledBody)

	actual, err := flavors.Create(context.TODO(), fake.ServiceClient(fakeServer), flavors.CreateOpts{
		Name:            "Basic",
		Description:     "A basic standalone Octavia load balancer.",
		Enabled:         ptr.To(false),
		FlavorProfileId: "9daa2768-74e7-4d13-bf5d-1b8e0dc239e1",
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, FlavorDisabled, *actual)
}

func TestRequiredCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := flavors.Create(context.TODO(), fake.ServiceClient(fakeServer), flavors.CreateOpts{})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestGetFlavor(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorGetSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := flavors.Get(context.TODO(), client, "5548c807-e6e8-43d7-9ea4-b38d34dd74a0").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, FlavorDb, *actual)
}

func TestDeleteFlavor(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorDeletionSuccessfully(t, fakeServer)

	res := flavors.Delete(context.TODO(), fake.ServiceClient(fakeServer), "5548c807-e6e8-43d7-9ea4-b38d34dd74a0")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateFlavor(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleFlavorUpdateSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := flavors.Update(context.TODO(), client, "5548c807-e6e8-43d7-9ea4-b38d34dd74a0", flavors.UpdateOpts{
		Name:        ptr.To("Basic v2"),
		Description: ptr.To("Rename flavor"),
		Enabled:     ptr.To(true),
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, FlavorUpdated, *actual)
}
