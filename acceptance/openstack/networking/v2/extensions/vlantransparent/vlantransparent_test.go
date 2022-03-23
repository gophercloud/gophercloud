//go:build acceptance || networking || vlantransparent
// +build acceptance networking vlantransparent

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	networkingv2 "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/common/extensions"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestVLANTransparentCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	extension, err := extensions.Get(client, "vlan-transparent").Extract()
	if err != nil {
		t.Skip("This test requires vlan-transparent Neutron extension")
	}
	tools.PrintResource(t, extension)

	// Create a VLAN transparent network.
	network, err := CreateVLANTransparentNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networkingv2.DeleteNetwork(t, client, network.ID)

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
