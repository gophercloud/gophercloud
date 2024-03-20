//go:build acceptance || networking || vlantransparent

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networkingv2 "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/common/extensions"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestVLANTransparentCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	extension, err := extensions.Get(context.TODO(), client, "vlan-transparent").Extract()
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
