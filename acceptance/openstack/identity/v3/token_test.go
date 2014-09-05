// +build acceptance

package v3

import (
	"testing"

	"github.com/rackspace/gophercloud/openstack"
	tokens3 "github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
	"github.com/rackspace/gophercloud/openstack/utils"
)

func TestGetToken(t *testing.T) {
	// Obtain credentials from the environment.
	ao, err := utils.AuthOptions()
	if err != nil {
		t.Fatalf("Unable to acquire credentials: %v", err)
	}

	// Trim out unused fields.
	ao.Username, ao.TenantID, ao.TenantName = "", "", ""

	// Create an unauthenticated client.
	provider, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		t.Fatalf("Unable to instantiate client: %v", err)
	}

	// Create a service client.
	service := openstack.NewIdentityV3(provider)

	// Use the service to create a token.
	result, err := tokens3.Create(service, ao, nil)
	if err != nil {
		t.Fatalf("Unable to get token: %v", err)
	}

	token, err := result.TokenID()
	if err != nil {
		t.Fatalf("Unable to extract token from response: %v", err)
	}

	t.Logf("Acquired token: %s", token)
}
