package quotas

import "github.com/rackspace/gophercloud"

const version = "v2.0"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(version, "quotas")
}
