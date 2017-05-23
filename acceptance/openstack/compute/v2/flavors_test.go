// +build acceptance compute flavors

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

func TestFlavorsList(t *testing.T) {
	t.Logf("** Default flavors (same as Project flavors): **")
	t.Logf("")
	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := flavors.ListDetail(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve flavors: %v", err)
	}

	allFlavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		t.Fatalf("Unable to extract flavor results: %v", err)
	}

	for _, flavor := range allFlavors {
		tools.PrintResource(t, flavor)
	}

	flavorTypes := [3]flavors.FlavorType{flavors.Project, flavors.Private, flavors.All}
	for _, flavorType := range flavorTypes {
		t.Logf("** %s flavors: **", flavorType)
		t.Logf("")
		allPages, err := flavors.ListDetail(client, flavors.ListOpts{FlavorType: flavorType}).AllPages()
		if err != nil {
			t.Fatalf("Unable to retrieve flavors: %v", err)
		}

		allFlavors, err := flavors.ExtractFlavors(allPages)
		if err != nil {
			t.Fatalf("Unable to extract flavor results: %v", err)
		}

		for _, flavor := range allFlavors {
			tools.PrintResource(t, flavor)
			t.Logf("")
		}
	}

}

func TestFlavorsGet(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	flavor, err := flavors.Get(client, choices.FlavorID).Extract()
	if err != nil {
		t.Fatalf("Unable to get flavor information: %v", err)
	}

	tools.PrintResource(t, flavor)
}
