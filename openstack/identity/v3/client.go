package v3

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
)

// Client abstracts the connection information necessary to make API calls to Identity v3
// resources.
type Client gophercloud.ServiceClient

var (
	nilClient = Client{}
)

// NewClient attempts to authenticate to the v3 identity endpoint. Returns a populated
// IdentityV3Client on success or an error on failure.
func NewClient(authOptions gophercloud.AuthOptions) (*Client, error) {
	client := Client{Options: authOptions}

	result, err := tokens.Create(&client, nil)
	if err != nil {
		return nil, err
	}

	// Assign the token and return.
	client.TokenID = result.TokenID()
	return &client, nil
}
