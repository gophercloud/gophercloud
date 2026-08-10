//go:build acceptance || networking || loadbalancer || providers

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/providers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestProvidersList(t *testing.T) {
	client, err := clients.NewLoadBalancerV2Client()
	if err != nil {
		t.Fatalf("Unable to create a loadbalancer client: %v", err)
	}

	allPages, err := providers.List(client, nil).AllPages(context.TODO())
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

func TestProviderCapabilities(t *testing.T) {
	client, err := clients.NewLoadBalancerV2Client()
	if err != nil {
		t.Fatalf("Unable to create a loadbalancer client: %v", err)
	}

	allPages, err := providers.List(client, nil).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list providers: %v", err)
	}

	allProviders, err := providers.ExtractProviders(allPages)
	if err != nil {
		t.Fatalf("Unable to extract providers: %v", err)
	}
	provider := ""
	for _, candidate := range allProviders {
		if candidate.Name == "amphora" {
			provider = candidate.Name
			break
		}
	}
	if provider == "" {
		t.Skip("The amphora provider is not enabled")
	}

	flavorCapabilities, err := providers.ListFlavorCapabilities(context.TODO(), client, provider, nil).Extract()
	if err != nil {
		t.Fatalf("Unable to list flavor capabilities for provider %q: %v", provider, err)
	}
	expectedFlavorCapability := providers.Capability{
		Name: "loadbalancer_topology",
		Description: "The load balancer topology. One of: SINGLE - One amphora per load balancer. " +
			"ACTIVE_STANDBY - Two amphora per load balancer.",
	}
	found := false
	for _, capability := range flavorCapabilities {
		tools.PrintResource(t, capability)
		if capability == expectedFlavorCapability {
			found = true
			break
		}
	}
	th.AssertTrue(t, found)

	availabilityZoneCapabilities, err := providers.ListAvailabilityZoneCapabilities(context.TODO(), client, provider, nil).Extract()
	if err != nil {
		t.Fatalf("Unable to list availability zone capabilities for provider %q: %v", provider, err)
	}
	expectedAvailabilityZoneCapability := providers.Capability{
		Name:        "compute_zone",
		Description: "The compute availability zone.",
	}
	found = false
	for _, capability := range availabilityZoneCapabilities {
		tools.PrintResource(t, capability)
		if capability == expectedAvailabilityZoneCapability {
			found = true
			break
		}
	}
	th.AssertTrue(t, found)
}
