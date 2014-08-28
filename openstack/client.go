package openstack

import (
	"github.com/rackspace/gophercloud"
	identity3 "github.com/rackspace/gophercloud/openstack/identity/v3"
)

// Client provides access to service clients for this OpenStack cloud.
type Client gophercloud.ServiceClient

// NewClient authenticates to an OpenStack cloud with the provided credentials.
// It first queries the root identity endpoint to determine which versions of the identity service are supported, then chooses
// the most recent identity service available to proceed.
func NewClient(authOptions gophercloud.AuthOptions) (*Client, error) {
	return nil, nil
}

// IdentityV3 explicitly accesses the v3 identity service.
func (client *Client) IdentityV3() (*identity3.Client, error) {
	return nil, nil
}
