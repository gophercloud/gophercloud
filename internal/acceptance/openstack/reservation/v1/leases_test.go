//go:build acceptance || reservation || leases

package v1

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/leases"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

var hostReservation = leases.HostReservationOpts{
	Min: 1,
	Max: 1,
}

func TestLeasesList(t *testing.T) {
	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	allPages, err := leases.List(client, leases.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allLeases, err := leases.ExtractLeases(allPages)
	th.AssertNoErr(t, err)

	for _, lease := range allLeases {
		tools.PrintResource(t, lease)
		tools.PrintResource(t, lease.StartDate)
		tools.PrintResource(t, lease.EndDate)
		tools.PrintResource(t, lease.CreatedAt)
		tools.PrintResource(t, lease.UpdatedAt)
		tools.PrintResource(t, lease.Reservations)
		tools.PrintResource(t, lease.Events)
	}
}

func TestLeasesCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	defer RequireFreepoolHost(t, client)()

	lease, err := CreateLease(t, client, tools.RandomString("ACPTTEST", 8), hostReservation)
	th.AssertNoErr(t, err)
	// Blazar refuses to remove a host that a lease still holds, so the lease
	// has to go before any host enrolled for it.
	defer DeleteLease(t, client, lease)

	tools.PrintResource(t, lease)
	tools.PrintResource(t, lease.StartDate)
	tools.PrintResource(t, lease.EndDate)

	newLease, err := leases.Get(context.TODO(), client, lease.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, lease.ID, newLease.ID)
	th.AssertEquals(t, lease.Name, newLease.Name)

	// Blazar schedules a start and an end event for every lease.
	th.AssertEquals(t, true, len(newLease.Events) >= 2)

	allPages, err := leases.List(client, leases.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allLeases, err := leases.ExtractLeases(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, l := range allLeases {
		if l.ID == lease.ID {
			found = true
		}
	}
	th.AssertEquals(t, true, found)

	// The lease has not started yet, so its start date can still move.
	newStartDate := lease.StartDate.Add(10 * time.Minute)
	newEndDate := lease.EndDate.Add(time.Hour)
	newBeforeEndDate := newEndDate.Add(-30 * time.Minute)

	updateOpts := leases.UpdateOpts{
		Name:          tools.RandomString("ACPTTEST", 8),
		StartDate:     &newStartDate,
		EndDate:       &newEndDate,
		BeforeEndDate: &newBeforeEndDate,
	}

	updatedLease, err := leases.Update(context.TODO(), client, lease.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, updatedLease)
	tools.PrintResource(t, updatedLease.StartDate)
	tools.PrintResource(t, updatedLease.EndDate)

	th.AssertEquals(t, updateOpts.Name, updatedLease.Name)
	th.AssertEquals(t, newStartDate, updatedLease.StartDate)
	th.AssertEquals(t, newEndDate, updatedLease.EndDate)

	// Blazar moves the events along with the dates they belong to.
	events := make(map[string]leases.Event, len(updatedLease.Events))
	for _, event := range updatedLease.Events {
		events[event.EventType] = event
	}

	th.AssertEquals(t, newStartDate, events["start_lease"].Time)
	th.AssertEquals(t, newEndDate, events["end_lease"].Time)
	th.AssertEquals(t, newBeforeEndDate, events["before_end_lease"].Time)
}

// A lease reserves the hosts it needs as soon as it is created, so the
// allocation is visible before the lease starts.
func TestLeasesHostReservation(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	defer RequireFreepoolHost(t, client)()

	lease, err := CreateLease(t, client, tools.RandomString("ACPTTEST", 8), hostReservation)
	th.AssertNoErr(t, err)
	defer DeleteLease(t, client, lease)

	reservation := lease.Reservations[0]
	tools.PrintResource(t, reservation)

	th.AssertEquals(t, leases.ResourceTypeHost, reservation.ResourceType)
	th.AssertEquals(t, lease.ID, reservation.LeaseID)
	th.AssertEquals(t, false, reservation.MissingResources)
	th.AssertEquals(t, 1, *reservation.Min)
	th.AssertEquals(t, 1, *reservation.Max)

	// The instance reservation fields are absent for a host reservation.
	th.AssertEquals(t, true, reservation.Amount == nil)
	th.AssertEquals(t, true, reservation.VCPUs == nil)
}

// Blazar creates the flavor, aggregate and server group backing an instance
// reservation as soon as the lease is created.
func TestLeasesInstanceReservation(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	defer RequireFreepoolHost(t, client)()

	lease, err := CreateLease(t, client, tools.RandomString("ACPTTEST", 8), leases.InstanceReservationOpts{
		Amount:   1,
		VCPUs:    1,
		MemoryMB: 512,
		DiskGB:   1,
		Affinity: gophercloud.Disabled,
	})
	th.AssertNoErr(t, err)
	defer DeleteLease(t, client, lease)

	reservation := lease.Reservations[0]
	tools.PrintResource(t, reservation)

	th.AssertEquals(t, leases.ResourceTypeInstance, reservation.ResourceType)
	th.AssertEquals(t, false, reservation.MissingResources)
	th.AssertEquals(t, 1, *reservation.Amount)
	th.AssertEquals(t, 1, *reservation.VCPUs)
	th.AssertEquals(t, 512, *reservation.MemoryMB)
	th.AssertEquals(t, 1, *reservation.DiskGB)

	// A requested anti-affinity comes back with the server group enforcing it.
	th.AssertEquals(t, false, *reservation.Affinity)
	th.AssertEquals(t, true, reservation.ServerGroupID != nil)

	// Blazar names the reserved flavor after the reservation.
	th.AssertEquals(t, reservation.ID, *reservation.FlavorID)
	th.AssertEquals(t, true, reservation.AggregateID != nil)

	// The host reservation fields are absent for an instance reservation.
	th.AssertEquals(t, true, reservation.Min == nil)
	th.AssertEquals(t, true, reservation.Max == nil)

	updateOpts := leases.UpdateOpts{
		Reservations: []leases.UpdateReservationOpts{
			{
				ID:       reservation.ID,
				MemoryMB: 1024,
			},
		},
	}

	_, err = leases.Update(context.TODO(), client, lease.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	updatedLease, err := leases.Get(context.TODO(), client, lease.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, updatedLease.Reservations[0])

	th.AssertEquals(t, 1024, *updatedLease.Reservations[0].MemoryMB)
}

func TestLeasesFlavorReservation(t *testing.T) {
	clients.RequireAdmin(t)
	clients.SkipReleasesBelow(t, "stable/2024.2")

	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	computeClient, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	flavor := SmallestFlavor(t, computeClient)

	defer RequireFreepoolHost(t, client)()

	lease, err := CreateLease(t, client, tools.RandomString("ACPTTEST", 8), leases.FlavorInstanceReservationOpts{
		Amount:   1,
		FlavorID: flavor.ID,
	})
	th.AssertNoErr(t, err)
	defer DeleteLease(t, client, lease)

	reservation := lease.Reservations[0]
	tools.PrintResource(t, reservation)

	th.AssertEquals(t, leases.ResourceTypeFlavorInstance, reservation.ResourceType)
	th.AssertEquals(t, 1, *reservation.Amount)
	th.AssertEquals(t, flavor.VCPUs, *reservation.VCPUs)
	th.AssertEquals(t, flavor.RAM, *reservation.MemoryMB)
	th.AssertEquals(t, flavor.Disk, *reservation.DiskGB)

	// Blazar does not support affinity for a flavor reservation, and creates a
	// flavor of its own named after the reservation.
	th.AssertEquals(t, true, reservation.Affinity == nil)
	th.AssertEquals(t, reservation.ID, *reservation.FlavorID)
}
