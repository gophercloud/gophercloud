package networksegmentranges

import "github.com/gophercloud/gophercloud/v2"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("network_segment_ranges")
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("network_segment_ranges", id)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}
