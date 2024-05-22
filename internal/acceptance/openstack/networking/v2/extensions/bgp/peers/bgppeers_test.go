//go:build acceptance || networking || bgp || peers

package peers

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/peers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestBGPPeerCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a BGP Peer
	bgpPeerCreated, err := CreateBGPPeer(t, client)
	th.AssertNoErr(t, err)

	// Get a BGP Peer
	bgpPeerGot, err := peers.Get(context.TODO(), client, bgpPeerCreated.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bgpPeerCreated.ID, bgpPeerGot.ID)
	th.AssertEquals(t, bgpPeerCreated.Name, bgpPeerGot.Name)

	// Update a BGP Peer
	newBGPPeerName := tools.RandomString("TESTACC-BGPPEER-", 10)
	updateBGPOpts := peers.UpdateOpts{
		Name:     newBGPPeerName,
		Password: tools.MakeNewPassword(""),
	}
	bgpPeerUpdated, err := peers.Update(context.TODO(), client, bgpPeerGot.ID, updateBGPOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bgpPeerUpdated.Name, newBGPPeerName)
	t.Logf("Update BGP Peer, renamed from %s to %s", bgpPeerGot.Name, bgpPeerUpdated.Name)

	// List all BGP Peers
	allPages, err := peers.List(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allPeers, err := peers.ExtractBGPPeers(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved BGP Peers")
	tools.PrintResource(t, allPeers)
	th.AssertIntGreaterOrEqual(t, len(allPeers), 1)

	// Delete a BGP Peer
	t.Logf("Attempting to delete BGP Peer: %s", bgpPeerUpdated.Name)
	err = peers.Delete(context.TODO(), client, bgpPeerGot.ID).ExtractErr()
	th.AssertNoErr(t, err)

	_, err = peers.Get(context.TODO(), client, bgpPeerGot.ID).Extract()
	th.AssertErr(t, err)
	t.Logf("BGP Peer %s deleted", bgpPeerUpdated.Name)
}
