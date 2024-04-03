//go:build acceptance || containers || capsules

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/amphorae"
)

func TestAmphoraeList(t *testing.T) {
	clients.RequireAdmin(t)
	client, err := clients.NewLoadBalancerV2Client()
	if err != nil {
		t.Fatalf("Unable to create a loadbalancer client: %v", err)
	}

	allPages, err := amphorae.List(client, nil).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list amphorae: %v", err)
	}

	allAmphorae, err := amphorae.ExtractAmphorae(allPages)
	if err != nil {
		t.Fatalf("Unable to extract amphorae: %v", err)
	}

	for _, amphora := range allAmphorae {
		tools.PrintResource(t, amphora)
	}
}
