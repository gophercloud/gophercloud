// +build acceptance

package v2

import (
	"testing"

	"github.com/rackspace/gophercloud/openstack"
	tokens2 "github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
	"github.com/rackspace/gophercloud/openstack/utils"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestAuthenticate(t *testing.T) {
	// Obtain credentials from the environment.
	ao, err := utils.AuthOptions()
	th.AssertNoErr(t, err)

	// Trim out unused fields. Prefer authentication by API key to password.
	ao.UserID, ao.DomainID, ao.DomainName = "", "", ""
	if ao.APIKey != "" {
		ao.Password = ""
	}

	// Create an unauthenticated client.
	provider, err := openstack.NewClient(ao.IdentityEndpoint)
	th.AssertNoErr(t, err)

	// Create a service client.
	service := openstack.NewIdentityV2(provider)

	// Authenticated!
	result := tokens2.Create(service, ao)

	// Extract and print the token.
	token, err := result.ExtractToken()
	th.AssertNoErr(t, err)

	t.Logf("Acquired token: [%s]", token.ID)
	t.Logf("The token will expire at: [%s]", token.ExpiresAt.String())
	t.Logf("The token is valid for tenant: [%#v]", token.Tenant)

	// Extract and print the service catalog.
	catalog, err := result.ExtractServiceCatalog()
	th.AssertNoErr(t, err)

	t.Logf("Acquired service catalog listing [%d] services", len(catalog.Entries))
	for i, entry := range catalog.Entries {
		t.Logf("[%d]: name=[%s], type=[%s]", i, entry.Name, entry.Type)
		for _, endpoint := range entry.Endpoints {
			t.Logf("      - region=[%s] publicURL=[%s]", endpoint.Region, endpoint.PublicURL)
		}
	}
}
