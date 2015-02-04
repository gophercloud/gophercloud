package stacks

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
)

// Create accepts an os.CreateOpts struct and creates a new stack using the values
// provided.
func Create(c *gophercloud.ServiceClient, opts os.CreateOptsBuilder) os.CreateResult {
	return os.Create(c, opts)
}

// Adopt accepts an os.AdoptOpts struct and creates a new stack from existing stack
// resources using the values provided.
func Adopt(c *gophercloud.ServiceClient, opts os.AdoptOptsBuilder) os.AdoptResult {
	return os.Adopt(c, opts)
}
