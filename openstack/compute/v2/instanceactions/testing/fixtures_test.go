package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/instanceactions"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ListExpected represents an expected repsonse from a List request.
var ListExpected = []instanceactions.InstanceAction{
	{
		Action:       "stop",
		InstanceUUID: "fcd19ef2-b593-40b1-90a5-fc31063fa95c",
		Message:      "",
		ProjectID:    "6f70656e737461636b20342065766572",
		RequestID:    "req-f8a59f03-76dc-412f-92c2-21f8612be728",
		StartTime:    time.Date(2018, 04, 25, 1, 26, 29, 000000, time.UTC),
		UserID:       "admin",
	},
	{
		Action:       "create",
		InstanceUUID: "fcd19ef2-b593-40b1-90a5-fc31063fa95c",
		Message:      "test",
		ProjectID:    "6f70656e737461636b20342065766572",
		RequestID:    "req-50189019-626d-47fb-b944-b8342af09679",
		StartTime:    time.Date(2018, 04, 25, 1, 26, 25, 000000, time.UTC),
		UserID:       "admin",
	},
}

// HandleInstanceActionListSuccessfully sets up the test server to respond to a ListAddresses request.
func HandleInstanceActionListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/servers/asdfasdfasdf/os-instance-actions", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"instanceActions": [
				{
					"action": "stop",
					"instance_uuid": "fcd19ef2-b593-40b1-90a5-fc31063fa95c",
					"message": null,
					"project_id": "6f70656e737461636b20342065766572",
					"request_id": "req-f8a59f03-76dc-412f-92c2-21f8612be728",
					"start_time": "2018-04-25T01:26:29.000000",
					"user_id": "admin"
				},
				{
					"action": "create",
					"instance_uuid": "fcd19ef2-b593-40b1-90a5-fc31063fa95c",
					"message": "test",
					"project_id": "6f70656e737461636b20342065766572",
					"request_id": "req-50189019-626d-47fb-b944-b8342af09679",
					"start_time": "2018-04-25T01:26:25.000000",
					"user_id": "admin"
				}
			]
		}`)
	})
}

var (
	expectedUpdateAt    = time.Date(2018, 04, 25, 1, 26, 36, 0, time.UTC)
	expectedEventHost   = "compute"
	expectedEventHostID = "2091634baaccdc4c5a1d57069c833e402921df696b7f970791b12ec6"
	expectedEvents      = []instanceactions.Event{{
		Event:      "compute_stop_instance",
		Host:       &expectedEventHost,
		HostID:     &expectedEventHostID,
		Result:     "Success",
		StartTime:  time.Date(2018, 04, 25, 1, 26, 36, 0, time.UTC),
		FinishTime: time.Date(2018, 04, 25, 1, 26, 36, 0, time.UTC),
		Traceback:  "",
	}}
)

// GetExpected represents an expected repsonse from a Get request.
var GetExpected = instanceactions.InstanceActionDetail{
	Action:       "stop",
	InstanceUUID: "4bf3473b-d550-4b65-9409-292d44ab14a2",
	Message:      "",
	ProjectID:    "6f70656e737461636b20342065766572",
	RequestID:    "req-0d819d5c-1527-4669-bdf0-ffad31b5105b",
	StartTime:    time.Date(2018, 04, 25, 1, 26, 36, 0, time.UTC),
	UpdatedAt:    &expectedUpdateAt,
	UserID:       "admin",
	Events:       &expectedEvents,
}

// HandleInstanceActionGetSuccessfully sets up the test server to respond to a Get request.
func HandleInstanceActionGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/servers/asdfasdfasdf/os-instance-actions/okzeorkmkfs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"instanceAction": 
			{
				"action": "stop",
				"events": [
					{
						"event": "compute_stop_instance",
						"finish_time": "2018-04-25T01:26:36.00000",
						"host": "compute",
						"hostId": "2091634baaccdc4c5a1d57069c833e402921df696b7f970791b12ec6",
						"result": "Success",
						"start_time": "2018-04-25T01:26:36.00000",
						"traceback": null
					}
				],
				"instance_uuid": "4bf3473b-d550-4b65-9409-292d44ab14a2",
				"message": null,
				"project_id": "6f70656e737461636b20342065766572",
				"request_id": "req-0d819d5c-1527-4669-bdf0-ffad31b5105b",
				"start_time": "2018-04-25T01:26:36.00000",
				"updated_at": "2018-04-25T01:26:36.00000",
				"user_id": "admin"
			}
		}`)
	})
}
