package v2

// TODO

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fakeclient "github.com/rackspace/gophercloud/testhelper/client"
)

func HandleImageCreationSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)
		th.TestJSONRequest(t, r, `{
			"id": "e7db3b45-8db7-47ad-8109-3fb55c2c24fd",
			"name": "Ubuntu 12.10",
			"tags": [
				"ubuntu",
				"quantal"
			]
		}`)

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "queued",
			"name": "Ubuntu 12.10",
			"tags": ["ubuntu","quantal"],
			"container_format": "bare",
			"created_at": "2014-11-11T20:47:55Z",
			"disk_format": "qcow2",
			"updated_at": "2014-11-11T20:47:55Z",
			"visibility": "private",
			"self": "/v2/images/e7db3b45-8db7-47ad-8109-3fb55c2c24fd",
			"min_disk": 0,
			"protected": false,
			"id": "e7db3b45-8db7-47ad-8109-3fb55c2c24fd",
			"file": "/v2/images/e7db3b45-8db7-47ad-8109-3fb55c2c24fd/file",
			"owner": "b4eedccc6fb74fa8a7ad6b08382b852b",
			"min_ram": 0,
			"schema": "/v2/schemas/image",
			"size": "None",
			"checksum": "None",
			"virtual_size": "None"
		}`)
	})
}

func HandleImageGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/images/1bea47ed-f6a9-463b-b423-14b9cca9ad27", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "active",
			"name": "cirros-0.3.2-x86_64-disk",
			"tags": [],
			"container_format": "bare",
			"created_at": "2014-05-05T17:15:10Z",
			"disk_format": "qcow2",
			"updated_at": "2014-05-05T17:15:11Z",
			"visibility": "public",
			"self": "/v2/images/1bea47ed-f6a9-463b-b423-14b9cca9ad27",
			"min_disk": 0,
			"protected": false,
			"id": "1bea47ed-f6a9-463b-b423-14b9cca9ad27",
			"file": "/v2/images/1bea47ed-f6a9-463b-b423-14b9cca9ad27/file",
			"checksum": "64d7c1cd2b6f60c92c14662941cb7913",
			"owner": "5ef70662f8b34079a6eddb8da9d75fe8",
			"size": 13167616,
			"min_ram": 0,
			"schema": "/v2/schemas/image",
			"virtual_size": "None"
		}`)
	})
}

func HandleImageDeleteSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/images/1bea47ed-f6a9-463b-b423-14b9cca9ad27", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)
		
		w.WriteHeader(http.StatusNoContent)
	})
}
