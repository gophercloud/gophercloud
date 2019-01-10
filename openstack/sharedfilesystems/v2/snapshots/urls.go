package snapshots

import "github.com/gophercloud/gophercloud"

func listDetailURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("snapshots", "detail")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("snapshots", id)
}
