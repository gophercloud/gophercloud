package testing

import (
	"fmt"
	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/speaker"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
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
}

func TestUpdate(t *testing.T) {
}

func TestDelete(t *testing.T) {
}

func TestAddBGPPeer(t *testing.T) {
}

func TestRemoveBGPPeer(t *testing.T) {
}

func TestGetAdvertisedRoutes(t *testing.T) {
}

func TestAddGatewayNetwork(t *testing.T) {
}

func TestRemoveGatewayNetwork(t *testing.T) {
}
