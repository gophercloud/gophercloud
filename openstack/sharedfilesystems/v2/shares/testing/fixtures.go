package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

const (
	shareEndpoint = "/shares"
	shareID       = "011d21e2-fbc3-4e4a-9993-9ea223f73264"
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

// MockDeleteResponse creates a mock delete response
func MockDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}

var getResponse = `{
    "share": {
        "links": [
            {
                "href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
                "rel": "self"
            },
            {
                "href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
                "rel": "bookmark"
            }
        ],
        "availability_zone": "nova",
        "share_network_id": "713df749-aac0-4a54-af52-10f6c991e80c",
        "share_server_id": "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
        "snapshot_id": null,
        "id": "011d21e2-fbc3-4e4a-9993-9ea223f73264",
        "size": 1,
        "share_type": "25747776-08e5-494f-ab40-a64b9d20d8f7",
        "share_type_name": "default",
        "consistency_group_id": "9397c191-8427-4661-a2e8-b23820dc01d4",
        "project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
        "metadata": {
            "project": "my_app",
            "aim": "doc"
        },
        "status": "available",
        "description": "My custom share London",
        "host": "manila2@generic1#GENERIC1",
        "has_replicas": false,
        "replication_type": null,
        "task_state": null,
        "is_public": true,
        "snapshot_support": true,
        "name": "my_test_share",
        "created_at": "2015-09-18T10:25:24.000000",
        "share_proto": "NFS",
        "volume_type": "default",
        "source_cgsnapshot_member_id": null
    }
}`

// MockGetResponse creates a mock get response
func MockGetResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, getResponse)
	})
}

var shortListResponse = `{
	"shares": [{
		"id": "d94a8548-2079-4be0-b21c-0a887acd31ca",
		"links": [{
			"href": "http://172.18.198.54:8786/v1/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
			"rel": "self"
		}, {
			"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
			"rel": "bookmark"
		}],
		"name": "My_share"
	}, {
		"id": "406ea93b-32e9-4907-a117-148b3945749f",
		"links": [{
			"href": "http://172.18.198.54:8786/v1/16e1ab15c35a457e9c2b2aa189f544e1/shares/406ea93b-32e9-4907-a117-148b3945749f",
			"rel": "self"
		}, {
			"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/406ea93b-32e9-4907-a117-148b3945749f",
			"rel": "bookmark"
		}],
		"name": "Share1"
	}]
}`

// MockListResponse creates a mock list response
func MockListResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		r.ParseForm()
		marker := r.Form.Get("offset")
		switch marker {
		case "":
			fmt.Fprintf(w, shortListResponse)
		default:
			fmt.Fprintf(w, `{"shares": []}`)
		}
	})
}

// MockListDetailResponse creates a mock list response with details
func MockListDetailResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		r.ParseForm()
		marker := r.Form.Get("offset")
		switch marker {
		case "":
			fallthrough
		case "1":
			fmt.Fprintf(w, `
{
    "shares": [
        {
            "links": [
                {
                    "href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/f45cc5b2-d1bb-4a3e-ba5b-5c4125613adc",
                    "rel": "self"
                },
                {
                    "href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/f45cc5b2-d1bb-4a3e-ba5b-5c4125613adc",
                    "rel": "bookmark"
                }
            ],
            "availability_zone": "nova",
            "share_network_id": "f9b2e754-ac01-4466-86e1-5c569424754e",
            "export_locations": [],
            "share_server_id": "87d8943a-f5da-47a4-b2f2-ddfa6794aa82",
            "snapshot_id": null,
            "id": "f45cc5b2-d1bb-4a3e-ba5b-5c4125613adc",
            "size": 1,
            "share_type": "25747776-08e5-494f-ab40-a64b9d20d8f7",
            "share_type_name": "default",
            "export_location": null,
            "consistency_group_id": "9397c191-8427-4661-a2e8-b23820dc01d4",
            "project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
            "metadata": {},
            "status": "error",
            "access_rules_status": "active",
            "description": "There is a share description.",
            "host": "manila2@generic1#GENERIC1",
            "task_state": null,
            "is_public": true,
            "snapshot_support": true,
            "name": "my_share4",
            "has_replicas": false,
            "replication_type": null,
            "created_at": "2015-09-16T18:19:50.000000",
            "share_proto": "NFS",
            "volume_type": "default",
            "source_cgsnapshot_member_id": null
        }
	]
}`)
		case "2":
			fmt.Fprintf(w, `
{
    "shares": [
        {
            "links": [
                {
                    "href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/c4a2ced4-2c9f-4ae1-adaa-6171833e64df",
                    "rel": "self"
                },
                {
                    "href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/c4a2ced4-2c9f-4ae1-adaa-6171833e64df",
                    "rel": "bookmark"
                }
            ],
            "availability_zone": "nova",
            "share_network_id": "f9b2e754-ac01-4466-86e1-5c569424754e",
            "export_locations": [
                "10.254.0.5:/shares/share-50ad5e7b-f6f1-4b78-a651-0812cef2bb67"
            ],
            "share_server_id": "87d8943a-f5da-47a4-b2f2-ddfa6794aa82",
            "snapshot_id": null,
            "id": "c4a2ced4-2c9f-4ae1-adaa-6171833e64df",
            "size": 1,
            "share_type": "25747776-08e5-494f-ab40-a64b9d20d8f7",
            "share_type_name": "default",
            "export_location": "10.254.0.5:/shares/share-50ad5e7b-f6f1-4b78-a651-0812cef2bb67",
            "consistency_group_id": "9397c191-8427-4661-a2e8-b23820dc01d4",
            "project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
            "metadata": {},
            "status": "available",
            "access_rules_status": "active",
            "description": "Changed description.",
            "host": "manila2@generic1#GENERIC1",
            "task_state": null,
            "is_public": true,
            "snapshot_support": true,
            "name": "my_share3",
            "has_replicas": false,
            "replication_type": null,
            "created_at": "2015-09-16T17:26:28.000000",
            "share_proto": "NFS",
            "volume_type": "default",
            "source_cgsnapshot_member_id": null
        }
    ]
}`)
		default:
			fmt.Fprintf(w, `
{
"shares": []
}`)
		}
	})
}
