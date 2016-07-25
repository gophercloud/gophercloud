// +build acceptance compute servers

package v2

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/networks"
)

func TestNetworksList(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := networks.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	for _, network := range allNetworks {
		printNetwork(t, &network)
	}
}

func TestNetworksGet(t *testing.T) {
	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	networkID, err := getNetworkIDFromNetworks(t, client, choices.NetworkName)
	if err != nil {
		t.Fatal(err)
	}

	network, err := networks.Get(client, networkID).Extract()
	if err != nil {
		t.Fatalf("Unable to get network %s: %v", networkID, err)
	}

	printNetwork(t, network)
}

func getNetworkIDFromNetworkExtension(t *testing.T, client *gophercloud.ServiceClient, networkName string) (string, error) {
	allPages, err := networks.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	networkList, err := networks.ExtractNetworks(allPages)
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	networkID := ""
	for _, network := range networkList {
		t.Logf("Network: %v", network)
		if network.Label == networkName {
			networkID = network.ID
		}
	}

	t.Logf("Found network ID for %s: %s\n", networkName, networkID)

	return networkID, nil
}

func getNetworkIDFromNetworks(t *testing.T, client *gophercloud.ServiceClient, networkName string) (string, error) {
	allPages, err := networks.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	for _, network := range allNetworks {
		if network.Label == networkName {
			return network.ID, nil
		}
	}

	return "", fmt.Errorf("Failed to obtain network ID for network %s", networkName)
}

func printNetwork(t *testing.T, network *networks.Network) {
	t.Logf("Bridge: %s", network.Bridge)
	t.Logf("BridgeInterface: %s", network.BridgeInterface)
	t.Logf("Broadcast: %s", network.Broadcast)
	t.Logf("CIDR: %s", network.CIDR)
	t.Logf("CIDRv6: %s", network.CIDRv6)
	t.Logf("CreatedAt: %v", network.CreatedAt)
	t.Logf("Deleted: %t", network.Deleted)
	t.Logf("DeletedAt: %v", network.DeletedAt)
	t.Logf("DHCPStart: %s", network.DHCPStart)
	t.Logf("DNS1: %s", network.DNS1)
	t.Logf("DNS2: %s", network.DNS2)
	t.Logf("Gateway: %s", network.Gateway)
	t.Logf("Gatewayv6: %s", network.Gatewayv6)
	t.Logf("Host: %s", network.Host)
	t.Logf("ID: %s", network.ID)
	t.Logf("Injected: %t", network.Injected)
	t.Logf("Label: %s", network.Label)
	t.Logf("MultiHost: %t", network.MultiHost)
	t.Logf("Netmask: %s", network.Netmask)
	t.Logf("Netmaskv6: %s", network.Netmaskv6)
	t.Logf("Priority: %d", network.Priority)
	t.Logf("ProjectID: %s", network.ProjectID)
	t.Logf("RXTXBase: %d", network.RXTXBase)
	t.Logf("UpdatedAt: %v", network.UpdatedAt)
	t.Logf("VLAN: %d", network.VLAN)
	t.Logf("VPNPrivateAddress: %s", network.VPNPrivateAddress)
	t.Logf("VPNPublicAddress: %s", network.VPNPublicAddress)
	t.Logf("VPNPublicPort: %d", network.VPNPublicPort)
}
