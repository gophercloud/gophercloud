//go:build acceptance || networking || loadbalancer || providers
// +build acceptance networking loadbalancer providers

package v2

import (
	"testing"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack/loadbalancer/v2/providers"
)

func TestProvidersList(t *testing.T) {
	client, err := clients.NewLoadBalancerV2Client()
	if err != nil {
		t.Fatalf("Unable to create a loadbalancer client: %v", err)
	}

	allPages, err := providers.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list providers: %v", err)
	}

	allProviders, err := providers.ExtractProviders(allPages)
	if err != nil {
		t.Fatalf("Unable to extract providers: %v", err)
	}

	for _, provider := range allProviders {
		tools.PrintResource(t, provider)
	}
}
