package rackspace

import (
	"errors"

	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack"
)

const (
	// RackspaceUSIdentity is an identity endpoint located in the United States.
	RackspaceUSIdentity = "https://identity.api.rackspacecloud.com/v2.0/"

	// RackspaceUKIdentity is an identity endpoint located in the UK.
	RackspaceUKIdentity = "https://lon.identity.api.rackspacecloud.com/v2.0/"
)

// NewClient creates a client that's prepared to communicate with the Rackspace API, but is not
// yet authenticated. Most users will probably prefer using the AuthenticatedClient function
// instead.
//
// Provide the base URL of the identity endpoint you wish to authenticate against as "endpoint".
// Often, this will be either RackspaceUSIdentity or RackspaceUKIdentity.
func NewClient(endpoint string) (*gophercloud.ProviderClient, error) {
	return os.NewClient(endpoint)
}

// AuthenticatedClient logs in to Rackspace with the provided credentials and constructs a
// ProviderClient that's ready to operate.
//
// If the provided AuthOptions does not specify an explicit IdentityEndpoint, it will default to
// the canonical, production Rackspace US identity endpoint.
func AuthenticatedClient(options gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
	if options.IdentityEndpoint == "" {
		options.IdentityEndpoint = RackspaceUSIdentity
	}

	_, err := NewClient(options.IdentityEndpoint)
	if err != nil {
		return nil, err
	}

	return nil, errors.New("Incomplete")
}
