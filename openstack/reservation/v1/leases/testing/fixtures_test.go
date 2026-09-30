package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/leases"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// LeaseBody is a terminated lease holding an instance reservation.
const LeaseBody = `
{
  "id": "b179d3b5-6014-44a7-969b-0d333a969631",
  "name": "my_lease",
  "start_date": "2026-08-24T13:49:00.000000",
  "end_date": "2026-08-25T12:10:00.000000",
  "status": "TERMINATED",
  "degraded": false,
  "user_id": "0a75eb6437d147a3bb7edf7aa4a18b79",
  "project_id": "a2b6c11209974d6c916a190697d95f5f",
  "trust_id": "46648d282e634131984e9277bbdd1720",
  "created_at": "2026-08-24 13:48:47",
  "updated_at": "2026-08-25 12:10:09",
  "reservations": [
    {
      "id": "39355156-7d5a-499a-87f8-35331a17b08e",
      "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
      "status": "deleted",
      "missing_resources": false,
      "resources_changed": false,
      "resource_id": "94a8b5c4-17d9-4724-82d9-521ffa63f747",
      "resource_type": "virtual:instance",
      "amount": 1,
      "vcpus": 2,
      "memory_mb": 3750,
      "disk_gb": 20,
      "affinity": null,
      "resource_properties": "",
      "flavor_id": "39355156-7d5a-499a-87f8-35331a17b08e",
      "server_group_id": null,
      "aggregate_id": 119,
      "created_at": "2026-08-24 13:48:47",
      "updated_at": "2026-08-25 12:10:08"
    }
  ],
  "events": [
    {
      "id": "41f904f5-35cf-43da-a588-6abaf67bc823",
      "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
      "status": "DONE",
      "event_type": "end_lease",
      "time": "2026-08-25T12:10:00.000000",
      "created_at": "2026-08-24 13:48:48",
      "updated_at": "2026-08-25 12:10:08"
    },
    {
      "id": "6959f3d6-06cf-4700-bebb-f13c4db7e16e",
      "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
      "status": "DONE",
      "event_type": "start_lease",
      "time": "2026-08-24T13:49:00.000000",
      "created_at": "2026-08-24 13:48:48",
      "updated_at": "2026-08-24 13:49:06"
    },
    {
      "id": "84377370-676e-4167-b4c5-fb9ac505ead1",
      "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
      "status": "DONE",
      "event_type": "before_end_lease",
      "time": "2026-08-25T11:10:00.000000",
      "created_at": "2026-08-24 13:48:48",
      "updated_at": "2026-08-25 11:10:04"
    }
  ]
}
`

var (
	LeasesListResult = fmt.Sprintf(`{"leases": [%s]}`, LeaseBody)
	LeaseGetResult   = fmt.Sprintf(`{"lease": %s}`, LeaseBody)
)

// LeaseCreateResult is a newly created lease reserving one host.
const LeaseCreateResult = `
{
  "lease": {
    "created_at": "2026-09-23 16:09:33",
    "degraded": false,
    "end_date": "2026-09-23T17:10:00.000000",
    "events": [
      {
        "created_at": "2026-09-23 16:09:34",
        "event_type": "start_lease",
        "id": "44fb8d43-8b91-46b1-8482-77e9d6d0b398",
        "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
        "status": "UNDONE",
        "time": "2026-09-23T16:10:00.000000",
        "updated_at": null
      },
      {
        "created_at": "2026-09-23 16:09:34",
        "event_type": "before_end_lease",
        "id": "951847da-aa35-4c4d-9b04-f6ebadb2699d",
        "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
        "status": "UNDONE",
        "time": "2026-09-23T16:10:00.000000",
        "updated_at": null
      },
      {
        "created_at": "2026-09-23 16:09:34",
        "event_type": "end_lease",
        "id": "d2d84160-37a1-41a3-b5ef-3872c82e81d0",
        "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
        "status": "UNDONE",
        "time": "2026-09-23T17:10:00.000000",
        "updated_at": null
      }
    ],
    "id": "b179d3b5-6014-44a7-969b-0d333a969631",
    "name": "my_lease",
    "project_id": "a2b6c11209974d6c916a190697d95f5f",
    "reservations": [
      {
        "before_end": "default",
        "created_at": "2026-09-23 16:09:33",
        "hypervisor_properties": "",
        "id": "b4675fff-ec59-480b-9399-9a74dffbc5c1",
        "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
        "max": 1,
        "min": 1,
        "missing_resources": false,
        "resource_id": "24d2442a-1b65-4ed7-973e-e22fed5fc60f",
        "resource_properties": "",
        "resource_type": "physical:host",
        "resources_changed": false,
        "status": "pending",
        "updated_at": "2026-09-23 16:09:34"
      }
    ],
    "start_date": "2026-09-23T16:10:00.000000",
    "status": "PENDING",
    "trust_id": "46648d282e634131984e9277bbdd1720",
    "updated_at": "2026-09-23 16:09:34",
    "user_id": "0a75eb6437d147a3bb7edf7aa4a18b79"
  }
}
`

const LeaseCreateRequest = `
{
  "name": "my_lease",
  "start_date": "2026-10-01 08:00",
  "end_date": "2026-10-03 18:00",
  "before_end_date": "2026-10-03 16:30",
  "reservations": [
    {
      "resource_type": "physical:host",
      "min": 1,
      "max": 2,
      "hypervisor_properties": "[\">=\", \"$memory_mb\", \"8192\"]",
      "resource_properties": "",
      "before_end": "default"
    },
    {
      "resource_type": "virtual:instance",
      "amount": 2,
      "vcpus": 1,
      "memory_mb": 2048,
      "disk_gb": 20,
      "affinity": false,
      "resource_properties": "[\"==\", \"$gpu\", \"a100\"]"
    }
  ]
}
`

// LeaseCreateFlavorResult is a newly created lease reserving one instance of
// an existing flavor.
const LeaseCreateFlavorResult = `
{
  "lease": {
    "created_at": "2026-09-25 09:44:18",
    "degraded": false,
    "end_date": "2026-09-25T10:45:00.000000",
    "events": [
      {
        "created_at": "2026-09-25 09:44:18",
        "event_type": "end_lease",
        "id": "02f1f1c2-96ba-4626-b8b7-b86ed1981065",
        "lease_id": "8b0e2a6c-3f5d-4e1a-9c7b-2d4f6a8e0c13",
        "status": "UNDONE",
        "time": "2026-09-25T10:45:00.000000",
        "updated_at": null
      },
      {
        "created_at": "2026-09-25 09:44:18",
        "event_type": "start_lease",
        "id": "03c3f656-6252-458c-bbb4-7c2502008d08",
        "lease_id": "8b0e2a6c-3f5d-4e1a-9c7b-2d4f6a8e0c13",
        "status": "UNDONE",
        "time": "2026-09-25T09:45:00.000000",
        "updated_at": null
      },
      {
        "created_at": "2026-09-25 09:44:18",
        "event_type": "before_end_lease",
        "id": "2ffd2718-98d4-4239-a06e-46771542d595",
        "lease_id": "8b0e2a6c-3f5d-4e1a-9c7b-2d4f6a8e0c13",
        "status": "UNDONE",
        "time": "2026-09-25T09:45:00.000000",
        "updated_at": null
      }
    ],
    "id": "8b0e2a6c-3f5d-4e1a-9c7b-2d4f6a8e0c13",
    "name": "lease_baz",
    "project_id": "a2b6c11209974d6c916a190697d95f5f",
    "reservations": [
      {
        "affinity": null,
        "aggregate_id": 131,
        "amount": 1,
        "created_at": "2026-09-25 09:44:18",
        "disk_gb": 10,
        "flavor_id": "e0d58000-686d-4727-96b0-3e662bd2f7a0",
        "id": "e0d58000-686d-4727-96b0-3e662bd2f7a0",
        "lease_id": "8b0e2a6c-3f5d-4e1a-9c7b-2d4f6a8e0c13",
        "memory_mb": 1875,
        "missing_resources": false,
        "resource_id": "655ffd58-0d14-4f46-b8f5-57171e4f82c7",
        "resource_properties": "{\"id\": \"1e1a9b1e-1f0a-4d1e-9f1a-0b1c2d3e4f5a\", \"name\": \"flavor_foo\", \"ram\": 1875, \"disk\": 10, \"swap\": \"\", \"OS-FLV-EXT-DATA:ephemeral\": 0, \"OS-FLV-DISABLED:disabled\": false, \"vcpus\": 1, \"os-flavor-access:is_public\": false, \"rxtx_factor\": 1.0, \"extra_specs\": {}}",
        "resource_type": "flavor:instance",
        "resources_changed": false,
        "server_group_id": null,
        "status": "pending",
        "updated_at": "2026-09-25 09:44:18",
        "vcpus": 1
      }
    ],
    "start_date": "2026-09-25T09:45:00.000000",
    "status": "PENDING",
    "trust_id": "46648d282e634131984e9277bbdd1720",
    "updated_at": "2026-09-25 09:44:18",
    "user_id": "0a75eb6437d147a3bb7edf7aa4a18b79"
  }
}
`

// Blazar rejects any affinity other than null for a flavor reservation, so the
// key is always sent as an explicit null.
const LeaseCreateFlavorRequest = `
{
  "name": "lease_baz",
  "start_date": "2026-10-01 08:00",
  "end_date": "2026-10-03 18:00",
  "reservations": [
    {
      "resource_type": "flavor:instance",
      "amount": 1,
      "flavor_id": "1e1a9b1e-1f0a-4d1e-9f1a-0b1c2d3e4f5a",
      "affinity": null
    }
  ]
}
`

const LeaseUpdateRequest = `
{
  "name": "lease_renamed",
  "start_date": "2026-10-01 09:00",
  "end_date": "2026-10-04 18:00",
  "before_end_date": "2026-10-04 17:00",
  "reservations": [
    {
      "id": "b4675fff-ec59-480b-9399-9a74dffbc5c1",
      "min": 2,
      "max": 3,
      "hypervisor_properties": "",
      "resource_properties": ""
    },
    {
      "id": "57369989-d289-4ef0-a5a1-0ee932d4dcac",
      "amount": 3,
      "vcpus": 2,
      "memory_mb": 4096,
      "disk_gb": 40,
      "affinity": true
    }
  ]
}
`

// LeaseUpdateResult is the lease of LeaseCreateResult after renaming it and
// moving its end date.
const LeaseUpdateResult = `
{
  "lease": {
    "created_at": "2026-09-23 16:09:33",
    "degraded": false,
    "end_date": "2026-09-23T18:10:00.000000",
    "events": [
      {
        "created_at": "2026-09-23 16:09:34",
        "event_type": "start_lease",
        "id": "44fb8d43-8b91-46b1-8482-77e9d6d0b398",
        "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
        "status": "UNDONE",
        "time": "2026-09-23T16:10:00.000000",
        "updated_at": null
      },
      {
        "created_at": "2026-09-23 16:09:34",
        "event_type": "before_end_lease",
        "id": "951847da-aa35-4c4d-9b04-f6ebadb2699d",
        "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
        "status": "UNDONE",
        "time": "2026-09-23T17:10:00.000000",
        "updated_at": "2026-09-23 16:09:36"
      },
      {
        "created_at": "2026-09-23 16:09:34",
        "event_type": "end_lease",
        "id": "d2d84160-37a1-41a3-b5ef-3872c82e81d0",
        "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
        "status": "UNDONE",
        "time": "2026-09-23T18:10:00.000000",
        "updated_at": "2026-09-23 16:09:36"
      }
    ],
    "id": "b179d3b5-6014-44a7-969b-0d333a969631",
    "name": "lease_renamed",
    "project_id": "a2b6c11209974d6c916a190697d95f5f",
    "reservations": [
      {
        "before_end": "default",
        "created_at": "2026-09-23 16:09:33",
        "hypervisor_properties": "",
        "id": "b4675fff-ec59-480b-9399-9a74dffbc5c1",
        "lease_id": "b179d3b5-6014-44a7-969b-0d333a969631",
        "max": 1,
        "min": 1,
        "missing_resources": false,
        "resource_id": "24d2442a-1b65-4ed7-973e-e22fed5fc60f",
        "resource_properties": "",
        "resource_type": "physical:host",
        "resources_changed": false,
        "status": "pending",
        "updated_at": "2026-09-23 16:09:34"
      }
    ],
    "start_date": "2026-09-23T16:10:00.000000",
    "status": "UPDATING",
    "trust_id": "46648d282e634131984e9277bbdd1720",
    "updated_at": "2026-09-23 16:09:36",
    "user_id": "0a75eb6437d147a3bb7edf7aa4a18b79"
  }
}
`

// Pointer targets for the expected fixtures.
var (
	liveResourceProperties = ""
	liveFlavorID           = "39355156-7d5a-499a-87f8-35331a17b08e"

	createdHostProperties = ""
	createdBeforeEnd      = "default"
	createdUpdatedAt      = time.Date(2026, 9, 23, 16, 9, 34, 0, time.UTC)
	renamedUpdatedAt      = time.Date(2026, 9, 23, 16, 9, 36, 0, time.UTC)

	liveUpdatedAt            = time.Date(2026, 8, 25, 12, 10, 9, 0, time.UTC)
	liveReservationUpdatedAt = time.Date(2026, 8, 25, 12, 10, 8, 0, time.UTC)
	liveEndEventUpdatedAt    = time.Date(2026, 8, 25, 12, 10, 8, 0, time.UTC)
	liveStartEventUpdatedAt  = time.Date(2026, 8, 24, 13, 49, 6, 0, time.UTC)
	liveBeforeEndUpdatedAt   = time.Date(2026, 8, 25, 11, 10, 4, 0, time.UTC)
)

var ExpectedLease = leases.Lease{
	ID:        "b179d3b5-6014-44a7-969b-0d333a969631",
	Name:      "my_lease",
	StartDate: time.Date(2026, 8, 24, 13, 49, 0, 0, time.UTC),
	EndDate:   time.Date(2026, 8, 25, 12, 10, 0, 0, time.UTC),
	Status:    "TERMINATED",
	Degraded:  false,
	UserID:    "0a75eb6437d147a3bb7edf7aa4a18b79",
	ProjectID: "a2b6c11209974d6c916a190697d95f5f",
	TrustID:   "46648d282e634131984e9277bbdd1720",
	CreatedAt: time.Date(2026, 8, 24, 13, 48, 47, 0, time.UTC),
	UpdatedAt: &liveUpdatedAt,
	Reservations: []leases.Reservation{
		{
			ID:                 "39355156-7d5a-499a-87f8-35331a17b08e",
			LeaseID:            "b179d3b5-6014-44a7-969b-0d333a969631",
			Status:             "deleted",
			MissingResources:   false,
			ResourcesChanged:   false,
			ResourceID:         "94a8b5c4-17d9-4724-82d9-521ffa63f747",
			ResourceType:       leases.ResourceTypeInstance,
			Amount:             gophercloud.IntToPointer(1),
			VCPUs:              gophercloud.IntToPointer(2),
			MemoryMB:           gophercloud.IntToPointer(3750),
			DiskGB:             gophercloud.IntToPointer(20),
			Affinity:           nil,
			ResourceProperties: &liveResourceProperties,
			FlavorID:           &liveFlavorID,
			ServerGroupID:      nil,
			AggregateID:        gophercloud.IntToPointer(119),
			CreatedAt:          time.Date(2026, 8, 24, 13, 48, 47, 0, time.UTC),
			UpdatedAt:          &liveReservationUpdatedAt,
		},
	},
	Events: []leases.Event{
		{
			ID:        "41f904f5-35cf-43da-a588-6abaf67bc823",
			LeaseID:   "b179d3b5-6014-44a7-969b-0d333a969631",
			Status:    "DONE",
			EventType: "end_lease",
			Time:      time.Date(2026, 8, 25, 12, 10, 0, 0, time.UTC),
			CreatedAt: time.Date(2026, 8, 24, 13, 48, 48, 0, time.UTC),
			UpdatedAt: &liveEndEventUpdatedAt,
		},
		{
			ID:        "6959f3d6-06cf-4700-bebb-f13c4db7e16e",
			LeaseID:   "b179d3b5-6014-44a7-969b-0d333a969631",
			Status:    "DONE",
			EventType: "start_lease",
			Time:      time.Date(2026, 8, 24, 13, 49, 0, 0, time.UTC),
			CreatedAt: time.Date(2026, 8, 24, 13, 48, 48, 0, time.UTC),
			UpdatedAt: &liveStartEventUpdatedAt,
		},
		{
			ID:        "84377370-676e-4167-b4c5-fb9ac505ead1",
			LeaseID:   "b179d3b5-6014-44a7-969b-0d333a969631",
			Status:    "DONE",
			EventType: "before_end_lease",
			Time:      time.Date(2026, 8, 25, 11, 10, 0, 0, time.UTC),
			CreatedAt: time.Date(2026, 8, 24, 13, 48, 48, 0, time.UTC),
			UpdatedAt: &liveBeforeEndUpdatedAt,
		},
	},
}

var ExpectedLeasesList = []leases.Lease{ExpectedLease}

var ExpectedCreatedLease = leases.Lease{
	ID:        "b179d3b5-6014-44a7-969b-0d333a969631",
	Name:      "my_lease",
	StartDate: time.Date(2026, 9, 23, 16, 10, 0, 0, time.UTC),
	EndDate:   time.Date(2026, 9, 23, 17, 10, 0, 0, time.UTC),
	Status:    "PENDING",
	Degraded:  false,
	UserID:    "0a75eb6437d147a3bb7edf7aa4a18b79",
	ProjectID: "a2b6c11209974d6c916a190697d95f5f",
	TrustID:   "46648d282e634131984e9277bbdd1720",
	CreatedAt: time.Date(2026, 9, 23, 16, 9, 33, 0, time.UTC),
	UpdatedAt: &createdUpdatedAt,
	Reservations: []leases.Reservation{
		{
			ID:                   "b4675fff-ec59-480b-9399-9a74dffbc5c1",
			LeaseID:              "b179d3b5-6014-44a7-969b-0d333a969631",
			Status:               "pending",
			MissingResources:     false,
			ResourcesChanged:     false,
			ResourceID:           "24d2442a-1b65-4ed7-973e-e22fed5fc60f",
			ResourceType:         leases.ResourceTypeHost,
			Min:                  gophercloud.IntToPointer(1),
			Max:                  gophercloud.IntToPointer(1),
			HypervisorProperties: &createdHostProperties,
			ResourceProperties:   &createdHostProperties,
			BeforeEnd:            &createdBeforeEnd,
			CreatedAt:            time.Date(2026, 9, 23, 16, 9, 33, 0, time.UTC),
			UpdatedAt:            &createdUpdatedAt,
		},
	},
	Events: []leases.Event{
		{
			ID:        "44fb8d43-8b91-46b1-8482-77e9d6d0b398",
			LeaseID:   "b179d3b5-6014-44a7-969b-0d333a969631",
			Status:    "UNDONE",
			EventType: "start_lease",
			Time:      time.Date(2026, 9, 23, 16, 10, 0, 0, time.UTC),
			CreatedAt: time.Date(2026, 9, 23, 16, 9, 34, 0, time.UTC),
		},
		{
			ID:        "951847da-aa35-4c4d-9b04-f6ebadb2699d",
			LeaseID:   "b179d3b5-6014-44a7-969b-0d333a969631",
			Status:    "UNDONE",
			EventType: "before_end_lease",
			Time:      time.Date(2026, 9, 23, 16, 10, 0, 0, time.UTC),
			CreatedAt: time.Date(2026, 9, 23, 16, 9, 34, 0, time.UTC),
		},
		{
			ID:        "d2d84160-37a1-41a3-b5ef-3872c82e81d0",
			LeaseID:   "b179d3b5-6014-44a7-969b-0d333a969631",
			Status:    "UNDONE",
			EventType: "end_lease",
			Time:      time.Date(2026, 9, 23, 17, 10, 0, 0, time.UTC),
			CreatedAt: time.Date(2026, 9, 23, 16, 9, 34, 0, time.UTC),
		},
	},
}

// ExpectedUpdatedLease is ExpectedCreatedLease with the changes Blazar reports
// in LeaseUpdateResult.
var ExpectedUpdatedLease = func() leases.Lease {
	l := ExpectedCreatedLease
	l.Name = "lease_renamed"
	l.EndDate = time.Date(2026, 9, 23, 18, 10, 0, 0, time.UTC)
	l.Status = "UPDATING"
	l.UpdatedAt = &renamedUpdatedAt

	l.Events = append([]leases.Event(nil), ExpectedCreatedLease.Events...)
	l.Events[1].Time = time.Date(2026, 9, 23, 17, 10, 0, 0, time.UTC)
	l.Events[1].UpdatedAt = &renamedUpdatedAt
	l.Events[2].Time = time.Date(2026, 9, 23, 18, 10, 0, 0, time.UTC)
	l.Events[2].UpdatedAt = &renamedUpdatedAt

	return l
}()

func HandleListLeases(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/leases",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, LeasesListResult)
		})
}

func HandleListLeasesEmpty(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/leases",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, `{"leases": []}`)
		})
}

func HandleCreateLease(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/leases",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, LeaseCreateRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			fmt.Fprint(w, LeaseCreateResult)
		})
}

func HandleCreateFlavorLease(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/leases",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, LeaseCreateFlavorRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			fmt.Fprint(w, LeaseCreateFlavorResult)
		})
}

func HandleGetLease(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/leases/b179d3b5-6014-44a7-969b-0d333a969631",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, LeaseGetResult)
		})
}

func HandleUpdateLease(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/leases/b179d3b5-6014-44a7-969b-0d333a969631",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, LeaseUpdateRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, LeaseUpdateResult)
		})
}

func HandleDeleteLease(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/leases/b179d3b5-6014-44a7-969b-0d333a969631",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNoContent)
		})
}
