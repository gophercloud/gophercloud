//go:build acceptance || identity || endpoints

package v3

import (
	"context"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/endpoints"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/services"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestEndpointsList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := endpoints.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allEndpoints, err := endpoints.ExtractEndpoints(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, endpoint := range allEndpoints {
		tools.PrintResource(t, endpoint)

		if strings.Contains(endpoint.URL, "/v3") {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestEndpointsGet(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := endpoints.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allEndpoints, err := endpoints.ExtractEndpoints(allPages)
	th.AssertNoErr(t, err)

	endpoint := allEndpoints[0]
	e, err := endpoints.Get(context.TODO(), client, endpoint.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get endpoint: %v", err)
	}

	tools.PrintResource(t, e)

	th.AssertEquals(t, e.Name, e.Name)
}

func TestEndpointsNavigateCatalog(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	// Discover the service we're interested in.
	serviceListOpts := services.ListOpts{
		ServiceType: "compute",
	}

	allPages, err := services.List(client, serviceListOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServices, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, len(allServices), 1)

	computeService := allServices[0]
	tools.PrintResource(t, computeService)

	// Enumerate the endpoints available for this service.
	endpointListOpts := endpoints.ListOpts{
		Availability: gophercloud.AvailabilityPublic,
		ServiceID:    computeService.ID,
	}

	allPages, err = endpoints.List(client, endpointListOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allEndpoints, err := endpoints.ExtractEndpoints(allPages)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, len(allServices), 1)

	tools.PrintResource(t, allEndpoints[0])
}

func TestEndpointCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	service, err := CreateService(t, client, &services.CreateOpts{
		Type:  "endpoint-test",
		Name:  tools.RandomString("ACPTTEST", 8),
		Extra: map[string]any{},
	})
	th.AssertNoErr(t, err)
	defer DeleteService(t, client, service.ID)

	endpoint, err := CreateEndpoint(t, client, &endpoints.CreateOpts{
		Availability: gophercloud.Availability("internal"),
		ServiceID:    service.ID,
		URL:          "https://example.com",
	})
	th.AssertNoErr(t, err)
	defer DeleteEndpoint(t, client, endpoint.ID)

	tools.PrintResource(t, endpoint)
	tools.PrintResource(t, endpoint.URL)

	newEndpoint, err := endpoints.Update(context.TODO(), client, endpoint.ID, &endpoints.UpdateOpts{
		Name:        "new-endpoint",
		URL:         "https://example-updated.com",
		Description: "Updated Endpoint",
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, newEndpoint.URL, "https://example-updated.com")
	th.AssertEquals(t, newEndpoint.Description, "Updated Endpoint")
}
