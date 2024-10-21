//go:build acceptance || networking || bgp || bgpvpns

package bgpvpns

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2/extensions/layer3"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgpvpns"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestBGPVPNCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a BGP VPN
	bgpVpnCreated, err := CreateBGPVPN(t, client)
	th.AssertNoErr(t, err)

	// Get a BGP VPN
	bgpVpnGot, err := bgpvpns.Get(context.TODO(), client, bgpVpnCreated.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bgpVpnCreated.ID, bgpVpnGot.ID)
	th.AssertEquals(t, bgpVpnCreated.Name, bgpVpnGot.Name)

	// Update a BGP VPN
	newBGPVPNName := tools.RandomString("TESTACC-BGPVPN-", 10)
	updateBGPOpts := bgpvpns.UpdateOpts{
		Name: &newBGPVPNName,
	}
	bgpVpnUpdated, err := bgpvpns.Update(context.TODO(), client, bgpVpnGot.ID, updateBGPOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newBGPVPNName, bgpVpnUpdated.Name)
	t.Logf("Update BGP VPN, renamed from %s to %s", bgpVpnGot.Name, bgpVpnUpdated.Name)

	// List all BGP VPNs
	allPages, err := bgpvpns.List(client, bgpvpns.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allVPNs, err := bgpvpns.ExtractBGPVPNs(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved BGP VPNs")
	tools.PrintResource(t, allVPNs)
	th.AssertIntGreaterOrEqual(t, len(allVPNs), 1)

	// Delete a BGP VPN
	t.Logf("Attempting to delete BGP VPN: %s", bgpVpnUpdated.Name)
	err = bgpvpns.Delete(context.TODO(), client, bgpVpnUpdated.ID).ExtractErr()
	th.AssertNoErr(t, err)

	_, err = bgpvpns.Get(context.TODO(), client, bgpVpnGot.ID).Extract()
	th.AssertErr(t, err)
	t.Logf("BGP VPN %s deleted", bgpVpnUpdated.Name)
}

func TestBGPVPNNetworkAssociationCRD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a BGP VPN
	bgpVpnCreated, err := CreateBGPVPN(t, client)
	th.AssertNoErr(t, err)
	defer func() {
		err = bgpvpns.Delete(context.TODO(), client, bgpVpnCreated.ID).ExtractErr()
		th.AssertNoErr(t, err)
	}()

	// Create a Network
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	// Associate a Network with a BGP VPN
	assocOpts := bgpvpns.CreateNetworkAssociationOpts{
		NetworkID: network.ID,
	}
	assoc, err := bgpvpns.CreateNetworkAssociation(context.TODO(), client, bgpVpnCreated.ID, assocOpts).Extract()
	th.AssertNoErr(t, err)
	defer func() {
		err = bgpvpns.DeleteNetworkAssociation(context.TODO(), client, bgpVpnCreated.ID, assoc.ID).ExtractErr()
		th.AssertNoErr(t, err)
	}()
	th.AssertEquals(t, network.ID, assoc.NetworkID)

	// Get a Network Association
	assocGot, err := bgpvpns.GetNetworkAssociation(context.TODO(), client, bgpVpnCreated.ID, assoc.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, network.ID, assocGot.NetworkID)

	// List all Network Associations
	allPages, err := bgpvpns.ListNetworkAssociations(client, bgpVpnCreated.ID, bgpvpns.ListNetworkAssociationsOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allAssocs, err := bgpvpns.ExtractNetworkAssociations(allPages)
	th.AssertNoErr(t, err)
	t.Logf("Retrieved Network Associations")
	tools.PrintResource(t, allAssocs)
	th.AssertIntGreaterOrEqual(t, len(allAssocs), 1)

	// Get BGP VPN with associations
	getBgpVpn, err := bgpvpns.Get(context.TODO(), client, bgpVpnCreated.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getBgpVpn)
}

func TestBGPVPNRouterAssociationCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a BGP VPN
	bgpVpnCreated, err := CreateBGPVPN(t, client)
	th.AssertNoErr(t, err)
	defer func() {
		err = bgpvpns.Delete(context.TODO(), client, bgpVpnCreated.ID).ExtractErr()
		th.AssertNoErr(t, err)
	}()

	// Create a Network
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	// Create a Router
	routerCreated, err := layer3.CreateRouter(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer layer3.DeleteRouter(t, client, routerCreated.ID)

	// Associate a Router with a BGP VPN
	assocOpts := bgpvpns.CreateRouterAssociationOpts{
		RouterID: routerCreated.ID,
	}
	assoc, err := bgpvpns.CreateRouterAssociation(context.TODO(), client, bgpVpnCreated.ID, assocOpts).Extract()
	th.AssertNoErr(t, err)
	defer func() {
		err = bgpvpns.DeleteRouterAssociation(context.TODO(), client, bgpVpnCreated.ID, assoc.ID).ExtractErr()
		th.AssertNoErr(t, err)
	}()
	th.AssertEquals(t, routerCreated.ID, assoc.RouterID)

	// Get a Router Association
	assocGot, err := bgpvpns.GetRouterAssociation(context.TODO(), client, bgpVpnCreated.ID, assoc.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, routerCreated.ID, assocGot.RouterID)

	// Update a Router Association
	assocUpdOpts := bgpvpns.UpdateRouterAssociationOpts{
		AdvertiseExtraRoutes: new(bool),
	}
	assocUpdate, err := bgpvpns.UpdateRouterAssociation(context.TODO(), client, bgpVpnCreated.ID, assoc.ID, assocUpdOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, routerCreated.ID, assocUpdate.RouterID)
	th.AssertEquals(t, false, assocUpdate.AdvertiseExtraRoutes)

	// List all Router Associations
	allPages, err := bgpvpns.ListRouterAssociations(client, bgpVpnCreated.ID, bgpvpns.ListRouterAssociationsOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allAssocs, err := bgpvpns.ExtractRouterAssociations(allPages)
	th.AssertNoErr(t, err)
	t.Logf("Retrieved Router Associations")
	tools.PrintResource(t, allAssocs)
	th.AssertIntGreaterOrEqual(t, len(allAssocs), 1)

	// Get BGP VPN with associations
	getBgpVpn, err := bgpvpns.Get(context.TODO(), client, bgpVpnCreated.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getBgpVpn)
}

func TestBGPVPNPortAssociationCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a BGP VPN
	bgpVpnCreated, err := CreateBGPVPN(t, client)
	th.AssertNoErr(t, err)
	defer func() {
		err = bgpvpns.Delete(context.TODO(), client, bgpVpnCreated.ID).ExtractErr()
		th.AssertNoErr(t, err)
	}()

	// Create a Network
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	// Create port
	port, err := networking.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer networking.DeletePort(t, client, port.ID)

	// Associate a Port with a BGP VPN
	assocOpts := bgpvpns.CreatePortAssociationOpts{
		PortID: port.ID,
	}
	assoc, err := bgpvpns.CreatePortAssociation(context.TODO(), client, bgpVpnCreated.ID, assocOpts).Extract()
	th.AssertNoErr(t, err)
	defer func() {
		err = bgpvpns.DeletePortAssociation(context.TODO(), client, bgpVpnCreated.ID, assoc.ID).ExtractErr()
		th.AssertNoErr(t, err)
	}()
	th.AssertEquals(t, port.ID, assoc.PortID)

	// Get a Port Association
	assocGot, err := bgpvpns.GetPortAssociation(context.TODO(), client, bgpVpnCreated.ID, assoc.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, port.ID, assocGot.PortID)

	// Update a Port Association
	assocUpdOpts := bgpvpns.UpdatePortAssociationOpts{
		AdvertiseFixedIPs: new(bool),
	}
	assocUpdate, err := bgpvpns.UpdatePortAssociation(context.TODO(), client, bgpVpnCreated.ID, assoc.ID, assocUpdOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, port.ID, assocUpdate.PortID)
	th.AssertEquals(t, false, assocUpdate.AdvertiseFixedIPs)

	// List all Port Associations
	allPages, err := bgpvpns.ListPortAssociations(client, bgpVpnCreated.ID, bgpvpns.ListPortAssociationsOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allAssocs, err := bgpvpns.ExtractPortAssociations(allPages)
	th.AssertNoErr(t, err)
	t.Logf("Retrieved Port Associations")
	tools.PrintResource(t, allAssocs)
	th.AssertIntGreaterOrEqual(t, len(allAssocs), 1)

	// Get BGP VPN with associations
	getBgpVpn, err := bgpvpns.Get(context.TODO(), client, bgpVpnCreated.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getBgpVpn)
}
