//go:build acceptance || networking || loadbalancer || flavorprofiles

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/internal/ptr"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/availabilityzoneprofiles"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAvailabilityZoneProfilesList(t *testing.T) {
	client, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	allPages, err := availabilityzoneprofiles.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allAvailabilityZoneProfiles, err := availabilityzoneprofiles.ExtractFlavorProfiles(allPages)
	th.AssertNoErr(t, err)

	for _, availabilityzoneprofile := range allAvailabilityZoneProfiles {
		tools.PrintResource(t, availabilityzoneprofile)
	}
}

func TestAvailabilityZoneProfilesCRUD(t *testing.T) {
	lbClient, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	availabilityZoneProfile, err := CreateAvailabilityZonerProfile(t, lbClient)
	th.AssertNoErr(t, err)
	defer DeleteAvailabilityZoneProfile(t, lbClient, availabilityZoneProfile)

	tools.PrintResource(t, availabilityZoneProfile)

	th.AssertEquals(t, "amphora", availabilityZoneProfile.ProviderName)

	availabilityZoneProfileUpdateOpts := availabilityzoneprofiles.UpdateOpts{
		Name: ptr.To(tools.RandomString("TESTACCTUP-", 8)),
	}

	availabilityZoneProfileUpdated, err := availabilityzoneprofiles.Update(context.TODO(), lbClient, availabilityZoneProfile.ID, availabilityZoneProfileUpdateOpts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, *availabilityZoneProfileUpdateOpts.Name, availabilityZoneProfileUpdated.Name)

	t.Logf("Successfully updated availabiltyzoneprofile %s", availabilityZoneProfileUpdated.Name)
}
