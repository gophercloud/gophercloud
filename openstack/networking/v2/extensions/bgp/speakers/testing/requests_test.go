package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/speakers"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/bgp-speakers",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, ListBGPSpeakerResult)
		})
	count := 0

	speakers.List(fake.ServiceClient()).EachPage(
		func(page pagination.Page) (bool, error) {
			count++
			actual, err := speakers.ExtractBGPSpeakers(page)

			if err != nil {
				t.Errorf("Failed to extract BGP speakers: %v", err)
				return false, nil
			}
			expected := []speakers.BGPSpeaker{BGPSpeaker1}
			th.CheckDeepEquals(t, expected, actual)
			return true, nil
		})
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpSpeakerID := "ab01ade1-ae62-43c9-8a1f-3c24225b96d8"
	th.Mux.HandleFunc("/v2.0/bgp-speakers/"+bgpSpeakerID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetBGPSpeakerResult)
	})

	s, err := speakers.Get(fake.ServiceClient(), bgpSpeakerID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, *s, BGPSpeaker1)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/bgp-speakers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, CreateResponse)
	})

	opts := speakers.CreateOpts{
		IPVersion:                     6,
		AdvertiseFloatingIPHostRoutes: false,
		AdvertiseTenantNetworks:       true,
		Name:                          "gophercloud-testing-bgp-speaker",
		LocalAS:                       "2000",
		Networks:                      []string{},
	}
	r, err := speakers.Create(fake.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, r.Name, opts.Name)
	th.AssertEquals(t, r.LocalAS, 2000)
	th.AssertEquals(t, len(r.Networks), 0)
	th.AssertEquals(t, r.IPVersion, opts.IPVersion)
	th.AssertEquals(t, r.AdvertiseFloatingIPHostRoutes, opts.AdvertiseFloatingIPHostRoutes)
	th.AssertEquals(t, r.AdvertiseTenantNetworks, opts.AdvertiseTenantNetworks)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpSpeakerID := "ab01ade1-ae62-43c9-8a1f-3c24225b96d8"
	th.Mux.HandleFunc("/v2.0/bgp-speakers/"+bgpSpeakerID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	err := speakers.Delete(fake.ServiceClient(), bgpSpeakerID).ExtractErr()
	th.AssertNoErr(t, err)
}
