package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavors"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const FlavorsListBody = `
{
	"flavors": [
        {
            "id": "4c82a610-8c7f-4a72-8cca-42f584e3f6d1",
            "name": "Basic",
            "description": "A basic standalone Octavia load balancer.",
            "enabled": true,
            "flavor_profile_id": "bdba88c7-beab-4fc9-a5dd-3635de59185b"
        },
		{
            "id": "0af3b9cc-9284-44c2-9494-0ec337fa31bb",
            "name": "Advance",
            "description": "A advance standalone Octavia load balancer.",
            "enabled": false,
            "flavor_profile_id": "c221abc6-a845-45a0-925c-27110c9d7bdc"
        }
    ]
}
`

const SingleFlavorBody = `
{
	"flavor": {
		"id": "5548c807-e6e8-43d7-9ea4-b38d34dd74a0",
		"name": "Basic",
		"description": "A basic standalone Octavia load balancer.",
		"enabled": true,
		"flavor_profile_id": "9daa2768-74e7-4d13-bf5d-1b8e0dc239e1"
	}
}
`

const PostUpdateFlavorBody = `
{
	"flavor": {
		"id": "5548c807-e6e8-43d7-9ea4-b38d34dd74a0",
		"name": "Basic v2",
		"description": "Rename flavor",
		"enabled": false,
		"flavor_profile_id": "9daa2768-74e7-4d13-bf5d-1b8e0dc239e1"
	}
}
`

var (
	FlavorBasic = flavors.Flavor{
		ID:              "4c82a610-8c7f-4a72-8cca-42f584e3f6d1",
		Name:            "Basic",
		Description:     "A basic standalone Octavia load balancer.",
		Enabled:         true,
		FlavorProfileId: "bdba88c7-beab-4fc9-a5dd-3635de59185b",
	}

	FlavorAdvance = flavors.Flavor{
		ID:              "0af3b9cc-9284-44c2-9494-0ec337fa31bb",
		Name:            "Advance",
		Description:     "A advance standalone Octavia load balancer.",
		Enabled:         false,
		FlavorProfileId: "c221abc6-a845-45a0-925c-27110c9d7bdc",
	}

	FlavorDb = flavors.Flavor{
		ID:              "5548c807-e6e8-43d7-9ea4-b38d34dd74a0",
		Name:            "Basic",
		Description:     "A basic standalone Octavia load balancer.",
		Enabled:         true,
		FlavorProfileId: "9daa2768-74e7-4d13-bf5d-1b8e0dc239e1",
	}

	FlavorUpdated = flavors.Flavor{
		ID:              "5548c807-e6e8-43d7-9ea4-b38d34dd74a0",
		Name:            "Basic v2",
		Description:     "Rename flavor",
		Enabled:         false,
		FlavorProfileId: "9daa2768-74e7-4d13-bf5d-1b8e0dc239e1",
	}
)

func HandleFlavorListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/flavors", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, FlavorsListBody)
		case "3a0d060b-fcec-4250-9ab6-940b806a12dd":
			fmt.Fprintf(w, `{ "flavors": [] }`)
		default:
			t.Fatalf("/v2.0/lbaas/flavors invoked with unexpected marker=[%s]", marker)
		}
	})
}

func HandleFlavorCreationSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/v2.0/lbaas/flavors", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
			"flavor": {
				"name": "Basic",
				"description": "A basic standalone Octavia load balancer.",
				"enabled": true,
				"flavor_profile_id": "9daa2768-74e7-4d13-bf5d-1b8e0dc239e1"
			}
		}`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, response)
	})
}

func HandleFlavorGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/flavors/5548c807-e6e8-43d7-9ea4-b38d34dd74a0", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, SingleFlavorBody)
	})
}

func HandleFlavorDeletionSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/flavors/5548c807-e6e8-43d7-9ea4-b38d34dd74a0", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleFlavorUpdateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v2.0/lbaas/flavors/5548c807-e6e8-43d7-9ea4-b38d34dd74a0", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `{
			"flavor": {
				"name": "Basic v2",
				"description": "Rename flavor",
				"enabled": true
			}
		}`)

		fmt.Fprintf(w, PostUpdateFlavorBody)
	})
}
