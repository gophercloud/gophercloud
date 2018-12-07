package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vlantransparent"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

// VLANTransparentNetwork represents OpenStack V2 Networking Network with the
// "vlan-transparent" extension enabled.
type VLANTransparentNetwork struct {
	networks.Network
	vlantransparent.TransparentExt
}

// ListVLANTransparentNetworks will list networks with the "vlan-transparent"
// extension. An error will be returned networks could not be listed.
func ListVLANTransparentNetworks(t *testing.T, client *gophercloud.ServiceClient) ([]*VLANTransparentNetwork, error) {
	iTrue := true
	networkListOpts := networks.ListOpts{}
	listOpts := vlantransparent.ListOptsExt{
		ListOptsBuilder: networkListOpts,
		VLANTransparent: &iTrue,
	}

	var allNetworks []*VLANTransparentNetwork

	t.Log("Attempting to list VLAN-transparent networks")

	allPages, err := networks.List(client, listOpts).AllPages()
	if err != nil {
		return nil, err
	}
	err = networks.ExtractNetworksInto(allPages, &allNetworks)
	if err != nil {
		return nil, err
	}

	t.Log("Successfully retrieved networks.")

	return allNetworks, nil
}
