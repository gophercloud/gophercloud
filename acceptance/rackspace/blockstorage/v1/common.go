package v1

import (
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/rackspace"
	th "github.com/rackspace/gophercloud/testhelper"
)

func newClient() (*gophercloud.ServiceClient, error) {
	opts, err := rackspace.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	provider, err := rackspace.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	return rackspace.NewBlockStorageV1(provider, gophercloud.EndpointOpts{
		Region: "IAD",
	})
}

func setup(t *testing.T) *gophercloud.ServiceClient {
	client, err := newClient()
	th.AssertNoErr(t, err)

	return client
}
