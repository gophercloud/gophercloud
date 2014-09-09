// +build acceptance

package v3

import (
	"testing"

	"github.com/rackspace/gophercloud"
	endpoints3 "github.com/rackspace/gophercloud/openstack/identity/v3/endpoints"
	services3 "github.com/rackspace/gophercloud/openstack/identity/v3/services"
)

func TestListEndpoints(t *testing.T) {
	// Create a service client.
	serviceClient := createAuthenticatedClient(t)
	if serviceClient == nil {
		return
	}

	// Use the service to list all available endpoints.
	results, err := endpoints3.List(serviceClient, endpoints3.ListOpts{})
	if err != nil {
		t.Fatalf("Unexpected error while listing endpoints: %v", err)
	}

	err = gophercloud.EachPage(results, func(page gophercloud.Collection) bool {
		t.Logf("--- Page ---")

		for _, endpoint := range endpoints3.AsEndpoints(page) {
			t.Logf("Endpoint: %8s %10s %9s %s",
				endpoint.ID,
				endpoint.Interface,
				endpoint.Name,
				endpoint.URL)
		}

		return true
	})
	if err != nil {
		t.Errorf("Unexpected error while iterating endpoint pages: %v", err)
	}
}

func TestNavigateCatalog(t *testing.T) {
	// Create a service client.
	client := createAuthenticatedClient(t)

	// Discover the service we're interested in.
	computeResults, err := services3.List(client, services3.ListOpts{ServiceType: "compute"})
	if err != nil {
		t.Fatalf("Unexpected error while listing services: %v", err)
	}

	allServices, err := gophercloud.AllPages(computeResults)
	if err != nil {
		t.Fatalf("Unexpected error while traversing service results: %v", err)
	}

	computeServices := services3.AsServices(allServices)

	if len(computeServices) != 1 {
		t.Logf("%d compute services are available at this endpoint.", len(computeServices))
		return
	}
	computeService := computeServices[0]

	// Enumerate the endpoints available for this service.
	endpointResults, err := endpoints3.List(client, endpoints3.ListOpts{
		Interface: gophercloud.InterfacePublic,
		ServiceID: computeService.ID,
	})

	allEndpoints, err := gophercloud.AllPages(endpointResults)
	if err != nil {
		t.Fatalf("Unexpected error while listing endpoints: %v", err)
	}

	endpoints := endpoints3.AsEndpoints(allEndpoints)

	if len(endpoints) != 1 {
		t.Logf("%d endpoints are available for the service %v.", len(endpoints), computeService.Name)
		return
	}

	endpoint := endpoints[0]
	t.Logf("Success. The compute endpoint is at %s.", endpoint.URL)
}
