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
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := hypervisors.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list hypervisors: %v", err)
	}

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	if err != nil {
		t.Fatalf("Unable to extract hypervisors")
	}

	for _, h := range allHypervisors {
		tools.PrintResource(t, h)
	}
}

func TestHypervisorsGet(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	hypervisorID, err := getHypervisorID(t, client)
	if err != nil {
		t.Fatal(err)
	}

	hypervisor, err := hypervisors.Get(client, hypervisorID).Extract()
	if err != nil {
		t.Fatalf("Unable to get hypervisor: %v", err)
	}

	tools.PrintResource(t, hypervisor)

	th.AssertEquals(t, hypervisorID, hypervisor.ID)
}

func TestHypervisorsGetStatistics(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	hypervisorsStats, err := hypervisors.GetStatistics(client).Extract()
	if err != nil {
		t.Fatalf("Unable to get hypervisors statistics: %v", err)
	}

	tools.PrintResource(t, hypervisorsStats)

	if hypervisorsStats.Count == 0 {
		t.Fatalf("Unable to get hypervisor stats")
	}
}

func TestHypervisorsGetUptime(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	hypervisorID, err := getHypervisorID(t, client)
	if err != nil {
		t.Fatal(err)
	}

	hypervisor, err := hypervisors.GetUptime(client, hypervisorID).Extract()
	if err != nil {
		t.Fatalf("Unable to hypervisor uptime: %v", err)
	}

	tools.PrintResource(t, hypervisor)

	th.AssertEquals(t, hypervisorID, hypervisor.ID)
}

func getHypervisorID(t *testing.T, client *gophercloud.ServiceClient) (int, error) {
	allPages, err := hypervisors.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list hypervisors: %v", err)
	}

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	if err != nil {
		t.Fatalf("Unable to extract hypervisors")
	}

	if len(allHypervisors) > 0 {
		return allHypervisors[0].ID, nil
	}

	return 0, fmt.Errorf("Unable to get hypervisor ID")
}
