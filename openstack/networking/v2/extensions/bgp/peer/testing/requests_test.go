package testing

import (
	"fmt"
	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/peer"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"net/http"
	"testing"
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
			fmt.Fprintf(w, ListBGPPeersResult)
		})
	count := 0

	peer.List(fake.ServiceClient()).EachPage(
		func(page pagination.Page) (bool, error) {
			count++
			actual, err := peer.ExtractBGPPeers(page)

			if err != nil {
				t.Errorf("Failed to extract BGP Peers: %v", err)
				return false, nil
			}
			expected := []peer.BGPPeer{BGPPeer1, BGPPeer2}
			th.CheckDeepEquals(t, expected, actual)
			return true, nil
		})
}
