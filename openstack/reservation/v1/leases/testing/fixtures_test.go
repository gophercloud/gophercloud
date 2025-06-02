package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const LeaseListBody = `
{
    "leases": [
        {
            "id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
            "name": "lease_foo",
            "start_date": "2017-12-26T12:00:00.000000",
            "end_date": "2017-12-27T12:00:00.000000",
            "status":"PENDING",
            "degraded": false,
            "user_id": "5434f637520d4c17bbf254af034b0320",
            "project_id": "aa45f56901ef45ee95e3d211097c0ea3",
            "trust_id": "b442a580b9504ababf305bf2b4c49512",
            "created_at": "2017-12-27 10:00:00",
            "updated_at": null,
            "reservations": [
                {
                    "id": "087bc740-6d2d-410b-9d47-c7b2b55a9d36",
                    "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                    "status": "pending",
                    "missing_resources": false,
                    "resources_changed": false,
                    "resource_id": "5e6c0e6e-f1e6-490b-baaf-50deacbbe371",
                    "resource_type": "physical:host",
                    "min": 4,
                    "max": 6,
                    "hypervisor_properties": "[\">=\", \"$vcpus\", \"4\"]",
                    "resource_properties": "",
                    "before_end": "default",
                    "created_at": "2017-12-27 10:00:00",
                    "updated_at": null
                },
                {
                    "id": "ddc45423-f863-4e4e-8e7a-51d27cfec962",
                    "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                    "status": "pending",
                    "missing_resources": false,
                    "resources_changed": false,
                    "resource_id": "0b901727-cca2-43ed-bcc8-c21b0982dcb1",
                    "resource_type": "virtual:instance",
                    "amount": 4,
                    "vcpus": 2,
                    "memory_mb": 4096,
                    "disk_gb": 100,
                    "affinity": false,
                    "resource_properties": "",
                    "flavor_id": "ddc45423-f863-4e4e-8e7a-51d27cfec962",
                    "server_group_id": "33cdfc42-5a04-4fcc-b190-1abebaa056bb",
                    "aggregate_id": 11,
                    "created_at": "2017-12-27 10:00:00",
                    "updated_at": null
                }
            ],
            "events": [
                {
                    "id": "188a8584-f832-4df9-9a4a-51e6364420ff",
                    "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                    "status": "UNDONE",
                    "event_type": "start_lease",
                    "time": "2017-12-26T12:00:00.000000",
                    "created_at": "2017-12-27 10:00:00",
                    "updated_at": null
                },
                {
                    "id": "277d6436-dfcb-4eae-ae5e-ac7fa9c2fd56",
                    "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                    "status": "UNDONE",
                    "event_type": "end_lease",
                    "time": "2017-12-27T12:00:00.000000",
                    "created_at": "2017-12-27 10:00:00",
                    "updated_at": null
                }
            ]
        }
    ]
}
`

const LeaseGetBody = `
{
    "lease": {
        "id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
        "name": "lease_foo",
        "start_date": "2017-12-26T12:00:00.000000",
        "end_date": "2017-12-27T12:00:00.000000",
        "status":"PENDING",
        "degraded": false,
        "user_id": "5434f637520d4c17bbf254af034b0320",
        "project_id": "aa45f56901ef45ee95e3d211097c0ea3",
        "trust_id": "b442a580b9504ababf305bf2b4c49512",
        "created_at": "2017-12-27 10:00:00",
        "updated_at": null,
        "reservations": [
            {
                "id": "087bc740-6d2d-410b-9d47-c7b2b55a9d36",
                "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                "status": "pending",
                "missing_resources": false,
                "resources_changed": false,
                "resource_id": "5e6c0e6e-f1e6-490b-baaf-50deacbbe371",
                "resource_type": "physical:host",
                "min": 4,
                "max": 6,
                "hypervisor_properties": "[\">=\", \"$vcpus\", \"4\"]",
                "resource_properties": "",
                "before_end": "default",
                "created_at": "2017-12-27 10:00:00",
                "updated_at": null
            },
            {
                "id": "ddc45423-f863-4e4e-8e7a-51d27cfec962",
                "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                "status": "pending",
                "missing_resources": false,
                "resources_changed": false,
                "resource_id": "0b901727-cca2-43ed-bcc8-c21b0982dcb1",
                "resource_type": "virtual:instance",
                "amount": 4,
                "vcpus": 2,
                "memory_mb": 4096,
                "disk_gb": 100,
                "affinity": false,
                "resource_properties": "",
                "flavor_id": "ddc45423-f863-4e4e-8e7a-51d27cfec962",
                "server_group_id": "33cdfc42-5a04-4fcc-b190-1abebaa056bb",
                "aggregate_id": 11,
                "created_at": "2017-12-27 10:00:00",
                "updated_at": null
            }
        ],
        "events": [
            {
                "id": "188a8584-f832-4df9-9a4a-51e6364420ff",
                "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                "status": "UNDONE",
                "event_type": "start_lease",
                "time": "2017-12-26T12:00:00.000000",
                "created_at": "2017-12-27 10:00:00",
                "updated_at": null
            },
            {
                "id": "277d6436-dfcb-4eae-ae5e-ac7fa9c2fd56",
                "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                "status": "UNDONE",
                "event_type": "end_lease",
                "time": "2017-12-27T12:00:00.000000",
                "created_at": "2017-12-27 10:00:00",
                "updated_at": null
            },
            {
                "id": "f583af71-ca21-4b66-87de-52211d118029",
                "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                "status": "UNDONE",
                "time": "2017-12-27T11:00:00.000000",
                "event_type": "before_end_lease",
                "created_at": "2017-12-27 10:00:00",
                "updated_at": null
            }
        ]
    }
}
`

const LeaseCreateBody = `
{
    "name": "lease_foo",
    "start_date": "2017-12-26 12:00",
    "end_date": "2017-12-27 12:00",
    "before_end_date": "2017-12-27 11:00",
    "reservations": [
        {
            "resource_type": "physical:host",
            "min": 4,
            "max": 6,
            "hypervisor_properties": "[\">=\", \"$vcpus\", \"4\"]",
            "resource_properties": "",
            "before_end": "default"
        },
        {
            "resource_type": "virtual:instance",
            "amount": 4,
            "vcpus": 2,
            "memory_mb": 4096,
            "disk_gb": 100,
            "affinity": false,
            "resource_properties": ""
        }
    ],
    "events": []
}
`

const LeaseCreateResponse = `
{
    "lease": {
        "id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
        "name": "lease_foo",
        "start_date": "2017-12-26T12:00:00.000000",
        "end_date": "2017-12-27T12:00:00.000000",
        "status":"PENDING",
        "degraded": false,
        "user_id": "5434f637520d4c17bbf254af034b0320",
        "project_id": "aa45f56901ef45ee95e3d211097c0ea3",
        "trust_id": "b442a580b9504ababf305bf2b4c49512",
        "created_at": "2017-12-27 10:00:00",
        "updated_at": null,
        "reservations": [
            {
                "id": "087bc740-6d2d-410b-9d47-c7b2b55a9d36",
                "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                "status": "pending",
                "missing_resources": false,
                "resources_changed": false,
                "resource_id": "5e6c0e6e-f1e6-490b-baaf-50deacbbe371",
                "resource_type": "physical:host",
                "min": 4,
                "max": 6,
                "hypervisor_properties": "[\">=\", \"$vcpus\", \"4\"]",
                "resource_properties": "",
                "before_end": "default",
                "created_at": "2017-12-27 10:00:00",
                "updated_at": null
            },
            {
                "id": "ddc45423-f863-4e4e-8e7a-51d27cfec962",
                "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                "status": "pending",
                "missing_resources": false,
                "resources_changed": false,
                "resource_id": "0b901727-cca2-43ed-bcc8-c21b0982dcb1",
                "resource_type": "virtual:instance",
                "amount": 4,
                "vcpus": 2,
                "memory_mb": 4096,
                "disk_gb": 100,
                "affinity": false,
                "resource_properties": "",
                "flavor_id": "ddc45423-f863-4e4e-8e7a-51d27cfec962",
                "server_group_id": "33cdfc42-5a04-4fcc-b190-1abebaa056bb",
                "aggregate_id": 11,
                "created_at": "2017-12-27 10:00:00",
                "updated_at": null
            }
        ],
        "events": [
            {
                "id": "188a8584-f832-4df9-9a4a-51e6364420ff",
                "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                "status": "UNDONE",
                "event_type": "start_lease",
                "time": "2017-12-26T12:00:00.000000",
                "created_at": "2017-12-27 10:00:00",
                "updated_at": null
            },
            {
                "id": "277d6436-dfcb-4eae-ae5e-ac7fa9c2fd56",
                "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                "status": "UNDONE",
                "event_type": "end_lease",
                "time": "2017-12-27T12:00:00.000000",
                "created_at": "2017-12-27 10:00:00",
                "updated_at": null
            },
            {
                "id": "f583af71-ca21-4b66-87de-52211d118029",
                "lease_id": "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
                "status": "UNDONE",
                "time": "2017-12-27T11:00:00.000000",
                "event_type": "before_end_lease",
                "created_at": "2017-12-27 10:00:00",
                "updated_at": null
            }
        ]
    }
}
`

func MockListResponse(t *testing.T) {
	th.Mux.HandleFunc("/leases", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprint(w, LeaseListBody)
		case "1":
			fmt.Fprint(w, `{"leases": []}`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

func MockGetResponse(t *testing.T) {
	th.Mux.HandleFunc("/leases/6ee55c78-ac52-41a6-99af-2d2d73bcc466", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, LeaseGetBody)
	})
}

func MockCreateResponse(t *testing.T) {
	th.Mux.HandleFunc("/leases", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, LeaseCreateBody)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprint(w, LeaseCreateResponse)
	})
}
