package speakers

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/speakers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func CreateBGPSpeaker(t *testing.T, client *gophercloud.ServiceClient) (*speakers.BGPSpeaker, error) {
	iTrue := true
	opts := speakers.CreateOpts{
		IPVersion:                     4,
		AdvertiseFloatingIPHostRoutes: new(bool),
		AdvertiseTenantNetworks:       &iTrue,
		Name:                          tools.RandomString("TESTACC-BGPSPEAKER-", 8),
		LocalAS:                       3000,
	}

	t.Logf("Attempting to create BGP Speaker: %s", opts.Name)
	bgpSpeaker, err := speakers.Create(context.TODO(), client, opts).Extract()
	if err != nil {
		return bgpSpeaker, err
	}

	th.AssertEquals(t, bgpSpeaker.Name, opts.Name)
	th.AssertEquals(t, bgpSpeaker.LocalAS, opts.LocalAS)
	th.AssertEquals(t, bgpSpeaker.IPVersion, opts.IPVersion)
	th.AssertEquals(t, bgpSpeaker.AdvertiseTenantNetworks, *opts.AdvertiseTenantNetworks)
	th.AssertEquals(t, bgpSpeaker.AdvertiseFloatingIPHostRoutes, *opts.AdvertiseFloatingIPHostRoutes)
	t.Logf("Successfully created BGP Speaker")
	tools.PrintResource(t, bgpSpeaker)
	return bgpSpeaker, err
}
