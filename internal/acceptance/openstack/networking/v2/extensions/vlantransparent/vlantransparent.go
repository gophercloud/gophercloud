package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vlantransparent"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
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

	allPages, err := networks.List(client, listOpts).AllPages(context.TODO())
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

// CreateVLANTransparentNetwork will create a network with the
// "vlan-transparent" extension. An error will be returned if the network could
// not be created.
func CreateVLANTransparentNetwork(t *testing.T, client *gophercloud.ServiceClient) (*VLANTransparentNetwork, error) {
	networkName := tools.RandomString("TESTACC-", 8)
	networkCreateOpts := networks.CreateOpts{
		Name: networkName,
	}

	iTrue := true
	createOpts := vlantransparent.CreateOptsExt{
		CreateOptsBuilder: &networkCreateOpts,
		VLANTransparent:   &iTrue,
	}

	t.Logf("Attempting to create a VLAN-transparent network: %s", networkName)

	var network VLANTransparentNetwork
	err := networks.Create(context.TODO(), client, createOpts).ExtractInto(&network)
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created the network.")

	th.AssertEquals(t, networkName, network.Name)

	return &network, nil
}
