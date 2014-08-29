package v3

import (
	"time"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
)

// Client abstracts the connection information necessary to make API calls to Identity v3 resources.
// It exists mainly to adhere to the IdentityService interface.
type Client struct {
	gophercloud.ServiceClient
}

// Token models a token acquired from the tokens/ API resource.
type Token struct {
	ID        string
	ExpiresAt time.Time
}

// NewClient creates a new client associated with the v3 identity service of a provider.
func NewClient(provider *gophercloud.ProviderClient) *Client {
	return &Client{
		ServiceClient: gophercloud.ServiceClient{
			ProviderClient: *provider,
			Endpoint:       provider.IdentityEndpoint,
		},
	}
}

// Authenticate provides the supplied credentials to an identity v3 endpoint and attempts to acquire a token.
func (c *Client) Authenticate(authOptions gophercloud.AuthOptions) (*Token, error) {
	c.ServiceClient.ProviderClient.Options = authOptions

	result, err := tokens.Create(&c.ServiceClient, nil)
	if err != nil {
		return nil, err
	}

	tokenID, err := result.TokenID()
	if err != nil {
		return nil, err
	}

	expiresAt, err := result.ExpiresAt()
	if err != nil {
		return nil, err
	}

	return &Token{
		ID:        tokenID,
		ExpiresAt: expiresAt,
	}, nil
}
