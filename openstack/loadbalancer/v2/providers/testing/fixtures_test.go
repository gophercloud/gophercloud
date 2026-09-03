package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/providers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ProvidersListBody contains the canned body of a provider list response.
const ProvidersListBody = `
{
	"providers":[
	         {
			"name": "amphora",
			"description": "The Octavia Amphora driver."
		},
		{
			"name": "ovn",
			"description": "The Octavia OVN driver"
		}
	]
}
`

// FlavorCapabilitiesListBody contains the canned body of a provider flavor
// capabilities response.
const FlavorCapabilitiesListBody = `
{
	"flavor_capabilities": [
		{
			"name": "loadbalancer_topology",
			"description": "The load balancer topology."
		},
		{
			"name": "compute_flavor",
			"description": "The compute driver flavor ID."
		}
	]
}
`

// AvailabilityZoneCapabilitiesListBody contains the canned body of a provider
// availability zone capabilities response.
const AvailabilityZoneCapabilitiesListBody = `
{
	"availability_zone_capabilities": [
		{
			"name": "compute_zone",
			"description": "The compute availability zone."
		},
		{
			"name": "volume_zone",
			"description": "The volume availability zone."
		}
	]
}
`

var (
	ProviderAmphora = providers.Provider{
		Name:        "amphora",
		Description: "The Octavia Amphora driver.",
	}
	ProviderOVN = providers.Provider{
		Name:        "ovn",
		Description: "The Octavia OVN driver",
	}
	FlavorCapabilities = []providers.Capability{
		{
			Name:        "loadbalancer_topology",
			Description: "The load balancer topology.",
		},
		{
			Name:        "compute_flavor",
			Description: "The compute driver flavor ID.",
		},
	}
	AvailabilityZoneCapabilities = []providers.Capability{
		{
			Name:        "compute_zone",
			Description: "The compute availability zone.",
		},
		{
			Name:        "volume_zone",
			Description: "The volume availability zone.",
		},
	}
)

// HandleProviderListSuccessfully sets up the test server to respond to a provider List request.
func HandleProviderListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v2.0/lbaas/providers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprint(w, ProvidersListBody)
		default:
			t.Fatalf("/v2.0/lbaas/providers invoked with unexpected marker=[%s]", marker)
		}
	})
}

// HandleFlavorCapabilitiesListSuccessfully sets up the test server to respond
// to a provider flavor capabilities request.
func HandleFlavorCapabilitiesListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v2.0/lbaas/providers/amphora/flavor_capabilities", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.AssertDeepEquals(t, []string{"name", "description"}, r.URL.Query()["fields"])

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, FlavorCapabilitiesListBody)
	})
}

// HandleAvailabilityZoneCapabilitiesListSuccessfully sets up the test server
// to respond to a provider availability zone capabilities request.
func HandleAvailabilityZoneCapabilitiesListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v2.0/lbaas/providers/amphora/availability_zone_capabilities", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.AssertDeepEquals(t, []string{"name", "description"}, r.URL.Query()["fields"])

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, AvailabilityZoneCapabilitiesListBody)
	})
}
