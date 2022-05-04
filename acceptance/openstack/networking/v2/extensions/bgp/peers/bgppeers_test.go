package peers

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/peers"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func CreateBGPPeer(t *testing.T, client *gophercloud.ServiceClient) (*peers.BGPPeer, error) {
	var opts peers.CreateOpts
	opts.AuthType = "md5"
	opts.Password = tools.MakeNewPassword("")
	opts.RemoteAS = tools.RandomInt(1000, 2000)
	opts.Name = tools.RandomString("TESTACC-BGPPEER-", 8)
	opts.PeerIP = "192.168.0.1"

	t.Logf("Attempting to create BGP Peer: %s", opts.Name)
	bgpPeer, err := peers.Create(client, opts).Extract()
	if err != nil {
		return bgpPeer, err
	}

	t.Logf("Successfully created BGP Peer")
	th.AssertEquals(t, bgpPeer.Name, opts.Name)
	th.AssertEquals(t, bgpPeer.RemoteAS, opts.RemoteAS)
	th.AssertEquals(t, bgpPeer.PeerIP, opts.PeerIP)
	return bgpPeer, err
}

func TestBGPPeerCRD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	allPages, err := peers.List(client).AllPages()
	th.AssertNoErr(t, err)

	allPeers, err := peers.ExtractBGPPeers(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved BGP Peers")
	tools.PrintResource(t, allPeers)

	bgpPeerCreated, err := CreateBGPPeer(t, client)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, bgpPeerCreated)

	bgpPeerGot, err := peers.Get(client, bgpPeerCreated.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bgpPeerCreated.ID, bgpPeerGot.ID)
	th.AssertEquals(t, bgpPeerCreated.Name, bgpPeerGot.Name)

	t.Logf("Attempting to delete BGP Peer: %s", bgpPeerGot.Name)
	err = peers.Delete(client, bgpPeerGot.ID).ExtractErr()
	th.AssertNoErr(t, err)

	bgpPeerGot, err = peers.Get(client, bgpPeerGot.ID).Extract()
	th.AssertErr(t, err)
	t.Logf("BGP Peer %s deleted", bgpPeerCreated.Name)
}
