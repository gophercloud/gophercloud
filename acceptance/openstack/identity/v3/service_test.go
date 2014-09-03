// +build acceptance

package v3

import (
	"testing"

	"github.com/rackspace/gophercloud/openstack"
	services3 "github.com/rackspace/gophercloud/openstack/identity/v3/services"
	"github.com/rackspace/gophercloud/openstack/utils"
)

func TestListServices(t *testing.T) {
	// Obtain credentials from the environment.
	ao, err := utils.AuthOptions()
	if err != nil {
		t.Fatalf("Unable to acquire credentials: %v", err)
	}

	// Trim out unused fields.
	ao.TenantID, ao.TenantName = "", ""

	// Create an authenticated client.
	providerClient, err := openstack.AuthenticatedClient(ao)
	if err != nil {
		t.Fatalf("Unable to instantiate client: %v", err)
	}

	// Create a service client.
	serviceClient := openstack.NewIdentityV3(providerClient)

	// Use the service to create a token.
	results, err := services3.List(serviceClient, services3.ListOpts{})
	if err != nil {
		t.Fatalf("Unable to get token: %v", err)
	}

	for _, service := range results.Services {
		t.Logf("Service: %32s %15s %10s %s", service.ID, service.Type, service.Name, *service.Description)
	}
}
