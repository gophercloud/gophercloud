// +build acceptance compute flavors

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

func TestFlavorsList(t *testing.T) {
	client, err := newClient()
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
		printFlavor(t, &flavor)
	}
}

func TestFlavorsGet(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	flavor, err := flavors.Get(client, choices.FlavorID).Extract()
	if err != nil {
		t.Fatalf("Unable to get flavor information: %v", err)
	}

	printFlavor(t, flavor)
}

func printFlavor(t *testing.T, flavor *flavors.Flavor) {
	t.Logf("ID: %s", flavor.ID)
	t.Logf("Name: %s", flavor.Name)
	t.Logf("RAM: %d", flavor.RAM)
	t.Logf("Disk: %d", flavor.Disk)
	t.Logf("Swap: %d", flavor.Swap)
	t.Logf("RxTxFactor: %f", flavor.RxTxFactor)
}
