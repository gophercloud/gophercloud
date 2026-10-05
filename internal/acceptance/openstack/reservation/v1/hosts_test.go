//go:build acceptance || reservation || hosts

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/hosts"
	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/leases"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestHostsList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	allPages, err := hosts.List(client, hosts.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allHosts, err := hosts.ExtractHosts(allPages)
	th.AssertNoErr(t, err)

	for _, host := range allHosts {
		tools.PrintResource(t, host)
		tools.PrintResource(t, host.CreatedAt)
		tools.PrintResource(t, host.UpdatedAt)
		tools.PrintResource(t, host.ExtraCapabilities)
	}
}

func TestHostsCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	computeClient, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	hypervisorHostname := HypervisorHostname(t, computeClient)

	host, err := CreateHost(t, client, hypervisorHostname)
	th.AssertNoErr(t, err)
	defer DeleteHost(t, client, host)

	tools.PrintResource(t, host)

	newHost, err := hosts.Get(context.TODO(), client, host.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, host.ID, newHost.ID)
	th.AssertEquals(t, "a100", newHost.ExtraCapabilities["gpu"])

	allPages, err := hosts.List(client, hosts.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allHosts, err := hosts.ExtractHosts(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, h := range allHosts {
		if h.ID == host.ID {
			found = true
		}
	}
	th.AssertEquals(t, true, found)

	updateOpts := hosts.UpdateOpts{
		ExtraCapabilities: map[string]any{"gpu": "h100"},
	}

	updatedHost, err := hosts.Update(context.TODO(), client, host.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, updatedHost)
	th.AssertEquals(t, "h100", updatedHost.ExtraCapabilities["gpu"])
}

// Blazar removes an extra capability when it is updated to null.
func TestHostsRemoveCapability(t *testing.T) {
	clients.RequireAdmin(t)
	clients.SkipReleasesBelow(t, "stable/2026.2")

	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	computeClient, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	host, err := CreateHost(t, client, HypervisorHostname(t, computeClient))
	th.AssertNoErr(t, err)
	defer DeleteHost(t, client, host)

	updateOpts := hosts.UpdateOpts{
		ExtraCapabilities: map[string]any{"gpu": nil},
	}

	updatedHost, err := hosts.Update(context.TODO(), client, host.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, updatedHost)

	_, ok := updatedHost.ExtraCapabilities["gpu"]
	th.AssertEquals(t, false, ok)
}

func TestHostsAllocations(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	defer RequireFreepoolHost(t, client)()

	lease, err := CreateLease(t, client, tools.RandomString("ACPTTEST", 8), leases.HostReservationOpts{
		Min: 1,
		Max: 1,
	})
	th.AssertNoErr(t, err)
	defer DeleteLease(t, client, lease)

	reservationID := lease.Reservations[0].ID

	allPages, err := hosts.ListAllocations(client, hosts.ListAllocationsOpts{LeaseID: lease.ID}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allAllocations, err := hosts.ExtractAllocations(allPages)
	th.AssertNoErr(t, err)

	// Every host is listed, but only the one reserved by the lease holds a
	// reservation.
	var hostID string
	for _, allocation := range allAllocations {
		tools.PrintResource(t, allocation)
		if len(allocation.Reservations) > 0 {
			hostID = allocation.ResourceID
		}
	}
	th.AssertEquals(t, true, hostID != "")

	allocation, err := hosts.GetAllocation(context.TODO(), client, hostID, hosts.GetAllocationOpts{ReservationID: reservationID}).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, allocation)
	tools.PrintResource(t, allocation.Reservations[0].StartDate)
	tools.PrintResource(t, allocation.Reservations[0].EndDate)

	th.AssertEquals(t, hostID, allocation.ResourceID)
	th.AssertEquals(t, 1, len(allocation.Reservations))
	th.AssertEquals(t, reservationID, allocation.Reservations[0].ID)
	th.AssertEquals(t, lease.ID, allocation.Reservations[0].LeaseID)
	th.AssertEquals(t, true, lease.StartDate.Equal(allocation.Reservations[0].StartDate))
	th.AssertEquals(t, true, lease.EndDate.Equal(allocation.Reservations[0].EndDate))
}

func TestHostsResourceProperties(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	defer RequireFreepoolHost(t, client)()

	allPages, err := hosts.ListResourceProperties(client, hosts.ListResourcePropertiesOpts{Detail: true, All: true}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allProperties, err := hosts.ExtractResourceProperties(allPages)
	th.AssertNoErr(t, err)

	for _, property := range allProperties {
		tools.PrintResource(t, property)
	}

	if len(allProperties) == 0 {
		t.Skip("no resource property to update")
	}

	name := allProperties[0].Property

	// Blazar reports every property as public even with Detail, so tell a
	// private one by its absence from the list without All.
	allPages, err = hosts.ListResourceProperties(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	publicProperties, err := hosts.ExtractResourceProperties(allPages)
	th.AssertNoErr(t, err)

	private := true
	for _, property := range publicProperties {
		if property.Property == name {
			private = false
		}
	}

	defer func() {
		_, err := hosts.UpdateResourceProperty(context.TODO(), client, name, hosts.UpdateResourcePropertyOpts{Private: private}).Extract()
		th.AssertNoErr(t, err)
	}()

	updated, err := hosts.UpdateResourceProperty(context.TODO(), client, name, hosts.UpdateResourcePropertyOpts{Private: !private}).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, updated)
	tools.PrintResource(t, updated.CreatedAt)
	tools.PrintResource(t, updated.UpdatedAt)

	th.AssertEquals(t, name, updated.PropertyName)
	th.AssertEquals(t, !private, updated.Private)
}
