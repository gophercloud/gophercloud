package bgpvpns

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgpvpns"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func CreateBGPVPN(t *testing.T, client *gophercloud.ServiceClient) (*bgpvpns.BGPVPN, error) {
	opts := bgpvpns.CreateOpts{
		Name: tools.RandomString("TESTACC-BGPVPN-", 10),
	}

	t.Logf("Attempting to create BGP VPN: %s", opts.Name)
	bgpVpn, err := bgpvpns.Create(context.TODO(), client, opts).Extract()
	if err != nil {
		return bgpVpn, err
	}

	th.AssertEquals(t, bgpVpn.Name, opts.Name)
	t.Logf("Successfully created BGP VPN")
	tools.PrintResource(t, bgpVpn)
	return bgpVpn, err
}
