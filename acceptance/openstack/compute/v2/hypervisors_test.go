// +build acceptance compute hypervisors

package v2

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestHypervisorsList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	allPages, err := hypervisors.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	th.AssertNoErr(t, err)

	for _, h := range allHypervisors {
		tools.PrintResource(t, h)
	}
}

func TestHypervisorsGet(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	hypervisorID, err := getHypervisorID(t, client)
	th.AssertNoErr(t, err)

	hypervisor, err := hypervisors.Get(client, hypervisorID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, hypervisor)

	th.AssertEquals(t, hypervisorID, hypervisor.ID)
}

func TestHypervisorsGetStatistics(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	hypervisorsStats, err := hypervisors.GetStatistics(client).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, hypervisorsStats)

	if hypervisorsStats.Count == 0 {
		t.Fatalf("Unable to get hypervisor stats")
	}
}

func TestHypervisorsGetUptime(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	hypervisorID, err := getHypervisorID(t, client)
	th.AssertNoErr(t, err)

	hypervisor, err := hypervisors.GetUptime(client, hypervisorID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, hypervisor)

	th.AssertEquals(t, hypervisorID, hypervisor.ID)
}

func TestHypervisorListQuery(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	client.Microversion = "2.53"

	server, err := CreateMicroversionServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	iTrue := true
	listOpts := hypervisors.ListOpts{
		WithServers: &iTrue,
	}

	allPages, err := hypervisors.List(client, listOpts).AllPages()
	th.AssertNoErr(t, err)

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	th.AssertNoErr(t, err)

	hypervisor := allHypervisors[0]
	if len(*hypervisor.Servers) < 1 {
		t.Fatalf("hypervisor.Servers length should be >= 1")
	}
}

func getHypervisorID(t *testing.T, client *gophercloud.ServiceClient) (string, error) {
	allPages, err := hypervisors.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	th.AssertNoErr(t, err)

	if len(allHypervisors) > 0 {
		return allHypervisors[0].ID, nil
	}

	return "", fmt.Errorf("Unable to get hypervisor ID")
}
