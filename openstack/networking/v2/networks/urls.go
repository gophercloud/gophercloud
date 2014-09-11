package networks

import (
	"github.com/rackspace/gophercloud"
)

func APIVersionsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("")
}
