package speakers

import (
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/speakers"
	th "github.com/gophercloud/gophercloud/testhelper"
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
	t.Logf("Successfully created BGP Speaker")
	th.AssertEquals(t, bgpSpeaker.Name, opts.Name)
	th.AssertEquals(t, bgpSpeaker.LocalAS, localas)
	th.AssertEquals(t, bgpSpeaker.IPVersion, opts.IPVersion)
	th.AssertEquals(t, bgpSpeaker.AdvertiseTenantNetworks, opts.AdvertiseTenantNetworks)
	th.AssertEquals(t, bgpSpeaker.AdvertiseFloatingIPHostRoutes, opts.AdvertiseFloatingIPHostRoutes)
	return bgpSpeaker, err
}

func TestBGPSpeakerCRD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a BGP Speaker
	bgpSpeakerCreated, err := CreateBGPSpeaker(t, client)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, bgpSpeakerCreated)

	// Get a BGP Speaker
	bgpSpeakerGot, err := speakers.Get(client, bgpSpeakerCreated.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bgpSpeakerCreated.ID, bgpSpeakerGot.ID)
	th.AssertEquals(t, bgpSpeakerCreated.Name, bgpSpeakerGot.Name)

	// List BGP Speakers
	allPages, err := speakers.List(client).AllPages()
	th.AssertNoErr(t, err)
	allSpeakers, err := speakers.ExtractBGPSpeakers(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved BGP Speakers")
	tools.PrintResource(t, allSpeakers)
	th.AssertIntGreaterOrEqual(t, len(allSpeakers), 1)

	// Delete a BGP Speaker
	t.Logf("Attempting to delete BGP Speaker: %s", bgpSpeakerGot.Name)
	err = speakers.Delete(client, bgpSpeakerGot.ID).ExtractErr()
	th.AssertNoErr(t, err)

	// Confirm the BGP Speaker is deleted
	bgpSpeakerGot, err = speakers.Get(client, bgpSpeakerGot.ID).Extract()
	th.AssertErr(t, err)
	t.Logf("BGP Speaker %s deleted", bgpSpeakerCreated.Name)
}
