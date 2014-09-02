package openstack

import (
	"errors"

	"github.com/rackspace/gophercloud"
	identity3 "github.com/rackspace/gophercloud/openstack/identity/v3"
	"github.com/rackspace/gophercloud/openstack/utils"
)

// Client provides access to service clients for this OpenStack cloud.
type Client struct {
	gophercloud.ProviderClient
}

const (
	v20 = "v2.0"
	v30 = "v3.0"
)

// NewClient authenticates to an OpenStack cloud with the provided credentials.
// It first queries the root identity endpoint to determine which versions of the identity service are supported, then chooses
// the most recent identity service available to proceed.
func NewClient(authOptions gophercloud.AuthOptions) (*Client, error) {
	versions := []*utils.Version{
		&utils.Version{ID: v20, Priority: 20},
		&utils.Version{ID: v30, Priority: 30},
	}

	chosen, endpoint, err := utils.ChooseVersion(authOptions.IdentityEndpoint, versions)
	if err != nil {
		return nil, err
	}

	client := Client{
		ProviderClient: gophercloud.ProviderClient{
			IdentityEndpoint: endpoint,
			Options:          authOptions,
		},
	}

	switch chosen.ID {
	case v20:
	case v30:
		identityClient := identity3.NewClient(&client.ProviderClient)
		token, err := identityClient.Authenticate(authOptions)
		if err != nil {
			return nil, err
		}

		client.ProviderClient.TokenID = token.ID
	default:
		// The switch must be out of sync with "versions".
		return nil, errors.New("Wat")
	}

	return &client, nil
}

// NewIdentityV3 explicitly accesses the v3 identity service.
func (client *Client) NewIdentityV3() (*identity3.Client, error) {
	return identity3.NewClient(&client.ProviderClient), nil
}
