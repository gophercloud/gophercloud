package schedulerstats

import "github.com/gophercloud/gophercloud"

func poolsListURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("scheduler-stats", "pools")
}

func poolsListDetailURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("scheduler-stats", "pools", "detail")
}
