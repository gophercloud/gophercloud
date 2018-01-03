// +build acceptance networking subnetpools

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/subnetpools"
)

func TestSubnetPoolsCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	// Create a subnetpool
	subnetPool, err := CreateSubnetPool(t, client)
	if err != nil {
		t.Fatalf("Unable to create a subnetpool: %v", err)
	}

	tools.PrintResource(t, subnetPool)
}

func TestSubnetPoolsList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	allPages, err := subnetpools.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list subnetpools: %v", err)
	}

	allSubnetPools, err := subnetpools.ExtractSubnetPools(allPages)
	if err != nil {
		t.Fatalf("Unable to extract subnetpools: %v", err)
	}

	for _, subnetpool := range allSubnetPools {
		tools.PrintResource(t, subnetpool)
	}
}
