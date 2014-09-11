package networks

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

const TokenID = "123"

func ServiceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{
			TokenID: TokenID,
		},
		Endpoint: testhelper.Endpoint(),
	}
}
