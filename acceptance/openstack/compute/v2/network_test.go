// +build acceptance compute servers

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/networks"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestNetworksList(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	allPages, err := networks.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		t.Fatalf("Unable to list networks: %v", err)
	}

	var found bool
	for _, network := range allNetworks {
		tools.PrintResource(t, network)

		if network.Label == choices.NetworkName {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestNetworksGet(t *testing.T) {
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	networkID, err := GetNetworkIDFromNetworks(t, client, choices.NetworkName)
	if err != nil {
		t.Fatal(err)
	}

	network, err := networks.Get(client, networkID).Extract()
	if err != nil {
		t.Fatalf("Unable to get network %s: %v", networkID, err)
	}

	tools.PrintResource(t, network)

	th.AssertEquals(t, network.Label, choices.NetworkName)
}
