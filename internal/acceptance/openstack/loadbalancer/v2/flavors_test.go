//go:build acceptance || networking || loadbalancer || flavors

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavors"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestFlavorsList(t *testing.T) {
	client, err := clients.NewLoadBalancerV2Client()
	if err != nil {
		t.Fatalf("Unable to create a loadbalancer client: %v", err)
	}

	allPages, err := flavors.List(client, nil).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list flavors: %v", err)
	}

	allFlavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		t.Fatalf("Unable to extract flavors: %v", err)
	}

	for _, flavor := range allFlavors {
		tools.PrintResource(t, flavor)
	}
}

func TestFlavorsCRUD(t *testing.T) {
	lbClient, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	flavorProfile, err := CreateFlavorProfile(t, lbClient)
	th.AssertNoErr(t, err)
	defer DeleteFlavorProfile(t, lbClient, flavorProfile)

	tools.PrintResource(t, flavorProfile)

	th.AssertEquals(t, "amphora", flavorProfile.ProviderName)

	flavor, err := CreateFlavor(t, lbClient, flavorProfile)
	th.AssertNoErr(t, err)
	defer DeleteFlavor(t, lbClient, flavor)

	tools.PrintResource(t, flavor)

	th.AssertEquals(t, flavor.FlavorProfileId, flavorProfile.ID)

	flavorUpdateOpts := flavors.UpdateOpts{
		Name: tools.RandomString("TESTACCTUP-", 8),
	}

	flavorUpdated, err := flavors.Update(context.TODO(), lbClient, flavor.ID, flavorUpdateOpts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, flavorUpdateOpts.Name, flavorUpdated.Name)

	t.Logf("Successfully updated flavor %s", flavorUpdated.Name)
}
