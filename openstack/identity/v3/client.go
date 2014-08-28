package v3

import "github.com/rackspace/gophercloud"

// Client abstracts the connection information necessary to make API calls to Identity v3
// resources.
type Client struct {
	gophercloud.ServiceClient

	// TokenID is redudant storage for an active token.
	// The Identity service occasionally needs to access the assigned token directly, but I don't want to export it from all
	// service clients unless we absolutely need to.
	TokenID string
}

var (
	nilClient = Client{}
)

// NewClient attempts to authenticate to the v3 identity endpoint. Returns a populated
// IdentityV3Client on success or an error on failure.
func NewClient(authOptions gophercloud.AuthOptions) (*Client, error) {
	return &nilClient, nil
}
