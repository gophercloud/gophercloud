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

var (
	ProviderAmphora = providers.Provider{
		Name:        "amphora",
		Description: "The Octavia Amphora driver.",
	}
	ProviderOVN = providers.Provider{
		Name:        "ovn",
		Description: "The Octavia OVN driver",
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
