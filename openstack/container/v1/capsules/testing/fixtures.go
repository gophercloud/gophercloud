package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	fakeclient "github.com/gophercloud/gophercloud/testhelper/client"
)

// HandleImageGetSuccessfully test setup
func HandleCapsuleGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/capsules/cc654059-1a77-47a3-bfcf-715bde5aad9e", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"uuid": "cc654059-1a77-47a3-bfcf-715bde5aad9e",
			"status": "Running",
			"id": 1,
			"user_id": "d33b18c384574fd2a3299447aac285f0",
			"project_id": "6b8ffef2a0ac42ee87887b9cc98bdf68",
			"cpu": 1,
			"memory": "1024M",
			"meta_name": "test",
			"meta_labels": {"web": "app"},
			"created_at": "2018-01-12 09:37:25+00:00",
			"updated_at": "2018-01-12 09:37:26+00:00",
			"links": [
				{
				  "href": "http://10.10.10.10/v1/capsules/cc654059-1a77-47a3-bfcf-715bde5aad9e",
				  "rel": "self"
				},
				{
				  "href": "http://10.10.10.10/capsules/cc654059-1a77-47a3-bfcf-715bde5aad9e",
				  "rel": "bookmark"
				}
			],
			"capsule_version": "beta",
			"restart_policy":  "always",
			"containers_uuids": ["1739e28a-d391-4fd9-93a5-3ba3f29a4c9b"],
			"addresses": {
				"b1295212-64e1-471d-aa01-25ff46f9818d": [
					{
						"version": 4,
						"preserve_on_delete": false,
						"addr": "172.24.4.11",
						"port": "8439060f-381a-4386-a518-33d5a4058636",
						"subnet_id": "4a2bcd64-93ad-4436-9f48-3a7f9b267e0a"
					}
				]
			},
			"volumes_info": {
				"67618d54-dd55-4f7e-91b3-39ffb3ba7f5f": [
					"1739e28a-d391-4fd9-93a5-3ba3f29a4c9b"
				]
			},
			"host": "test-host",
			"status_reason": "No reason",
			"containers": [
				{
					"addresses": {
						"b1295212-64e1-471d-aa01-25ff46f9818d": [
							{
								"version": 4,
								"preserve_on_delete": false,
								"addr": "172.24.4.11",
								"port": "8439060f-381a-4386-a518-33d5a4058636",
								"subnet_id": "4a2bcd64-93ad-4436-9f48-3a7f9b267e0a"
							}
						]
					},
					"image": "test",
					"labels": {"foo": "bar"},
					"created_at": "2018-01-12 09:37:25+00:00",
					"updated_at": "2018-01-12 09:37:26+00:00",
					"workdir": "/root",
					"disk": 0,
					"id": 1,
					"security_groups": ["default"],
					"image_pull_policy": "ifnotpresent",
					"task_state": "Creating",
					"user_id": "d33b18c384574fd2a3299447aac285f0",
					"project_id": "6b8ffef2a0ac42ee87887b9cc98bdf68",
					"uuid": "1739e28a-d391-4fd9-93a5-3ba3f29a4c9b",
					"hostname": "test-hostname",
					"environment": {"USER1": "test"},
					"memory": "1024M",
					"status": "Running",
					"auto_remove": false,
					"container_id": "5109ebe2ca595777e994416208bd681b561b25ce493c34a234a1b68457cb53fb",
					"websocket_url": "ws://10.10.10.10/",
					"auto_heal": false,
					"host": "test-host",
					"image_driver": "docker",
					"status_detail": "Just created",
					"status_reason": "No reason",
					"websocket_token": "2ba16a5a-552f-422f-b511-bd786102691f",
					"name": "test-demo-omicron-13",
					"restart_policy": {
						"MaximumRetryCount": "0",
						"Name": "always"
					},
					"ports": [80],
					"meta": {"key1": "value1"},
					"command": "testcmd",
					"capsule_id": 1,
					"runtime": "runc",
					"cpu": 1,
					"interactive": true
				}
			]
		}`)
	})
}

// HandleCapsuleCreateSuccessfully creates an HTTP handler at `/capsules` on the test handler mux
// that responds with a `Create` response.
func HandleCapsuleCreateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/capsules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, `{}`)
	})
}
