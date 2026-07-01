package testing

import (
	"fmt"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/migrations"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
	"net/http"
	"testing"
	"time"
)

// ListExpected represents an expected response from a List request.
var ListExpected = []migrations.Migration{
	{
		Id:                2,
		CreatedAt:         time.Date(2024, 11, 28, 8, 15, 6, 000000, time.UTC),
		DestCompute:       "hci0",
		DestHost:          "192.168.1.41",
		DestNode:          "hci0",
		InstanceID:        "6ba1f91a-50cc-48cd-9ce6-310991acb08d",
		NewInstanceTypeId: 1,
		OldInstanceTypeId: 1,
		SourceCompute:     "compute0",
		SourceNode:        "compute0",
		Status:            "completed",
		UpdatedAt:         time.Date(2024, 11, 28, 8, 15, 28, 000000, time.UTC),
		MigrationType:     "live-migration",
		Uuid:              "dde90a82-3059-4369-8bda-e0ba92713c54",
		UserId:            "admin",
		ProjectId:         "admin",
	},
	{
		Id:                1,
		CreatedAt:         time.Date(2024, 11, 28, 8, 10, 02, 000000, time.UTC),
		DestCompute:       "compute0",
		DestHost:          "192.168.1.42",
		DestNode:          "compute0",
		InstanceID:        "6ba1f91a-50cc-48cd-9ce6-310991acb08d",
		NewInstanceTypeId: 1,
		OldInstanceTypeId: 1,
		SourceCompute:     "hci0",
		SourceNode:        "hci0",
		Status:            "completed",
		UpdatedAt:         time.Date(2024, 11, 28, 8, 10, 34, 000000, time.UTC),
		MigrationType:     "live-migration",
		Uuid:              "dde90a82-3059-4369-8bda-e0ba92713c54",
		UserId:            "admin",
		ProjectId:         "admin",
	},
}

// HandleMigrationListSuccessfully sets up the test server to respond to a List request.
func HandleMigrationListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/os-migrations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"migrations": [
				{
					"id": 2, 
					"uuid": "dde90a82-3059-4369-8bda-e0ba92713c54", 
					"source_compute": "compute0", 
					"dest_compute": "hci0", 
					"source_node": "compute0", 
					"dest_node": "hci0", 
					"dest_host": "192.168.1.41", 
					"old_instance_type_id": 1, 
					"new_instance_type_id": 1, 
					"instance_uuid": "6ba1f91a-50cc-48cd-9ce6-310991acb08d", 
					"status": "completed", 
					"migration_type": "live-migration",
					"user_id": "admin", 
					"project_id": "admin", 
					"created_at": "2024-11-28T08:15:06.000000", 
					"updated_at": "2024-11-28T08:15:28.000000"
				},
				{
					"id": 1, 
					"uuid": "dde90a82-3059-4369-8bda-e0ba92713c54", 
					"source_compute": "hci0", 
					"dest_compute": "compute0", 
					"source_node": "hci0", 
					"dest_node": "compute0", 
					"dest_host": "192.168.1.42", 
					"old_instance_type_id": 1, 
					"new_instance_type_id": 1, 
					"instance_uuid": "6ba1f91a-50cc-48cd-9ce6-310991acb08d", 
					"status": "completed", 
					"migration_type": "live-migration",
					"user_id": "admin", 
					"project_id": "admin", 
					"created_at": "2024-11-28T08:10:02.000000", 
					"updated_at": "2024-11-28T08:10:34.000000"
				}
			], 
			"migrations_links": [
				{
					"rel": "next", 
					"href": "http://192.168.122.161:8774/v2.1/os-migrations?instance_uuid=6ba1f91a-50cc-48cd-9ce6-310991acb08d&limit=1&marker=dde90a82-3059-4369-8bda-e0ba92713c54"
				}
			]
		}`)
	})
}
