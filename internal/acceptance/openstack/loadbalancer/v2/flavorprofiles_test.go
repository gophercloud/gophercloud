//go:build acceptance || networking || loadbalancer || flavorprofiles
// +build acceptance networking loadbalancer flavorprofiles

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/flavorprofiles"

	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestFlavorProfilesList(t *testing.T) {
	client, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	allPages, err := flavorprofiles.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allFlavorProfiles, err := flavorprofiles.ExtractFlavorProfiles(allPages)
	th.AssertNoErr(t, err)

	for _, flavorprofile := range allFlavorProfiles {
		tools.PrintResource(t, flavorprofile)
	}
}

func TestFlavorProfilesCRUD(t *testing.T) {
	lbClient, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	flavorProfile, err := CreateFlavorProfile(t, lbClient)
	th.AssertNoErr(t, err)
	defer DeleteFlavorProfile(t, lbClient, flavorProfile)

	tools.PrintResource(t, flavorProfile)

	th.AssertEquals(t, "amphora", flavorProfile.ProviderName)

	flavorProfileUpdateOpts := flavorprofiles.UpdateOpts{
		Name: tools.RandomString("TESTACCTUP-", 8),
	}

	flavorProfileUpdated, err := flavorprofiles.Update(lbClient, flavorProfile.ID, flavorProfileUpdateOpts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, flavorProfileUpdateOpts.Name, flavorProfileUpdated.Name)

	t.Logf("Successfully updated flavorprofile %s", flavorProfileUpdated.Name)
}
