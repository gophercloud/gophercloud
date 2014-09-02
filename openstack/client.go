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

// AuthenticatedClient logs in to an OpenStack cloud found at the identity endpoint specified by authOptions, acquires a token, and
// returns a Client instance that's ready to operate.
// It first queries the root identity endpoint to determine which versions of the identity service are supported, then chooses
// the most recent identity service available to proceed.
func AuthenticatedClient(authOptions gophercloud.AuthOptions) (*Client, error) {
	client := NewClient(authOptions)
	err := client.Authenticate()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewClient prepares an unauthenticated Client instance.
// Most users will probably prefer using the AuthenticatedClient function instead.
// This is useful if you wish to explicitly control the version of the identity service that's used for authentication explicitly,
// for example.
func NewClient(authOptions gophercloud.AuthOptions) *Client {
	return &Client{
		ProviderClient: gophercloud.ProviderClient{
			Options: authOptions,
		},
	}
}

// Authenticate or re-authenticate against the most recent identity service supported at the provided endpoint.
func (client *Client) Authenticate() error {
	versions := []*utils.Version{
		&utils.Version{ID: v20, Priority: 20},
		&utils.Version{ID: v30, Priority: 30},
	}

	chosen, endpoint, err := utils.ChooseVersion(client.ProviderClient.Options.IdentityEndpoint, versions)
	if err != nil {
		return err
	}

	switch chosen.ID {
	case v20:
		return client.authenticateV2(endpoint)
	case v30:
		return client.authenticateV3(endpoint)
	default:
		// The switch statement must be out of date from the versions list.
		return errors.New("Wat")
	}
}

// AuthenticateV2 acquires a token explicitly from the v2.0 identity API.
func (client *Client) AuthenticateV2() error {
	endpoint := client.ProviderClient.Options.IdentityEndpoint + "/v2.0"
	return client.authenticateV2(endpoint)
}

func (client *Client) authenticateV2(endpoint string) error {
	return errors.New("Not implemented yet.")
}

// AuthenticateV3 acquires a token explicitly from the v3.0 identity API.
func (client *Client) AuthenticateV3() error {
	endpoint := client.ProviderClient.Options.IdentityEndpoint + "/v3"
	return client.authenticateV3(endpoint)
}

func (client *Client) authenticateV3(endpoint string) error {
	identityClient := identity3.NewClient(&client.ProviderClient, endpoint)
	token, err := identityClient.GetToken(client.ProviderClient.Options)
	if err != nil {
		return err
	}

	client.ProviderClient.TokenID = token.ID

	return nil
}

// NewIdentityV3 explicitly accesses the v3 identity service.
func (client *Client) NewIdentityV3() (*identity3.Client, error) {
	endpoint := client.ProviderClient.Options.IdentityEndpoint + "/v3"
	return identity3.NewClient(&client.ProviderClient, endpoint), nil
}
