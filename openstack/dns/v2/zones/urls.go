package zones

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("zones")
}

func ZoneURL(client *gophercloud.ServiceClient, parts ...string) string {
	return client.ServiceURL(append([]string{"zones"}, parts...)...)
}
