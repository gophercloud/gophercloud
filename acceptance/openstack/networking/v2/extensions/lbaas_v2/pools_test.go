// +build acceptance networking loadbalancer pools

package lbaas_v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/pools"
)

func TestPoolsList(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/pike")
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a loadbalancer client: %v", err)
	}

	allPages, err := pools.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list pools: %v", err)
	}

	allPools, err := pools.ExtractPools(allPages)
	if err != nil {
		t.Fatalf("Unable to extract pools: %v", err)
	}

	for _, pool := range allPools {
		tools.PrintResource(t, pool)
	}
}
