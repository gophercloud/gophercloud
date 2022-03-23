package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/limits"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

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

// FirstLimit is the first limit in the List request.
var FirstLimit = limits.Limit{
	ResourceName: "volume",
	Links: map[string]interface{}{
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
	Links: map[string]interface{}{
		"self": "http://10.3.150.25/identity/v3/limits/3229b3849f584faea483d6851f7aab05",
	},
	ServiceID:     "9408080f1970482aa0e38bc2d4ea34b7",
	ProjectID:     "3a705b9f56bb439381b43c4fe59dccce",
	ID:            "3229b3849f584faea483d6851f7aab05",
	ResourceLimit: 5,
}

// ExpectedLimitsSlice is the slice of limits expected to be returned from ListOutput.
var ExpectedLimitsSlice = []limits.Limit{FirstLimit, SecondLimit}

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
