package instancelogs

import "github.com/gophercloud/gophercloud"

func baseURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "log")
}
