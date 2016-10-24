package testing

import (
	"fmt"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
	"net/http"
	"testing"
)

const (
	shareEndpoint = "/shares"
)

var createRequest = `{
		"share": {
			"name": "my_test_share",
			"size": 1,
			"share_proto": "NFS"
		}
	}`

var createResponse = `{
		"share": {
			"name": "my_test_share",
			"share_proto": "NFS",
			"size": 1,
			"status": null,
			"share_server_id": null,
			"project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
			"share_type": "25747776-08e5-494f-ab40-a64b9d20d8f7",
			"share_type_name": "default",
			"availability_zone": null,
			"created_at": "2015-09-18T10:25:24.533287",
			"export_location": null,
			"links": [
				{
					"href": "http://172.18.198.54:8786/v1/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
					"rel": "self"
				},
				{
					"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
					"rel": "bookmark"
				}
			],
			"share_network_id": null,
			"export_locations": [],
			"host": null,
			"access_rules_status": "active",
			"has_replicas": false,
			"replication_type": null,
			"task_state": null,
			"snapshot_support": true,
			"consistency_group_id": "9397c191-8427-4661-a2e8-b23820dc01d4",
			"source_cgsnapshot_member_id": null,
			"volume_type": "default",
			"snapshot_id": null,
			"is_public": true,
			"metadata": {
				"project": "my_app",
				"aim": "doc"
			},
			"id": "011d21e2-fbc3-4e4a-9993-9ea223f73264",
			"description": "My custom share London"
		}
	}`

// MockCreateResponse creates a mock response
func MockCreateResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, createRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, createResponse)
	})
}
