package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/limits"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const GetEnforcementModelOutput = `
{
    "model": {
        "description": "Limit enforcement and validation does not take project hierarchy into consideration.",
        "name": "flat"
    }
}
`

// ListOutput provides a single page of List results.
const ListOutput = `
{
    "links": {
        "self": "http://10.3.150.25/identity/v3/limits",
        "previous": null,
        "next": null
    },
    "limits": [
        {
            "resource_name": "volume",
            "region_id": null,
            "links": {
                "self": "http://10.3.150.25/identity/v3/limits/25a04c7a065c430590881c646cdcdd58"
            },
            "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
            "project_id": "3a705b9f56bb439381b43c4fe59dccce",
            "domain_id": null,
            "id": "25a04c7a065c430590881c646cdcdd58",
            "resource_limit": 11,
            "description": "Number of volumes for project 3a705b9f56bb439381b43c4fe59dccce"
        },
        {
            "resource_name": "snapshot",
            "region_id": "RegionOne",
            "links": {
                "self": "http://10.3.150.25/identity/v3/limits/3229b3849f584faea483d6851f7aab05"
            },
            "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
            "project_id": "3a705b9f56bb439381b43c4fe59dccce",
            "domain_id": null,
            "id": "3229b3849f584faea483d6851f7aab05",
            "resource_limit": 5,
            "description": null
        }
    ]
}
`

const CreateRequest = `
{
    "limits":[
        {
            "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
            "project_id": "3a705b9f56bb439381b43c4fe59dccce",
            "region_id": "RegionOne",
            "resource_name": "snapshot",
            "resource_limit": 5
        },
        {
            "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
            "domain_id": "edbafc92be354ffa977c58aa79c7bdb2",
            "resource_name": "volume",
            "resource_limit": 11,
            "description": "Number of volumes for project 3a705b9f56bb439381b43c4fe59dccce"
        }
    ]
}
`

const GetOutput = `
{
    "limit": {
        "resource_name": "volume",
        "region_id": null,
        "links": {
            "self": "http://10.3.150.25/identity/v3/limits/25a04c7a065c430590881c646cdcdd58"
        },
        "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
        "project_id": "3a705b9f56bb439381b43c4fe59dccce",
        "id": "25a04c7a065c430590881c646cdcdd58",
        "resource_limit": 11,
        "description": "Number of volumes for project 3a705b9f56bb439381b43c4fe59dccce"
    }
}
`

const UpdateRequest = `
{
    "limit": {
        "resource_limit": 5,
        "description": "Number of snapshots for project 3a705b9f56bb439381b43c4fe59dccce"
    }
}
`

const UpdateOutput = `
{
    "limit": {
        "resource_name": "snapshot",
        "region_id": "RegionOne",
        "links": {
            "self": "http://10.3.150.25/identity/v3/limits/3229b3849f584faea483d6851f7aab05"
        },
        "service_id": "9408080f1970482aa0e38bc2d4ea34b7",
        "project_id": "3a705b9f56bb439381b43c4fe59dccce",
        "id": "3229b3849f584faea483d6851f7aab05",
        "resource_limit": 5,
        "description": "Number of snapshots for project 3a705b9f56bb439381b43c4fe59dccce"
    }
}
`

// Model is the enforcement model in the GetEnforcementModel request.
var Model = limits.EnforcementModel{
	Name:        "flat",
	Description: "Limit enforcement and validation does not take project hierarchy into consideration.",
}

const CreateOutput = ListOutput

// FirstLimit is the first limit in the List request.
var FirstLimit = limits.Limit{
	ResourceName: "volume",
	Links: map[string]any{
		"self": "http://10.3.150.25/identity/v3/limits/25a04c7a065c430590881c646cdcdd58",
	},
	ServiceID:     "9408080f1970482aa0e38bc2d4ea34b7",
	ProjectID:     "3a705b9f56bb439381b43c4fe59dccce",
	ID:            "25a04c7a065c430590881c646cdcdd58",
	ResourceLimit: 11,
	Description:   "Number of volumes for project 3a705b9f56bb439381b43c4fe59dccce",
}

// SecondLimit is the second limit in the List request.
var SecondLimit = limits.Limit{
	ResourceName: "snapshot",
	RegionID:     "RegionOne",
	Links: map[string]any{
		"self": "http://10.3.150.25/identity/v3/limits/3229b3849f584faea483d6851f7aab05",
	},
	ServiceID:     "9408080f1970482aa0e38bc2d4ea34b7",
	ProjectID:     "3a705b9f56bb439381b43c4fe59dccce",
	ID:            "3229b3849f584faea483d6851f7aab05",
	ResourceLimit: 5,
}

// SecondLimitUpdated is the updated limit in the Update request.
var SecondLimitUpdated = limits.Limit{
	ResourceName: "snapshot",
	RegionID:     "RegionOne",
	Links: map[string]any{
		"self": "http://10.3.150.25/identity/v3/limits/3229b3849f584faea483d6851f7aab05",
	},
	ServiceID:     "9408080f1970482aa0e38bc2d4ea34b7",
	ProjectID:     "3a705b9f56bb439381b43c4fe59dccce",
	ID:            "3229b3849f584faea483d6851f7aab05",
	ResourceLimit: 5,
	Description:   "Number of snapshots for project 3a705b9f56bb439381b43c4fe59dccce",
}

// ExpectedLimitsSlice is the slice of limits expected to be returned from ListOutput.
var ExpectedLimitsSlice = []limits.Limit{FirstLimit, SecondLimit}

// HandleGetEnforcementModelSuccessfully creates an HTTP handler at `/limits/model` on the
// test handler mux that responds with a enforcement model.
func HandleGetEnforcementModelSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/limits/model", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetEnforcementModelOutput)
	})
}

// HandleListLimitsSuccessfully creates an HTTP handler at `/limits` on the
// test handler mux that responds with a list of two limits.
func HandleListLimitsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/limits", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleCreateLimitSuccessfully creates an HTTP handler at `/limits` on the
// test handler mux that tests limit creation.
func HandleCreateLimitSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/limits", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, CreateOutput)
	})
}

// HandleGetLimitSuccessfully creates an HTTP handler at `/limits` on the
// test handler mux that responds with a single limit.
func HandleGetLimitSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/limits/25a04c7a065c430590881c646cdcdd58", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleUpdateLimitSuccessfully creates an HTTP handler at `/limits` on the
// test handler mux that tests limit update.
func HandleUpdateLimitSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/limits/3229b3849f584faea483d6851f7aab05", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateRequest)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdateOutput)
	})
}

// HandleDeleteLimitSuccessfully creates an HTTP handler at `/limits` on the
// test handler mux that tests limit deletion.
func HandleDeleteLimitSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/limits/3229b3849f584faea483d6851f7aab05", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}
