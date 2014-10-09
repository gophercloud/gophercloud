// +build acceptance

package v2

import (
	"os"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/rackspace"
	th "github.com/rackspace/gophercloud/testhelper"
)

func rackspaceAuthOptions(t *testing.T) gophercloud.AuthOptions {
	// Obtain credentials from the environment.
	options := gophercloud.AuthOptions{
		Username: os.Getenv("RS_USERNAME"),
		APIKey:   os.Getenv("RS_APIKEY"),
	}

	if options.Username == "" {
		t.Fatal("Please provide a Rackspace username as RS_USERNAME.")
	}
	if options.APIKey == "" {
		t.Fatal("Please provide a Rackspace API key as RS_APIKEY.")
	}

	return options
}

func createClient(t *testing.T, auth bool) *gophercloud.ServiceClient {
	ao := rackspaceAuthOptions(t)

	provider, err := rackspace.NewClient(ao.IdentityEndpoint)
	th.AssertNoErr(t, err)

	if auth {
		err = rackspace.Authenticate(provider, ao)
		th.AssertNoErr(t, err)
	}

	return rackspace.NewIdentityV2(provider)
}

func unauthenticatedClient(t *testing.T) *gophercloud.ServiceClient {
	return createClient(t, false)
}

func authenticatedClient(t *testing.T) *gophercloud.ServiceClient {
	return createClient(t, true)
}
