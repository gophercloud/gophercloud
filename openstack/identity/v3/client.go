package v3

import (
	"time"

	"github.com/rackspace/gophercloud"
)

// Client abstracts the connection information necessary to make API calls to Identity v3 resources.
// It exists mainly to adhere to the IdentityService interface.
type Client gophercloud.ServiceClient

// Token models a token acquired from the tokens/ API resource.
type Token struct {
	ID        string
	ExpiresAt time.Time
}

// NewClient creates a new client associated with the v3 identity service of a provider.
func NewClient(provider *gophercloud.ProviderClient) *Client {
	return &Client{
		ProviderClient: *provider,
		Endpoint:       provider.IdentityEndpoint + "v3/",
	}
}

// Authenticate provides the supplied credentials to an identity v3 endpoint and attempts to acquire a token.
func (c *Client) Authenticate(authOptions gophercloud.AuthOptions) (*Token, error) {
	return nil, nil
}
