package testing

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/leases"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

var (
	startDate     = time.Date(2026, 10, 1, 8, 0, 0, 0, time.UTC)
	endDate       = time.Date(2026, 10, 3, 18, 0, 0, 0, time.UTC)
	beforeEndDate = time.Date(2026, 10, 3, 16, 30, 0, 0, time.UTC)

	beforeEndDefault = "default"
	beforeEndNone    = ""
)

func TestListLeases(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListLeases(t, fakeServer)

	allPages, err := leases.List(client.ServiceClient(fakeServer), leases.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := leases.ExtractLeases(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedLeasesList, actual)
	th.AssertEquals(t, 1, len(actual[0].Reservations))
	th.AssertEquals(t, 3, len(actual[0].Events))
}

// Blazar returns an empty collection rather.
func TestListLeasesEmpty(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListLeasesEmpty(t, fakeServer)

	allPages, err := leases.List(client.ServiceClient(fakeServer), leases.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := leases.ExtractLeases(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(actual))
}

func TestCreateLease(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCreateLease(t, fakeServer)

	createOpts := leases.CreateOpts{
		Name:          "my_lease",
		StartDate:     startDate,
		EndDate:       endDate,
		BeforeEndDate: &beforeEndDate,
		Reservations: []leases.ReservationOptsBuilder{
			leases.HostReservationOpts{
				Min:                  1,
				Max:                  2,
				HypervisorProperties: `[">=", "$memory_mb", "8192"]`,
				BeforeEnd:            &beforeEndDefault,
			},
			leases.InstanceReservationOpts{
				Amount:             2,
				VCPUs:              1,
				MemoryMB:           2048,
				DiskGB:             20,
				Affinity:           gophercloud.Disabled,
				ResourceProperties: `["==", "$gpu", "a100"]`,
			},
		},
	}

	actual, err := leases.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExpectedCreatedLease, actual)

	// A host reservation carries the host fields and none of the instance ones.
	host := actual.Reservations[0]
	th.AssertEquals(t, leases.ResourceTypeHost, host.ResourceType)
	th.AssertEquals(t, "default", *host.BeforeEnd)
	th.AssertEquals(t, true, host.Amount == nil)
	th.AssertEquals(t, true, host.Affinity == nil)
	th.AssertEquals(t, true, host.FlavorID == nil)
	th.AssertEquals(t, true, host.AggregateID == nil)
	th.AssertEquals(t, false, actual.UpdatedAt == nil)
	th.AssertEquals(t, false, host.UpdatedAt == nil)

	byType := make(map[string]leases.Event, len(actual.Events))
	for _, event := range actual.Events {
		th.AssertEquals(t, "UNDONE", event.Status)
		th.AssertEquals(t, true, event.UpdatedAt == nil)
		byType[event.EventType] = event
	}

	th.AssertEquals(t, actual.StartDate, byType["start_lease"].Time)
	th.AssertEquals(t, actual.StartDate, byType["before_end_lease"].Time)
	th.AssertEquals(t, actual.EndDate, byType["end_lease"].Time)
}

func TestCreateFlavorLease(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCreateFlavorLease(t, fakeServer)

	createOpts := leases.CreateOpts{
		Name:      "lease_baz",
		StartDate: startDate,
		EndDate:   endDate,
		Reservations: []leases.ReservationOptsBuilder{
			leases.FlavorInstanceReservationOpts{
				Amount:   1,
				FlavorID: "1e1a9b1e-1f0a-4d1e-9f1a-0b1c2d3e4f5a",
			},
		},
	}

	actual, err := leases.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)

	reservation := actual.Reservations[0]
	th.AssertEquals(t, leases.ResourceTypeFlavorInstance, reservation.ResourceType)

	// Blazar copies the size of the source flavor onto the reservation.
	th.AssertEquals(t, 1, *reservation.Amount)
	th.AssertEquals(t, 1, *reservation.VCPUs)
	th.AssertEquals(t, 1875, *reservation.MemoryMB)
	th.AssertEquals(t, 10, *reservation.DiskGB)
	th.AssertEquals(t, 131, *reservation.AggregateID)

	// The source flavor is kept as JSON in resource_properties, while the
	// flavor to boot from is the one Blazar names after the reservation.
	var sourceFlavor struct {
		ID string `json:"id"`
	}
	th.AssertNoErr(t, json.Unmarshal([]byte(*reservation.ResourceProperties), &sourceFlavor))
	th.AssertEquals(t, "1e1a9b1e-1f0a-4d1e-9f1a-0b1c2d3e4f5a", sourceFlavor.ID)
	th.AssertEquals(t, reservation.ID, *reservation.FlavorID)

	// Blazar does not support affinity for a flavor reservation.
	th.AssertEquals(t, true, reservation.Affinity == nil)
	th.AssertEquals(t, true, reservation.ServerGroupID == nil)
}

// Blazar requires a name and both dates, so the request should not reach the
// server without them.
func TestCreateLeaseMissingRequired(t *testing.T) {
	reservations := []leases.ReservationOptsBuilder{
		leases.HostReservationOpts{Min: 1, Max: 1},
	}

	for name, opts := range map[string]leases.CreateOpts{
		"no name":         {StartDate: startDate, EndDate: endDate, Reservations: reservations},
		"no start date":   {Name: "my_lease", EndDate: endDate, Reservations: reservations},
		"no end date":     {Name: "my_lease", StartDate: startDate, Reservations: reservations},
		"no reservations": {Name: "my_lease", StartDate: startDate, EndDate: endDate},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := opts.ToLeaseCreateMap()
			th.AssertEquals(t, true, err != nil)
		})
	}
}

// A nil reservation builder must be reported rather than panic.
func TestCreateLeaseNilReservation(t *testing.T) {
	createOpts := leases.CreateOpts{
		Name:         "my_lease",
		StartDate:    startDate,
		EndDate:      endDate,
		Reservations: []leases.ReservationOptsBuilder{nil},
	}

	_, err := createOpts.ToLeaseCreateMap()
	th.AssertEquals(t, true, err != nil)
}

// Blazar requires an amount and a flavor for a flavor reservation.
func TestFlavorInstanceReservationOptsMissingFlavor(t *testing.T) {
	_, err := leases.FlavorInstanceReservationOpts{Amount: 4}.ToReservationMap()
	th.AssertEquals(t, true, err != nil)
}

// Blazar requires hypervisor_properties and resource_properties to be present,
// so an unset filter must be sent as an empty string rather than omitted.
func TestHostReservationOptsSendsEmptyFilters(t *testing.T) {
	b, err := leases.HostReservationOpts{Min: 1, Max: 1}.ToReservationMap()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "", b["hypervisor_properties"])
	th.AssertEquals(t, "", b["resource_properties"])
	th.AssertEquals(t, leases.ResourceTypeHost, b["resource_type"])
}

func TestHostReservationOptsBeforeEnd(t *testing.T) {
	b, err := leases.HostReservationOpts{Min: 1, Max: 1}.ToReservationMap()
	th.AssertNoErr(t, err)
	_, ok := b["before_end"]
	th.AssertEquals(t, false, ok)

	b, err = leases.HostReservationOpts{Min: 1, Max: 1, BeforeEnd: &beforeEndNone}.ToReservationMap()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "", b["before_end"])

	b, err = leases.HostReservationOpts{Min: 1, Max: 1, BeforeEnd: &beforeEndDefault}.ToReservationMap()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "default", b["before_end"])
}

// Blazar requires affinity to be present, so an unset placement must be sent
// as an explicit null rather than omitted.
func TestInstanceReservationOptsSendsNullAffinity(t *testing.T) {
	b, err := leases.InstanceReservationOpts{
		Amount:   4,
		VCPUs:    2,
		MemoryMB: 4096,
		DiskGB:   100,
	}.ToReservationMap()
	th.AssertNoErr(t, err)

	affinity, ok := b["affinity"]
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, true, affinity == nil)
}

func TestGetLease(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetLease(t, fakeServer)

	actual, err := leases.Get(context.TODO(), client.ServiceClient(fakeServer), "b179d3b5-6014-44a7-969b-0d333a969631").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExpectedLease, actual)
}

// Blazar returns an explicit null for the instance reservation fields that are
// unset, but omits the host-only fields entirely. Both must unmarshal to nil.
func TestGetLeaseNullFields(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetLease(t, fakeServer)

	lease, err := leases.Get(context.TODO(), client.ServiceClient(fakeServer), "b179d3b5-6014-44a7-969b-0d333a969631").Extract()
	th.AssertNoErr(t, err)

	reservation := lease.Reservations[0]

	// Explicit nulls.
	th.AssertEquals(t, true, reservation.Affinity == nil)
	th.AssertEquals(t, true, reservation.ServerGroupID == nil)

	// Omitted keys.
	th.AssertEquals(t, true, reservation.Min == nil)
	th.AssertEquals(t, true, reservation.Max == nil)
	th.AssertEquals(t, true, reservation.HypervisorProperties == nil)
	th.AssertEquals(t, true, reservation.BeforeEnd == nil)

	// The always-present booleans are not pointers, so a false must survive.
	th.AssertEquals(t, false, lease.Degraded)
	th.AssertEquals(t, false, reservation.MissingResources)
	th.AssertEquals(t, false, reservation.ResourcesChanged)

	// A modified lease reports both timestamp formats: the lease dates and the
	// event times carry microseconds and a T separator, the audit timestamps
	// carry neither.
	th.AssertEquals(t, time.Date(2026, 8, 24, 13, 49, 0, 0, time.UTC), lease.StartDate)
	th.AssertEquals(t, time.Date(2026, 8, 25, 12, 10, 0, 0, time.UTC), lease.EndDate)
	th.AssertEquals(t, time.Date(2026, 8, 24, 13, 48, 47, 0, time.UTC), lease.CreatedAt)
	th.AssertEquals(t, time.Date(2026, 8, 25, 12, 10, 9, 0, time.UTC), *lease.UpdatedAt)
	th.AssertEquals(t, time.Date(2026, 8, 24, 13, 48, 47, 0, time.UTC), reservation.CreatedAt)
	th.AssertEquals(t, time.Date(2026, 8, 25, 12, 10, 8, 0, time.UTC), *reservation.UpdatedAt)

	// Blazar does not order the events, so look them up by type rather than by
	// position.
	th.AssertEquals(t, 3, len(lease.Events))

	byType := make(map[string]leases.Event, len(lease.Events))
	for _, event := range lease.Events {
		byType[event.EventType] = event
	}

	for _, eventType := range []string{"start_lease", "end_lease", "before_end_lease"} {
		event, ok := byType[eventType]
		th.AssertEquals(t, true, ok)
		th.AssertEquals(t, "DONE", event.Status)
		th.AssertEquals(t, time.Date(2026, 8, 24, 13, 48, 48, 0, time.UTC), event.CreatedAt)
	}

	th.AssertEquals(t, time.Date(2026, 8, 24, 13, 49, 0, 0, time.UTC), byType["start_lease"].Time)
	th.AssertEquals(t, time.Date(2026, 8, 25, 12, 10, 0, 0, time.UTC), byType["end_lease"].Time)
	th.AssertEquals(t, time.Date(2026, 8, 25, 11, 10, 0, 0, time.UTC), byType["before_end_lease"].Time)
	th.AssertEquals(t, time.Date(2026, 8, 24, 13, 49, 6, 0, time.UTC), *byType["start_lease"].UpdatedAt)
}

func TestUpdateLease(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleUpdateLease(t, fakeServer)

	newStartDate := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	newEndDate := time.Date(2026, 10, 4, 18, 0, 0, 0, time.UTC)
	newBeforeEndDate := time.Date(2026, 10, 4, 17, 0, 0, 0, time.UTC)
	clearedProperties := ""

	updateOpts := leases.UpdateOpts{
		Name:          "lease_renamed",
		StartDate:     &newStartDate,
		EndDate:       &newEndDate,
		BeforeEndDate: &newBeforeEndDate,
		Reservations: []leases.UpdateReservationOpts{
			{
				ID:                   "b4675fff-ec59-480b-9399-9a74dffbc5c1",
				Min:                  2,
				Max:                  3,
				HypervisorProperties: &clearedProperties,
				ResourceProperties:   &clearedProperties,
			},
			{
				ID:       "57369989-d289-4ef0-a5a1-0ee932d4dcac",
				Amount:   3,
				VCPUs:    2,
				MemoryMB: 4096,
				DiskGB:   40,
				Affinity: gophercloud.Enabled,
			},
		},
	}

	actual, err := leases.Update(context.TODO(), client.ServiceClient(fakeServer), "b179d3b5-6014-44a7-969b-0d333a969631", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExpectedUpdatedLease, actual)
	th.AssertEquals(t, "UPDATING", actual.Status)

	byType := make(map[string]leases.Event, len(actual.Events))
	for _, event := range actual.Events {
		byType[event.EventType] = event
	}

	th.AssertEquals(t, actual.EndDate, byType["end_lease"].Time)
	th.AssertEquals(t, actual.EndDate.Add(-time.Hour), byType["before_end_lease"].Time)
	th.AssertEquals(t, false, byType["end_lease"].UpdatedAt == nil)
	th.AssertEquals(t, true, byType["start_lease"].UpdatedAt == nil)
}

// Blazar treats an empty body as a no-op, so the request should not be sent
// at all.
func TestUpdateLeaseWithoutOpts(t *testing.T) {
	_, err := leases.UpdateOpts{}.ToLeaseUpdateMap()
	th.AssertEquals(t, true, err != nil)
}

// Blazar looks a reservation up by ID when updating a lease.
func TestUpdateLeaseReservationWithoutID(t *testing.T) {
	updateOpts := leases.UpdateOpts{
		Reservations: []leases.UpdateReservationOpts{{Min: 6, Max: 8}},
	}

	_, err := updateOpts.ToLeaseUpdateMap()
	th.AssertEquals(t, true, err != nil)
}

func TestDeleteLease(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleDeleteLease(t, fakeServer)

	err := leases.Delete(context.TODO(), client.ServiceClient(fakeServer), "b179d3b5-6014-44a7-969b-0d333a969631").ExtractErr()
	th.AssertNoErr(t, err)
}
