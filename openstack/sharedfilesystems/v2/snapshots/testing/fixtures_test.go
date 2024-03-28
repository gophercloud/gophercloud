package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const (
	snapshotEndpoint = "/snapshots"
	snapshotID       = "bc082e99-3bdb-4400-b95e-b85c7a41622c"
	shareID          = "19865c43-3b91-48c9-85a0-7ac4d6bb0efe"
)

var createRequest = `{
	"snapshot": {
		"share_id": "19865c43-3b91-48c9-85a0-7ac4d6bb0efe",
		"name": "test snapshot",
		"description": "test description"
	}
}`

var createResponse = `{
	"snapshot": {
		"status": "creating",
		"share_id": "19865c43-3b91-48c9-85a0-7ac4d6bb0efe",
		"description": "test description",
		"links": [
			{
				"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/9897f5ca-2559-4a4c-b761-d3439c0c9455",
				"rel": "self"
			},
			{
				"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/9897f5ca-2559-4a4c-b761-d3439c0c9455",
				"rel": "bookmark"
			}
		],
		"id": "bc082e99-3bdb-4400-b95e-b85c7a41622c",
		"size": 1,
		"user_id": "619e2ad074321cf246b03a89e95afee95fb26bb0b2d1fc7ba3bd30fcca25588a",
		"name": "test snapshot",
		"created_at": "2019-01-09T10:22:39.613550",
		"share_proto": "NFS",
		"project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
		"share_size": 1
	}
}`

// MockCreateResponse creates a mock response
func MockCreateResponse(t *testing.T) {
	th.Mux.HandleFunc(snapshotEndpoint, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, createRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, createResponse)
	})
}

// MockDeleteResponse creates a mock delete response
func MockDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc(snapshotEndpoint+"/"+snapshotID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}

var updateRequest = `{
		"snapshot": {
			"display_name": "my_new_test_snapshot",
			"display_description": ""
		}
}`

var updateResponse = `{
	"snapshot": {
		"status": "available",
		"share_id": "19865c43-3b91-48c9-85a0-7ac4d6bb0efe",
		"description": "",
		"links": [
			{
				"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/9897f5ca-2559-4a4c-b761-d3439c0c9455",
				"rel": "self"
			},
			{
				"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/9897f5ca-2559-4a4c-b761-d3439c0c9455",
				"rel": "bookmark"
			}
		],
		"id": "9897f5ca-2559-4a4c-b761-d3439c0c9455",
		"size": 1,
		"user_id": "619e2ad074321cf246b03a89e95afee95fb26bb0b2d1fc7ba3bd30fcca25588a",
		"name": "my_new_test_snapshot",
		"created_at": "2019-01-09T10:22:39.613550",
		"share_proto": "NFS",
		"project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
		"share_size": 1
	}
}`

func MockUpdateResponse(t *testing.T) {
	th.Mux.HandleFunc(snapshotEndpoint+"/"+snapshotID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, updateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, updateResponse)
	})
}

var getResponse = `{
	"snapshot": {
		"status": "available",
		"share_id": "19865c43-3b91-48c9-85a0-7ac4d6bb0efe",
		"description": null,
		"links": [
			{
				"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/bc082e99-3bdb-4400-b95e-b85c7a41622c",
				"rel": "self"
			},
			{
				"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/bc082e99-3bdb-4400-b95e-b85c7a41622c",
				"rel": "bookmark"
			}
		],
		"id": "bc082e99-3bdb-4400-b95e-b85c7a41622c",
		"size": 1,
		"user_id": "619e2ad074321cf246b03a89e95afee95fb26bb0b2d1fc7ba3bd30fcca25588a",
		"name": "new_app_snapshot",
		"created_at": "2019-01-06T11:11:02.000000",
		"share_proto": "NFS",
		"project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
		"share_size": 1
	}
}`

// MockGetResponse creates a mock get response
func MockGetResponse(t *testing.T) {
	th.Mux.HandleFunc(snapshotEndpoint+"/"+snapshotID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, getResponse)
	})
}

var listDetailResponse = `{
	"snapshots": [
		{
			"status": "available",
			"share_id": "19865c43-3b91-48c9-85a0-7ac4d6bb0efe",
			"description": null,
			"links": [
				{
					"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/bc082e99-3bdb-4400-b95e-b85c7a41622c",
					"rel": "self"
				},
				{
					"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/bc082e99-3bdb-4400-b95e-b85c7a41622c",
					"rel": "bookmark"
				}
			],
			"id": "bc082e99-3bdb-4400-b95e-b85c7a41622c",
			"size": 1,
			"user_id": "619e2ad074321cf246b03a89e95afee95fb26bb0b2d1fc7ba3bd30fcca25588a",
			"name": "new_app_snapshot",
			"created_at": "2019-01-06T11:11:02.000000",
			"share_proto": "NFS",
			"project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
			"share_size": 1
		}
	]
}`

var listDetailEmptyResponse = `{"snapshots": []}`

// MockListDetailResponse creates a mock detailed-list response
func MockListDetailResponse(t *testing.T) {
	th.Mux.HandleFunc(snapshotEndpoint+"/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("offset")

		switch marker {
		case "":
			fmt.Fprint(w, listDetailResponse)
		default:
			fmt.Fprint(w, listDetailEmptyResponse)
		}
	})
}

var resetStatusRequest = `{
                "reset_status": {
                        "status": "error"
                }
        }`

// MockResetStatusResponse creates a mock reset status snapshot response
func MockResetStatusResponse(t *testing.T) {
	th.Mux.HandleFunc(snapshotEndpoint+"/"+snapshotID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, resetStatusRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

var forceDeleteRequest = `{
                "force_delete": null
        }`

// MockForceDeleteResponse creates a mock force delete snapshot response
func MockForceDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc(snapshotEndpoint+"/"+snapshotID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, forceDeleteRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}
