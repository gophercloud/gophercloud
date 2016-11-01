package zones

import "github.com/gophercloud/gophercloud"

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("v2/zones")
}

func zoneURL(c *gophercloud.ServiceClient, zoneID string) string {
	return c.ServiceURL("v2/zones", zoneID)
}
