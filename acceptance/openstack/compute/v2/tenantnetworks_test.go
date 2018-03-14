// +build acceptance compute servers

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/tenantnetworks"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestTenantNetworksList(t *testing.T) {
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := tenantnetworks.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	allTenantNetworks, err := tenantnetworks.ExtractNetworks(allPages)
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	var found bool
	for _, network := range allTenantNetworks {
		tools.PrintResource(t, network)

		if network.Name == choices.NetworkName {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestTenantNetworksGet(t *testing.T) {
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	networkID, err := GetNetworkIDFromTenantNetworks(t, client, choices.NetworkName)
	if err != nil {
		t.Fatal(err)
	}

	network, err := tenantnetworks.Get(client, networkID).Extract()
	if err != nil {
		t.Fatalf("Unable to get network %s: %v", networkID, err)
	}

	tools.PrintResource(t, network)
}
