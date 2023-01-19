package speakers

import (
	"strconv"
	"testing"

	"github.com/bizflycloud/gophercloud"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack/networking/v2/extensions/bgp/speakers"
	th "github.com/bizflycloud/gophercloud/testhelper"
)

func CreateBGPSpeaker(t *testing.T, client *gophercloud.ServiceClient) (*speakers.BGPSpeaker, error) {
	opts := speakers.CreateOpts{
		IPVersion:                     4,
		AdvertiseFloatingIPHostRoutes: false,
		AdvertiseTenantNetworks:       true,
		Name:                          tools.RandomString("TESTACC-BGPSPEAKER-", 8),
		LocalAS:                       "3000",
		Networks:                      []string{},
	}

	t.Logf("Attempting to create BGP Speaker: %s", opts.Name)
	bgpSpeaker, err := speakers.Create(client, opts).Extract()
	if err != nil {
		return bgpSpeaker, err
	}

	localas, err := strconv.Atoi(opts.LocalAS)
	th.AssertEquals(t, bgpSpeaker.Name, opts.Name)
	th.AssertEquals(t, bgpSpeaker.LocalAS, localas)
	th.AssertEquals(t, bgpSpeaker.IPVersion, opts.IPVersion)
	th.AssertEquals(t, bgpSpeaker.AdvertiseTenantNetworks, opts.AdvertiseTenantNetworks)
	th.AssertEquals(t, bgpSpeaker.AdvertiseFloatingIPHostRoutes, opts.AdvertiseFloatingIPHostRoutes)
	t.Logf("Successfully created BGP Speaker")
	tools.PrintResource(t, bgpSpeaker)
	return bgpSpeaker, err
}
