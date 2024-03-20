//go:build acceptance || networking || bgp || speakers

package speakers

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	ap "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2/extensions/bgp/peers"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/peers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/speakers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestBGPSpeakerCRUD(t *testing.T) {
	clients.RequireAdmin(t)
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a BGP Speaker
	bgpSpeaker, err := CreateBGPSpeaker(t, client)
	th.AssertNoErr(t, err)

	// Create a BGP Peer
	bgpPeer, err := ap.CreateBGPPeer(t, client)
	th.AssertNoErr(t, err)

	// List BGP Speakers
	allPages, err := speakers.List(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allSpeakers, err := speakers.ExtractBGPSpeakers(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved BGP Speakers")
	tools.PrintResource(t, allSpeakers)
	th.AssertIntGreaterOrEqual(t, len(allSpeakers), 1)

	// Create a network
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	// Update BGP Speaker
	opts := speakers.UpdateOpts{
		Name:                          tools.RandomString("TESTACC-BGPSPEAKER-", 10),
		AdvertiseTenantNetworks:       false,
		AdvertiseFloatingIPHostRoutes: true,
	}
	speakerUpdated, err := speakers.Update(context.TODO(), client, bgpSpeaker.ID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, speakerUpdated.Name, opts.Name)
	t.Logf("Updated the BGP Speaker, name set from %s to %s", bgpSpeaker.Name, speakerUpdated.Name)

	// Get a BGP Speaker
	bgpSpeakerGot, err := speakers.Get(context.TODO(), client, bgpSpeaker.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bgpSpeaker.ID, bgpSpeakerGot.ID)
	th.AssertEquals(t, opts.Name, bgpSpeakerGot.Name)

	// AddBGPPeer
	addBGPPeerOpts := speakers.AddBGPPeerOpts{BGPPeerID: bgpPeer.ID}
	_, err = speakers.AddBGPPeer(context.TODO(), client, bgpSpeaker.ID, addBGPPeerOpts).Extract()
	th.AssertNoErr(t, err)
	speakerGot, err := speakers.Get(context.TODO(), client, bgpSpeaker.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bgpPeer.ID, speakerGot.Peers[0])
	t.Logf("Successfully added BGP Peer %s to BGP Speaker %s", bgpPeer.Name, speakerUpdated.Name)

	// RemoveBGPPeer
	removeBGPPeerOpts := speakers.RemoveBGPPeerOpts{BGPPeerID: bgpPeer.ID}
	err = speakers.RemoveBGPPeer(context.TODO(), client, bgpSpeaker.ID, removeBGPPeerOpts).ExtractErr()
	th.AssertNoErr(t, err)
	speakerGot, err = speakers.Get(context.TODO(), client, bgpSpeaker.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(speakerGot.Networks), 0)
	t.Logf("Successfully removed BGP Peer %s to BGP Speaker %s", bgpPeer.Name, speakerUpdated.Name)

	// GetAdvertisedRoutes
	pages, err := speakers.GetAdvertisedRoutes(client, bgpSpeaker.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	routes, err := speakers.ExtractAdvertisedRoutes(pages)
	th.AssertNoErr(t, err)
	th.AssertIntGreaterOrEqual(t, len(routes), 0)
	t.Logf("Successfully retrieved advertised routes")

	// AddGatewayNetwork
	optsAddGatewayNetwork := speakers.AddGatewayNetworkOpts{NetworkID: network.ID}
	r, err := speakers.AddGatewayNetwork(context.TODO(), client, bgpSpeaker.ID, optsAddGatewayNetwork).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, r.NetworkID, network.ID)
	t.Logf("Successfully added gateway network %s to BGP Speaker", network.ID)

	// RemoveGatewayNetwork
	optsRemoveGatewayNetwork := speakers.RemoveGatewayNetworkOpts{NetworkID: network.ID}
	err = speakers.RemoveGatewayNetwork(context.TODO(), client, bgpSpeaker.ID, optsRemoveGatewayNetwork).ExtractErr()
	th.AssertNoErr(t, err)
	t.Logf("Successfully removed gateway network %s to BGP Speaker", network.ID)

	// Delete a BGP Peer
	t.Logf("Delete the BGP Peer %s", bgpPeer.Name)
	err = peers.Delete(context.TODO(), client, bgpPeer.ID).ExtractErr()
	th.AssertNoErr(t, err)

	// Delete a BGP Speaker
	t.Logf("Delete the BGP Speaker %s", speakerUpdated.Name)
	err = speakers.Delete(context.TODO(), client, bgpSpeaker.ID).ExtractErr()
	th.AssertNoErr(t, err)
}
