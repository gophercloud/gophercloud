package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/trunks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, CreateResponse)
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
	_, err := trunks.Create(fake.ServiceClient(), options).Extract()
	if err == nil {
		t.Fatalf("Failed to detect missing parent PortID field")
	}
	options.PortID = "c373d2fa-3d3b-4492-924c-aff54dea19b6"
	n, err := trunks.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "ACTIVE")
	expectedTrunks, err := ExpectedTrunkSlice()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedTrunks[1], n)
}

func TestCreateNoSubports(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/trunks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateNoSubportsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, CreateNoSubportsResponse)
	})

	iTrue := true
	options := trunks.CreateOpts{
		Name:         "gophertrunk",
		Description:  "Trunk created by gophercloud",
		AdminStateUp: &iTrue,
		PortID:       "c373d2fa-3d3b-4492-924c-aff54dea19b6",
	}
	n, err := trunks.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "ACTIVE")
	th.AssertEquals(t, 0, len(n.Subports))
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/trunks/f6a9718c-5a64-43e3-944f-4deccad8e78c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := trunks.Delete(fake.ServiceClient(), "f6a9718c-5a64-43e3-944f-4deccad8e78c")
	th.AssertNoErr(t, res.Err)
}
