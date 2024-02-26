package peers

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/peers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func CreateBGPPeer(t *testing.T, client *gophercloud.ServiceClient) (*peers.BGPPeer, error) {
	var opts peers.CreateOpts
	opts.AuthType = "md5"
	opts.Password = tools.MakeNewPassword("")
	opts.RemoteAS = tools.RandomInt(1000, 2000)
	opts.Name = tools.RandomString("TESTACC-BGPPEER-", 8)
	opts.PeerIP = "192.168.0.1"

	t.Logf("Attempting to create BGP Peer: %s", opts.Name)
	bgpPeer, err := peers.Create(context.TODO(), client, opts).Extract()
	if err != nil {
		return bgpPeer, err
	}

	th.AssertEquals(t, bgpPeer.Name, opts.Name)
	th.AssertEquals(t, bgpPeer.RemoteAS, opts.RemoteAS)
	th.AssertEquals(t, bgpPeer.PeerIP, opts.PeerIP)
	t.Logf("Successfully created BGP Peer")
	tools.PrintResource(t, bgpPeer)
	return bgpPeer, err
}
