package pagination

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

func createClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: "abc123"},
		Endpoint: testhelper.Endpoint(),
	}
}
