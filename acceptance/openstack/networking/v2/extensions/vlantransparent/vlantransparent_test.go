// +build acceptance networking vlantransparent

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestVLANTransparentCRUD(t *testing.T) {
	t.Skip("We don't have VLAN transparent extension in OpenLab.")

	_, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a VLAN transparent network.
	// network, err := CreateVLANTransparentNetwork(t, client)
	// th.AssertNoErr(t, err)
	// defer DeleteNetwork(t, client, network.ID)

	// tools.PrintResource(t, network)

	// Update the created VLAN transparent network.
	// newNetwork, err := UpdateVLANTransparentNetwork(t, client, network.ID)
	// th.AssertNoErr(t, err)

	// tools.PrintResource(t, newNetwork)

	// Check that the created VLAN transparent network exists.
	// vlanTransparentNetworks, err := ListVLANTransparentNetworks(t, client)
	// th.AssertNoErr(t, err)

	// var found bool
	// for _, network := range vlanTransparentNetworks {
	// 	if network.ID == newNetwork.ID {
	// 		found = true
	// 	}
	// }

	// th.AssertEquals(t, found, true)
}
