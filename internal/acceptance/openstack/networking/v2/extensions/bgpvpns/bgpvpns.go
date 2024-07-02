package bgpvpns

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgpvpns"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func CreateBGPVPN(t *testing.T, client *gophercloud.ServiceClient) (*bgpvpns.BGPVPN, error) {
	opts := bgpvpns.CreateOpts{
		Name: tools.RandomString("TESTACC-BGPVPN-", 10),
	}

	t.Logf("Attempting to create BGP VPN: %s", opts.Name)
	bgpVpn, err := bgpvpns.Create(client, opts).Extract()
	if err != nil {
		return bgpVpn, err
	}

	th.AssertEquals(t, bgpVpn.Name, opts.Name)
	t.Logf("Successfully created BGP VPN")
	tools.PrintResource(t, bgpVpn)
	return bgpVpn, err
}
