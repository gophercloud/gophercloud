package testing

import (
	"fmt"
	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/speaker"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"io"
	"net/http"
	"testing"
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

	speaker.List(fake.ServiceClient()).EachPage(
		func(page pagination.Page) (bool, error) {
			count++
			actual, err := speaker.ExtractBGPSpeakers(page)

			if err != nil {
				t.Errorf("Failed to extract BGP Speaker: %v", err)
				return false, nil
			}
			expected := []speaker.BGPSpeaker{BGPSpeaker1}
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

	s, err := speaker.Get(fake.ServiceClient(), bgpSpeakerID).Extract()
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

	name := "gophercloud-testing-bgp-speaker"
	localas := "2000"
	networks := []string{}
	m := make(map[string]string)
	m["IPVersion"] = "6"
	m["AdvertiseFloatingIPHostRoutes"] = "false"
	opts := speaker.BuildCreateOpts(name, localas, networks, m)

	r, err := speaker.Create(fake.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, r.Name, opts.Name)
	th.AssertEquals(t, r.LocalAS, 2000)
	th.AssertEquals(t, len(r.Networks), 0)
	th.AssertEquals(t, r.IPVersion, opts.IPVersion)
	th.AssertEquals(t, r.AdvertiseFloatingIPHostRoutes, opts.AdvertiseFloatingIPHostRoutes)
	th.AssertEquals(t, r.AdvertiseTenantNetworks, opts.AdvertiseTenantNetworks)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpSpeakerID := "ab01ade1-ae62-43c9-8a1f-3c24225b96d8"
	th.Mux.HandleFunc("/v2.0/bgp-speakers/"+bgpSpeakerID, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, GetBGPSpeakerResult)
		} else if r.Method == "PUT" {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, UpdateBGPSpeakerRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, UpdateBGPSpeakerResponse)
		} else {
			panic("Unexpected Request")
		}
	})

	m := make(map[string]string)
	m["Name"] = "testing-bgp-speaker"
	m["AdvertiseTenantNetworks"] = "false"

	opts := speaker.BuildUpdateOpts(fake.ServiceClient(), bgpSpeakerID, m)
	r, err := speaker.Update(fake.ServiceClient(), bgpSpeakerID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, r.Name, m["Name"])
	th.AssertEquals(t, r.AdvertiseTenantNetworks, false)
	th.AssertEquals(t, r.AdvertiseFloatingIPHostRoutes, true)
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

	err := speaker.Delete(fake.ServiceClient(), bgpSpeakerID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestAddBGPPeer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpSpeakerID := "ab01ade1-ae62-43c9-8a1f-3c24225b96d8"
	bgpPeerID := "f5884c7c-71d5-43a3-88b4-1742e97674aa"
	th.Mux.HandleFunc("/v2.0/bgp-speakers/"+bgpSpeakerID+"/add_bgp_peer", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AddRemoveBGPPeerJSON)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, AddRemoveBGPPeerJSON)
	})

	_, err := speaker.AddBGPPeer(fake.ServiceClient(), bgpSpeakerID, bgpPeerID).Extract()
	th.AssertNoErr(t, err)
}

func TestRemoveBGPPeer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpSpeakerID := "ab01ade1-ae62-43c9-8a1f-3c24225b96d8"
	bgpPeerID := "f5884c7c-71d5-43a3-88b4-1742e97674aa"
	th.Mux.HandleFunc("/v2.0/bgp-speakers/"+bgpSpeakerID+"/remove_bgp_peer", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AddRemoveBGPPeerJSON)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
	})

	_, err := speaker.RemoveBGPPeer(fake.ServiceClient(), bgpSpeakerID, bgpPeerID).Extract()
	if err != io.EOF {
		th.AssertNoErr(t, err)
	}
}

func TestGetAdvertisedRoutes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpSpeakerID := "ab01ade1-ae62-43c9-8a1f-3c24225b96d8"
	th.Mux.HandleFunc("/v2.0/bgp-speakers/"+bgpSpeakerID+"/get_advertised_routes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetAdvertisedRoutesResult)
	})

	count := 0
	speaker.GetAdvertisedRoutes(fake.ServiceClient(), bgpSpeakerID).EachPage(
		func(page pagination.Page) (bool, error) {
			count++
			actual, err := speaker.ExtractAdvertisedRoutes(page)

			if err != nil {
				t.Errorf("Failed to extract Advertised route: %v", err)
				return false, nil
			}

			expected := []speaker.AdvertisedRoute{
				speaker.AdvertisedRoute{NextHop: "172.17.128.212", Destination: "172.17.129.192/27"},
				speaker.AdvertisedRoute{NextHop: "172.17.128.218", Destination: "172.17.129.0/27"},
				speaker.AdvertisedRoute{NextHop: "172.17.128.231", Destination: "172.17.129.160/27"},
			}
			th.CheckDeepEquals(t, count, 1)
			th.CheckDeepEquals(t, expected, actual)
			return true, nil
		})
}

func TestAddGatewayNetwork(t *testing.T) {
}

func TestRemoveGatewayNetwork(t *testing.T) {
}
