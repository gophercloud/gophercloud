package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/leases"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	allPages, err := leases.List(client.ServiceClient()).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := leases.ExtractLeases(allPages)
	th.AssertNoErr(t, err)

	expected := []leases.Lease{
		{
			ID:        "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
			Name:      "lease_foo",
			StartDate: time.Date(2017, 12, 26, 12, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2017, 12, 27, 12, 0, 0, 0, time.UTC),
			Status:    "PENDING",
			Degraded:  false,
			UserID:    "5434f637520d4c17bbf254af034b0320",
			ProjectID: "aa45f56901ef45ee95e3d211097c0ea3",
			TrustID:   "b442a580b9504ababf305bf2b4c49512",
			CreatedAt: time.Date(2017, 12, 27, 10, 0, 0, 0, time.UTC),
			Reservations: []leases.Reservation{
				{
					ID:                   "087bc740-6d2d-410b-9d47-c7b2b55a9d36",
					LeaseID:              "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
					Status:               "pending",
					MissingResources:     false,
					ResourcesChanged:     false,
					ResourceID:           "5e6c0e6e-f1e6-490b-baaf-50deacbbe371",
					ResourceType:         "physical:host",
					Min:                  4,
					Max:                  6,
					HypervisorProperties: "[\">=\", \"$vcpus\", \"4\"]",
					ResourceProperties:   "",
					BeforeEnd:            "default",
					CreatedAt:            time.Date(2017, 12, 27, 10, 0, 0, 0, time.UTC),
				},
				{
					ID:                 "ddc45423-f863-4e4e-8e7a-51d27cfec962",
					LeaseID:            "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
					Status:             "pending",
					MissingResources:   false,
					ResourcesChanged:   false,
					ResourceID:         "0b901727-cca2-43ed-bcc8-c21b0982dcb1",
					ResourceType:       "virtual:instance",
					Amount:             4,
					VCPUs:              2,
					MemoryMB:           4096,
					DiskGB:             100,
					Affinity:           false,
					ResourceProperties: "",
					FlavorID:           "ddc45423-f863-4e4e-8e7a-51d27cfec962",
					ServerGroupID:      "33cdfc42-5a04-4fcc-b190-1abebaa056bb",
					AggregateID:        11,
					CreatedAt:          time.Date(2017, 12, 27, 10, 0, 0, 0, time.UTC),
				},
			},
			Events: []leases.Event{
				{
					ID:        "188a8584-f832-4df9-9a4a-51e6364420ff",
					LeaseID:   "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
					Status:    "UNDONE",
					EventType: "start_lease",
					Time:      time.Date(2017, 12, 26, 12, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2017, 12, 27, 10, 0, 0, 0, time.UTC),
				},
				{
					ID:        "277d6436-dfcb-4eae-ae5e-ac7fa9c2fd56",
					LeaseID:   "6ee55c78-ac52-41a6-99af-2d2d73bcc466",
					Status:    "UNDONE",
					EventType: "end_lease",
					Time:      time.Date(2017, 12, 27, 12, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2017, 12, 27, 10, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	th.CheckDeepEquals(t, expected, actual)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	lease, err := leases.Get(context.TODO(), client.ServiceClient(), "6ee55c78-ac52-41a6-99af-2d2d73bcc466").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "lease_foo", lease.Name)
	th.AssertEquals(t, "6ee55c78-ac52-41a6-99af-2d2d73bcc466", lease.ID)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	lease, err := leases.Create(context.TODO(), client.ServiceClient(), leases.CreateOpts{
		Name:          "lease-001",
		StartDate:     "2017-12-26 12:00",
		EndDate:       "2025-05-15 12:00",
		BeforeEndDate: "2017-12-27 11:00",
		Reservations: []leases.Reservation{
			{
				ResourceType:         "physical:host",
				Min:                  4,
				Max:                  6,
				HypervisorProperties: "[\">=\", \"$vcpus\", \"4\"]",
				ResourceProperties:   "",
				BeforeEnd:            "default",
			},
			{
				ResourceType:       "virtual:instance",
				Amount:             4,
				VCPUs:              2,
				MemoryMB:           4096,
				DiskGB:             100,
				Affinity:           false,
				ResourceProperties: "",
			},
		},
		Events: []leases.Event{},
	}).Extract()
	t.Logf("%+v", lease)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "lease-001", lease.Name)
}
