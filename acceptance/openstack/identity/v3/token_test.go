// +build acceptance

package v3

import (
	"testing"

	"github.com/rackspace/gophercloud"
	identity3 "github.com/rackspace/gophercloud/openstack/identity/v3"
	"github.com/rackspace/gophercloud/openstack/utils"
)

func TestGetToken(t *testing.T) {
	// Obtain credentials from the environment.
	ao, err := utils.AuthOptions()
	if err != nil {
		t.Fatalf("Unable to acquire credentials: %v", err)
	}

	client := identity3.NewClient(&gophercloud.ProviderClient{
		Options: ao,
	}, ao.IdentityEndpoint+"/v3/")

	// Attempt to acquire a token.
	token, err := client.GetToken(ao)
	if err != nil {
		t.Fatalf("Unable to get token: %v", err)
	}

	t.Logf("Acquired token: %s", token.ID)
}
