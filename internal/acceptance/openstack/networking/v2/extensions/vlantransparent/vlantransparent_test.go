//go:build acceptance || networking || vlantransparent

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestVLANTransparentCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "vlan-transparent")

	// Create a VLAN transparent network.
	network, err := CreateVLANTransparentNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	tools.PrintResource(t, network)

	// Update the created VLAN transparent network.
	newNetwork, err := UpdateVLANTransparentNetwork(t, client, network.ID)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newNetwork)

	// Check that the created VLAN transparent network exists.
	vlanTransparentNetworks, err := ListVLANTransparentNetworks(t, client)
	th.AssertNoErr(t, err)

	var found bool
	for _, vlanTransparentNetwork := range vlanTransparentNetworks {
		if vlanTransparentNetwork.ID == network.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}
