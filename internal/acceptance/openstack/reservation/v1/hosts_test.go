//go:build acceptance || reservation || hosts

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/hosts"
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

// Blazar removes an extra capability when it is updated to null. This requires
// the 2026.2 release or later.
func TestHostsRemoveCapability(t *testing.T) {
	clients.RequireAdmin(t)

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
