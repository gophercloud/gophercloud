package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/registeredlimits"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ListOutput provides a single page of List results.
const ListOutput = `
{
    "links": {
        "self": "http://10.3.150.25/identity/v3/registered_limits",
        "previous": null,
        "next": null
    },
    "registered_limits": [
        {
            "resource_name": "volume",
            "region_id": "RegionOne",
            "links": {
                "self": "http://10.3.150.25/identity/v3/registered_limits/25a04c7a065c430590881c646cdcdd58"
            },
            "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
            "id": "25a04c7a065c430590881c646cdcdd58",
            "default_limit": 11,
            "description": "Number of volumes for service 9408080f1970482aa0e38bc2d4ea34b7"
        },
        {
            "resource_name": "snapshot",
            "region_id": "RegionOne",
            "links": {
                "self": "http://10.3.150.25/identity/v3/registered_limits/3229b3849f584faea483d6851f7aab05"
            },
            "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
            "id": "3229b3849f584faea483d6851f7aab05",
            "default_limit": 5,
            "description": null
        }
    ]
}
`

// GetOutput provides a Get result.
const GetOutput = `
{
    "registered_limit": {
        "id": "3229b3849f584faea483d6851f7aab05",
        "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
        "region_id": "RegionOne",
        "resource_name": "snapshot",
        "default_limit": 5,
        "description": null,
        "links": {
            "self": "http://10.3.150.25/identity/v3/registered_limits/3229b3849f584faea483d6851f7aab05"
        }
	}
}
`

const CreateRequest = `
{
    "registered_limits":[
        {
            "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
            "region_id": "RegionOne",
            "resource_name": "snapshot",
            "default_limit": 5
        },
        {
            "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
            "region_id": "RegionOne",
            "resource_name": "volume",
            "default_limit": 11,
            "description": "Number of volumes for service 9408080f1970482aa0e38bc2d4ea34b7"
        }
    ]
}
`

// UpdateRequest provides the input to an Update request.
const UpdateRequest = `
{
    "registered_limit": {
	    "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
		"default_limit": 15,
		"resource_name": "volumes"
    }
}
`

// UpdateOutput provides an Update response.
const UpdateOutput = `
{
    "registered_limit": {
        "id": "3229b3849f584faea483d6851f7aab05",
        "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
        "region_id": "RegionOne",
        "resource_name": "volumes",
        "default_limit": 15,
        "description": "Number of volumes for service 9408080f1970482aa0e38bc2d4ea34b7",
        "links": {
            "self": "http://10.3.150.25/identity/v3/registered_limits/3229b3849f584faea483d6851f7aab05"
        }
    }
}
`

const CreateOutput = ListOutput

// FirstLimit is the first limit in the List request.
var FirstRegisteredLimit = registeredlimits.RegisteredLimit{
	ResourceName: "volume",
	RegionID:     "RegionOne",
	Links: map[string]any{
		"self": "http://10.3.150.25/identity/v3/registered_limits/25a04c7a065c430590881c646cdcdd58",
	},
	ServiceID:    "9408080f1970482aa0e38bc2d4ea34b7",
	ID:           "25a04c7a065c430590881c646cdcdd58",
	DefaultLimit: 11,
	Description:  "Number of volumes for service 9408080f1970482aa0e38bc2d4ea34b7",
}

// SecondLimit is the second limit in the List request.
var SecondRegisteredLimit = registeredlimits.RegisteredLimit{
	ResourceName: "snapshot",
	RegionID:     "RegionOne",
	Links: map[string]any{
		"self": "http://10.3.150.25/identity/v3/registered_limits/3229b3849f584faea483d6851f7aab05",
	},
	ServiceID:    "9408080f1970482aa0e38bc2d4ea34b7",
	ID:           "3229b3849f584faea483d6851f7aab05",
	DefaultLimit: 5,
}

// UpdatedSecondRegisteredLimit is a Registered Limit Fixture.
var UpdatedSecondRegisteredLimit = registeredlimits.RegisteredLimit{
	ResourceName: "volumes",
	RegionID:     "RegionOne",
	Links: map[string]any{
		"self": "http://10.3.150.25/identity/v3/registered_limits/3229b3849f584faea483d6851f7aab05",
	},
	ServiceID:    "9408080f1970482aa0e38bc2d4ea34b7",
	ID:           "3229b3849f584faea483d6851f7aab05",
	DefaultLimit: 15,
	Description:  "Number of volumes for service 9408080f1970482aa0e38bc2d4ea34b7",
}

// ExpectedRegisteredLimitsSlice is the slice of registered_limits expected to be returned from ListOutput.
var ExpectedRegisteredLimitsSlice = []registeredlimits.RegisteredLimit{FirstRegisteredLimit, SecondRegisteredLimit}

// HandleListRegisteredLimitsSuccessfully creates an HTTP handler at `/registered_limits` on the
// test handler mux that responds with a list of two registered limits.
func HandleListRegisteredLimitsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/registered_limits", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleGetRegisteredLimitSuccessfully creates an HTTP handler at `/registered_limits` on the
// test handler mux that responds with a single project.
func HandleGetRegisteredLimitSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/registered_limits/3229b3849f584faea483d6851f7aab05", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleCreateRegisteredLimitSuccessfully creates an HTTP handler at `/registered_limits` on the
// test handler mux that tests registered limit creation.
func HandleCreateRegisteredLimitSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/registered_limits", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, CreateOutput)
	})
}

// HandleDeleteRegisteredLimitSuccessfully creates an HTTP handler at `/registered_limits` on the
// test handler mux that tests registered_limit deletion.
func HandleDeleteRegisteredLimitSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/registered_limits/3229b3849f584faea483d6851f7aab05", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleUpdateRegisteredLimitSuccessfully creates an HTTP handler at `/registered_limits` on the
// test handler mux that tests registered limits updates.
func HandleUpdateRegisteredLimitSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/registered_limits/3229b3849f584faea483d6851f7aab05", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateRequest)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdateOutput)
	})
}
