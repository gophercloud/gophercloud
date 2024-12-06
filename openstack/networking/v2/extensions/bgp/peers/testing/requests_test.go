package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/peers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/bgp-peers",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, ListBGPPeersResult)
		})
	count := 0

	err := peers.List(fake.ServiceClient()).EachPage(
		context.TODO(),
		func(_ context.Context, page pagination.Page) (bool, error) {
			count++
			actual, err := peers.ExtractBGPPeers(page)

			if err != nil {
				t.Errorf("Failed to extract BGP Peers: %v", err)
				return false, nil
			}
			expected := []peers.BGPPeer{BGPPeer1, BGPPeer2}
			th.CheckDeepEquals(t, expected, actual)
			return true, nil
		})
	th.AssertNoErr(t, err)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	bgpPeerID := "afacc0e8-6b66-44e4-be53-a1ef16033ceb"
	th.Mux.HandleFunc("/v2.0/bgp-peers/"+bgpPeerID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, GetBGPPeerResult)
	})

	s, err := peers.Get(context.TODO(), fake.ServiceClient(), bgpPeerID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, *s, BGPPeer1)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/bgp-peers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, CreateResponse)
	})

	var opts peers.CreateOpts
	opts.AuthType = "md5"
	opts.Password = "notSoStrong"
	opts.RemoteAS = 20000
	opts.Name = "gophercloud-testing-bgp-peer"
	opts.PeerIP = "192.168.0.1"

	r, err := peers.Create(context.TODO(), fake.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, r.AuthType, opts.AuthType)
	th.AssertEquals(t, r.RemoteAS, opts.RemoteAS)
	th.AssertEquals(t, r.PeerIP, opts.PeerIP)
}

func TestDelete(t *testing.T) {
	bgpPeerID := "afacc0e8-6b66-44e4-be53-a1ef16033ceb"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/bgp-peers/"+bgpPeerID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	err := peers.Delete(context.TODO(), fake.ServiceClient(), bgpPeerID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUpdate(t *testing.T) {
	bgpPeerID := "afacc0e8-6b66-44e4-be53-a1ef16033ceb"
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/bgp-peers/"+bgpPeerID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateBGPPeerRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateBGPPeerResponse)
	})

	var opts peers.UpdateOpts
	opts.Name = "test-rename-bgp-peer"
	opts.Password = "superStrong"

	r, err := peers.Update(context.TODO(), fake.ServiceClient(), bgpPeerID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, r.Name, opts.Name)
}
