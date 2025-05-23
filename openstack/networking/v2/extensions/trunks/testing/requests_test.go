package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/trunks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreateResponse)
	})

	iTrue := true
	options := trunks.CreateOpts{
		Name:         "gophertrunk",
		Description:  "Trunk created by gophercloud",
		AdminStateUp: &iTrue,
		Subports: []trunks.Subport{
			{
				SegmentationID:   1,
				SegmentationType: "vlan",
				PortID:           "28e452d7-4f8a-4be4-b1e6-7f3db4c0430b",
			},
			{
				SegmentationID:   2,
				SegmentationType: "vlan",
				PortID:           "4c8b2bff-9824-4d4c-9b60-b3f6621b2bab",
			},
		},
	}
	_, err := trunks.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	if err == nil {
		t.Fatalf("Failed to detect missing parent PortID field")
	}
	options.PortID = "c373d2fa-3d3b-4492-924c-aff54dea19b6"
	n, err := trunks.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "ACTIVE")
	expectedTrunks, err := ExpectedTrunkSlice()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedTrunks[1], n)
}

func TestCreateNoSubports(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/trunks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateNoSubportsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreateNoSubportsResponse)
	})

	iTrue := true
	options := trunks.CreateOpts{
		Name:         "gophertrunk",
		Description:  "Trunk created by gophercloud",
		AdminStateUp: &iTrue,
		PortID:       "c373d2fa-3d3b-4492-924c-aff54dea19b6",
	}
	n, err := trunks.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "ACTIVE")
	th.AssertEquals(t, 0, len(n.Subports))
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/trunks/f6a9718c-5a64-43e3-944f-4deccad8e78c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := trunks.Delete(context.TODO(), fake.ServiceClient(fakeServer), "f6a9718c-5a64-43e3-944f-4deccad8e78c")
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/trunks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponse)
	})

	client := fake.ServiceClient(fakeServer)
	count := 0

	err := trunks.List(client, trunks.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := trunks.ExtractTrunks(page)
		if err != nil {
			t.Errorf("Failed to extract trunks: %v", err)
			return false, err
		}

		expected, err := ExpectedTrunkSlice()
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/trunks/f6a9718c-5a64-43e3-944f-4deccad8e78c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetResponse)
	})

	n, err := trunks.Get(context.TODO(), fake.ServiceClient(fakeServer), "f6a9718c-5a64-43e3-944f-4deccad8e78c").Extract()
	th.AssertNoErr(t, err)
	expectedTrunks, err := ExpectedTrunkSlice()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expectedTrunks[1], n)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/trunks/f6a9718c-5a64-43e3-944f-4deccad8e78c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateResponse)
	})

	iFalse := false
	name := "updated_gophertrunk"
	description := "gophertrunk updated by gophercloud"
	options := trunks.UpdateOpts{
		Name:         &name,
		AdminStateUp: &iFalse,
		Description:  &description,
	}
	n, err := trunks.Update(context.TODO(), fake.ServiceClient(fakeServer), "f6a9718c-5a64-43e3-944f-4deccad8e78c", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Name, name)
	th.AssertEquals(t, n.AdminStateUp, iFalse)
	th.AssertEquals(t, n.Description, description)
}

func TestGetSubports(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/trunks/f6a9718c-5a64-43e3-944f-4deccad8e78c/get_subports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListSubportsResponse)
	})

	client := fake.ServiceClient(fakeServer)

	subports, err := trunks.GetSubports(context.TODO(), client, "f6a9718c-5a64-43e3-944f-4deccad8e78c").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedSubports, subports)
}

func TestMissingFields(t *testing.T) {
	iTrue := true
	opts := trunks.CreateOpts{
		Name:         "gophertrunk",
		PortID:       "c373d2fa-3d3b-4492-924c-aff54dea19b6",
		Description:  "Trunk created by gophercloud",
		AdminStateUp: &iTrue,
		Subports: []trunks.Subport{
			{
				SegmentationID:   1,
				SegmentationType: "vlan",
				PortID:           "28e452d7-4f8a-4be4-b1e6-7f3db4c0430b",
			},
			{
				SegmentationID:   2,
				SegmentationType: "vlan",
				PortID:           "4c8b2bff-9824-4d4c-9b60-b3f6621b2bab",
			},
			{
				PortID: "4c8b2bff-9824-4d4c-9b60-b3f6621b2bab",
			},
		},
	}

	_, err := opts.ToTrunkCreateMap()
	if err == nil {
		t.Fatalf("Failed to detect missing subport fields")
	}
}

func TestAddSubports(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/trunks/f6a9718c-5a64-43e3-944f-4deccad8e78c/add_subports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AddSubportsRequest)
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, AddSubportsResponse)
	})

	client := fake.ServiceClient(fakeServer)

	opts := trunks.AddSubportsOpts{
		Subports: ExpectedSubports,
	}

	trunk, err := trunks.AddSubports(context.TODO(), client, "f6a9718c-5a64-43e3-944f-4deccad8e78c", opts).Extract()
	th.AssertNoErr(t, err)
	expectedTrunk, err := ExpectedSubportsAddedTrunk()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expectedTrunk, trunk)
}

func TestRemoveSubports(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/trunks/f6a9718c-5a64-43e3-944f-4deccad8e78c/remove_subports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, RemoveSubportsRequest)
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, RemoveSubportsResponse)
	})

	client := fake.ServiceClient(fakeServer)

	opts := trunks.RemoveSubportsOpts{
		Subports: []trunks.RemoveSubport{
			{PortID: "28e452d7-4f8a-4be4-b1e6-7f3db4c0430b"},
			{PortID: "4c8b2bff-9824-4d4c-9b60-b3f6621b2bab"},
		},
	}
	trunk, err := trunks.RemoveSubports(context.TODO(), client, "f6a9718c-5a64-43e3-944f-4deccad8e78c", opts).Extract()

	th.AssertNoErr(t, err)
	expectedTrunk, err := ExpectedSubportsRemovedTrunk()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expectedTrunk, trunk)
}
