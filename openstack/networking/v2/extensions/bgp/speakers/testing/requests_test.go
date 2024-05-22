package testing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/speakers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
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

	err := speakers.List(fake.ServiceClient()).EachPage(
		context.TODO(),
		func(_ context.Context, page pagination.Page) (bool, error) {
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
	th.AssertNoErr(t, err)
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

	s, err := speakers.Get(context.TODO(), fake.ServiceClient(), bgpSpeakerID).Extract()
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
	r, err := speakers.Create(context.TODO(), fake.ServiceClient(), opts).Extract()
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

	err := speakers.Delete(context.TODO(), fake.ServiceClient(), bgpSpeakerID).ExtractErr()
	th.AssertNoErr(t, err)
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

	opts := speakers.UpdateOpts{
		Name:                          "testing-bgp-speaker",
		AdvertiseTenantNetworks:       false,
		AdvertiseFloatingIPHostRoutes: true,
	}

	r, err := speakers.Update(context.TODO(), fake.ServiceClient(), bgpSpeakerID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, r.Name, opts.Name)
	th.AssertEquals(t, r.AdvertiseTenantNetworks, opts.AdvertiseTenantNetworks)
	th.AssertEquals(t, r.AdvertiseFloatingIPHostRoutes, opts.AdvertiseFloatingIPHostRoutes)
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

	opts := speakers.AddBGPPeerOpts{BGPPeerID: bgpPeerID}
	r, err := speakers.AddBGPPeer(context.TODO(), fake.ServiceClient(), bgpSpeakerID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bgpPeerID, r.BGPPeerID)
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
		w.WriteHeader(http.StatusOK)
	})

	opts := speakers.RemoveBGPPeerOpts{BGPPeerID: bgpPeerID}
	err := speakers.RemoveBGPPeer(context.TODO(), fake.ServiceClient(), bgpSpeakerID, opts).ExtractErr()
	th.AssertEquals(t, err, io.EOF)
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
	err := speakers.GetAdvertisedRoutes(fake.ServiceClient(), bgpSpeakerID).EachPage(
		context.TODO(),
		func(_ context.Context, page pagination.Page) (bool, error) {
			count++
			actual, err := speakers.ExtractAdvertisedRoutes(page)

			if err != nil {
				t.Errorf("Failed to extract Advertised route: %v", err)
				return false, nil
			}

			expected := []speakers.AdvertisedRoute{
				{NextHop: "172.17.128.212", Destination: "172.17.129.192/27"},
				{NextHop: "172.17.128.218", Destination: "172.17.129.0/27"},
				{NextHop: "172.17.128.231", Destination: "172.17.129.160/27"},
			}
			th.CheckDeepEquals(t, count, 1)
			th.CheckDeepEquals(t, expected, actual)
			return true, nil
		})
	th.AssertNoErr(t, err)
}

func TestAddGatewayNetwork(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpSpeakerID := "ab01ade1-ae62-43c9-8a1f-3c24225b96d8"
	networkID := "ac13bb26-6219-49c3-a880-08847f6830b7"
	th.Mux.HandleFunc("/v2.0/bgp-speakers/"+bgpSpeakerID+"/add_gateway_network", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AddRemoveGatewayNetworkJSON)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, AddRemoveGatewayNetworkJSON)
	})

	opts := speakers.AddGatewayNetworkOpts{NetworkID: networkID}
	r, err := speakers.AddGatewayNetwork(context.TODO(), fake.ServiceClient(), bgpSpeakerID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, r.NetworkID, networkID)
}

func TestRemoveGatewayNetwork(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpSpeakerID := "ab01ade1-ae62-43c9-8a1f-3c24225b96d8"
	networkID := "ac13bb26-6219-49c3-a880-08847f6830b7"
	th.Mux.HandleFunc("/v2.0/bgp-speakers/"+bgpSpeakerID+"/remove_gateway_network", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AddRemoveGatewayNetworkJSON)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
	})

	opts := speakers.RemoveGatewayNetworkOpts{NetworkID: networkID}
	err := speakers.RemoveGatewayNetwork(context.TODO(), fake.ServiceClient(), bgpSpeakerID, opts).ExtractErr()
	th.AssertEquals(t, err, io.EOF)
}
