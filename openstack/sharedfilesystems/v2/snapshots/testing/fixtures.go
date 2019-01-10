package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

const (
	snapshotEndpoint = "/snapshots"
	snapshotID       = "bc082e99-3bdb-4400-b95e-b85c7a41622c"
)

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

		r.ParseForm()
		marker := r.Form.Get("offset")

		switch marker {
		case "":
			fmt.Fprint(w, listDetailResponse)
		default:
			fmt.Fprint(w, listDetailEmptyResponse)
		}
	})
}
