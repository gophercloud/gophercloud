// +build acceptance compute servers

package v2

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/tenantnetworks"
)

func TestTenantNetworksList(t *testing.T) {
	client, err := newClient()
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

	for _, network := range allTenantNetworks {
		printTenantNetwork(t, &network)
	}
}

func TestTenantNetworksGet(t *testing.T) {
	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	networkID, err := getNetworkIDFromTenantNetworks(t, client, choices.NetworkName)
	if err != nil {
		t.Fatal(err)
	}

	network, err := tenantnetworks.Get(client, networkID).Extract()
	if err != nil {
		t.Fatalf("Unable to get network %s: %v", networkID, err)
	}

	printTenantNetwork(t, network)
}

func getNetworkIDFromTenantNetworks(t *testing.T, client *gophercloud.ServiceClient, networkName string) (string, error) {
	allPages, err := tenantnetworks.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	allTenantNetworks, err := tenantnetworks.ExtractNetworks(allPages)
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	for _, network := range allTenantNetworks {
		if network.Name == networkName {
			return network.ID, nil
		}
	}

	return "", fmt.Errorf("Failed to obtain network ID for network %s", networkName)
}

func printTenantNetwork(t *testing.T, network *tenantnetworks.Network) {
	t.Logf("ID: %s", network.ID)
	t.Logf("Name: %s", network.Name)
	t.Logf("CIDR: %s", network.CIDR)
}
