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
