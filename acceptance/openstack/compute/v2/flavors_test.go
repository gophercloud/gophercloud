// +build acceptance compute flavors

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

func TestFlavorsList(t *testing.T) {
        t.Logf("Default flavors (same as IsPublic set to \"True\"):\t")
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
		PrintFlavor(t, &flavor)
                t.Logf("\t")
	}

        publicOptions := [3]string{"True", "False", "None"}
        for _, publicOption := range publicOptions {
                t.Logf("Flavors for IsPublic option set to \"%s\":\t", publicOption)
                allPages, err := flavors.ListDetail(client, flavors.ListOpts{IsPublic: publicOption}).AllPages()
                if err != nil {
                        t.Fatalf("Unable to retrieve flavors: %v", err)
                }

                allFlavors, err := flavors.ExtractFlavors(allPages)
                if err != nil {
                        t.Fatalf("Unable to extract flavor results: %v", err)
                }

                for _, flavor := range allFlavors {
                        PrintFlavor(t, &flavor)
                        t.Logf("\t")
                }
        }

}

func TestFlavorsGet(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err :=clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	flavor, err := flavors.Get(client, choices.FlavorID).Extract()
	if err != nil {
		t.Fatalf("Unable to get flavor information: %v", err)
	}

	PrintFlavor(t, flavor)
}
