package aggregates

import "github.com/gophercloud/gophercloud"

func aggregatesListURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("os-aggregates")
}

func aggregatesCreateURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("os-aggregates")
}

func aggregatesDeleteURL(c *gophercloud.ServiceClient, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID)
}

func aggregatesGetURL(c *gophercloud.ServiceClient, aggregateID string) string {
	return c.ServiceURL("os-aggregates", aggregateID)
}
