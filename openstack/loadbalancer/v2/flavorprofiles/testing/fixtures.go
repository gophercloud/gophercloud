package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavorprofiles"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const FlavorProfilesListBody = `
{
	"flavorprofiles": [
        {
            "id": "c55d080d-af45-47ee-b48c-4caa5e87724f",
            "name": "amphora-single",
            "provider_name": "amphora",
            "flavor_data": "{\"loadbalancer_topology\": \"SINGLE\"}"
        },
		{
            "id": "f78d2815-3714-4b6e-91d8-cf821ba01017",
            "name": "amphora-act-stdby",
            "provider_name": "amphora",
            "flavor_data": "{\"loadbalancer_topology\": \"ACTIVE_STANDBY\"}"
        }
    ]
}
`

const SingleFlavorProfileBody = `
{
	"flavorprofile": {
		"id": "dcd65be5-f117-4260-ab3d-b32cc5bd1272",
		"name": "amphora-test",
		"provider_name": "amphora",
		"flavor_data": "{\"loadbalancer_topology\": \"ACTIVE_STANDBY\"}"
	}
}
`

const PostUpdateFlavorBody = `
{
	"flavorprofile": {
		"id": "dcd65be5-f117-4260-ab3d-b32cc5bd1272",
		"name": "amphora-test-updated",
		"provider_name": "amphora",
		"flavor_data": "{\"loadbalancer_topology\": \"SINGLE\"}"
	}
}
`

var (
	FlavorProfileSingle = flavorprofiles.FlavorProfile{
		ID:           "c55d080d-af45-47ee-b48c-4caa5e87724f",
		Name:         "amphora-single",
		ProviderName: "amphora",
		FlavorData:   "{\"loadbalancer_topology\": \"SINGLE\"}",
	}

	FlavorProfileAct = flavorprofiles.FlavorProfile{
		ID:           "f78d2815-3714-4b6e-91d8-cf821ba01017",
		Name:         "amphora-act-stdby",
		ProviderName: "amphora",
		FlavorData:   "{\"loadbalancer_topology\": \"ACTIVE_STANDBY\"}",
	}

	FlavorDb = flavorprofiles.FlavorProfile{
		ID:           "dcd65be5-f117-4260-ab3d-b32cc5bd1272",
		Name:         "amphora-test",
		ProviderName: "amphora",
		FlavorData:   "{\"loadbalancer_topology\": \"ACTIVE_STANDBY\"}",
	}

	FlavorUpdated = flavorprofiles.FlavorProfile{
		ID:           "dcd65be5-f117-4260-ab3d-b32cc5bd1272",
		Name:         "amphora-test-updated",
		ProviderName: "amphora",
		FlavorData:   "{\"loadbalancer_topology\": \"SINGLE\"}",
	}
)

func HandleFlavorProfileListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/flavorprofiles", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, FlavorProfilesListBody)
		case "3a0d060b-fcec-4250-9ab6-940b806a12dd":
			fmt.Fprintf(w, `{ "flavors": [] }`)
		default:
			t.Fatalf("/v2.0/lbaas/flavors invoked with unexpected marker=[%s]", marker)
		}
	})
}

func HandleFlavorProfileCreationSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/v2.0/lbaas/flavorprofiles", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
			"flavorprofile": {
				"name": "amphora-test",
				"provider_name": "amphora",
				"flavor_data": "{\"loadbalancer_topology\": \"ACTIVE_STANDBY\"}"
			}
		}`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, response)
	})
}

func HandleFlavorProfileGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/flavorprofiles/dcd65be5-f117-4260-ab3d-b32cc5bd1272", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, SingleFlavorProfileBody)
	})
}

func HandleFlavorProfileDeletionSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/flavorprofiles/dcd65be5-f117-4260-ab3d-b32cc5bd1272", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleFlavorProfileUpdateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/flavorprofiles/dcd65be5-f117-4260-ab3d-b32cc5bd1272", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `{
			"flavorprofile": {
				"name": "amphora-test-updated",
				"provider_name": "amphora",
				"flavor_data": "{\"loadbalancer_topology\": \"SINGLE\"}"
			}
		}`)

		fmt.Fprintf(w, PostUpdateFlavorBody)
	})
}
