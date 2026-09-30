// Package v1 contains common functions for creating reservation-based
// resources for use in acceptance tests. See the `*_test.go` files for example
// usages.
package v1

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/hypervisors"
	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/hosts"
	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/leases"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// HypervisorHostname returns the name of a hypervisor that can be enrolled
// into the freepool. Blazar knows a host by the name Nova gives it, so there
// is no way to enrol one without asking Nova first.
func HypervisorHostname(t *testing.T, client *gophercloud.ServiceClient) string {
	allPages, err := hypervisors.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	th.AssertNoErr(t, err)

	if len(allHypervisors) == 0 {
		t.Skip("no hypervisor available to enrol into the freepool")
	}

	return allHypervisors[0].HypervisorHostname
}

// RequireFreepoolHost makes sure the Blazar freepool holds a host to reserve.
// If it is empty, a hypervisor is enrolled and the returned function removes it.
func RequireFreepoolHost(t *testing.T, client *gophercloud.ServiceClient) func() {
	allPages, err := hosts.List(client, hosts.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allHosts, err := hosts.ExtractHosts(allPages)
	th.AssertNoErr(t, err)

	if len(allHosts) > 0 {
		t.Logf("Reusing freepool host: %s", allHosts[0].HypervisorHostname)
		return func() {}
	}

	computeClient, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	host, err := CreateHost(t, client, HypervisorHostname(t, computeClient))
	th.AssertNoErr(t, err)

	return func() { DeleteHost(t, client, host) }
}

// CreateHost will enrol a compute host into the Blazar freepool. An error will
// be returned if the host was unable to be created.
func CreateHost(t *testing.T, client *gophercloud.ServiceClient, name string) (*hosts.Host, error) {
	t.Logf("Attempting to create host: %s", name)

	createOpts := hosts.CreateOpts{
		Name: name,
		ExtraCapabilities: map[string]any{
			"gpu": "a100",
		},
	}

	host, err := hosts.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return host, err
	}

	t.Logf("Created host: %s", host.ID)

	th.AssertEquals(t, name, host.HypervisorHostname)
	th.AssertEquals(t, "a100", host.ExtraCapabilities["gpu"])

	return host, nil
}

// DeleteHost will remove a host from the Blazar freepool. A fatal error will
// occur if the host could not be deleted.
func DeleteHost(t *testing.T, client *gophercloud.ServiceClient, host *hosts.Host) {
	err := hosts.Delete(context.TODO(), client, host.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete host %s: %v", host.ID, err)
	}

	t.Logf("Deleted host: %s", host.ID)
}

// CreateLease will create a lease holding the given reservation, starting in a
// minute and lasting an hour. An error will be returned if the lease was
// unable to be created.
func CreateLease(t *testing.T, client *gophercloud.ServiceClient, name string, reservation leases.ReservationOptsBuilder) (*leases.Lease, error) {
	t.Logf("Attempting to create lease: %s", name)

	reservationMap, err := reservation.ToReservationMap()
	th.AssertNoErr(t, err)

	startDate := time.Now().UTC().Add(time.Minute)
	endDate := startDate.Add(time.Hour)

	createOpts := leases.CreateOpts{
		Name:         name,
		StartDate:    startDate,
		EndDate:      endDate,
		Reservations: []leases.ReservationOptsBuilder{reservation},
	}

	lease, err := leases.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return lease, err
	}

	t.Logf("Created lease: %s", lease.ID)

	th.AssertEquals(t, name, lease.Name)
	th.AssertEquals(t, 1, len(lease.Reservations))
	th.AssertEquals(t, reservationMap["resource_type"], lease.Reservations[0].ResourceType)

	return lease, nil
}

// SmallestFlavor returns the Nova flavor with the least memory.
func SmallestFlavor(t *testing.T, computeClient *gophercloud.ServiceClient) flavors.Flavor {
	allPages, err := flavors.ListDetail(computeClient, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allFlavors, err := flavors.ExtractFlavors(allPages)
	th.AssertNoErr(t, err)

	if len(allFlavors) == 0 {
		t.Skip("no flavor available to reserve instances of")
	}

	smallest := allFlavors[0]
	for _, flavor := range allFlavors[1:] {
		if flavor.RAM < smallest.RAM {
			smallest = flavor
		}
	}

	return smallest
}

// DeleteLease will delete a lease and release the resources it reserved.
// A fatal error will occur if the lease could not be deleted.
func DeleteLease(t *testing.T, client *gophercloud.ServiceClient, lease *leases.Lease) {
	err := leases.Delete(context.TODO(), client, lease.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete lease %s: %v", lease.ID, err)
	}

	t.Logf("Deleted lease: %s", lease.ID)
}
