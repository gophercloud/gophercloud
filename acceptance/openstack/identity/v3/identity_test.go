package v3

import (
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/utils"
)

func createAuthenticatedClient(t *testing.T) *gophercloud.ServiceClient {
	// Obtain credentials from the environment.
	ao, err := utils.AuthOptions()
	if err != nil {
		t.Fatalf("Unable to acquire credentials: %v", err)
	}

	// Trim out unused fields.
	ao.Username, ao.TenantID, ao.TenantName = "", "", ""

	// Create an authenticated client.
	providerClient, err := openstack.AuthenticatedClient(ao)
	if err != nil {
		t.Fatalf("Unable to instantiate client: %v", err)
	}

	// Create a service client.
	return openstack.NewIdentityV3(providerClient)
}
