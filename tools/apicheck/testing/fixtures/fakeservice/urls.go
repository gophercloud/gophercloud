package fakeservice

import "github.com/gophercloud/gophercloud/v2"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("widgets")
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("widgets", id)
}

func actionURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("widgets", id, "action")
}

// deleteURL exercises nested url-func resolution: it delegates to resourceURL
// rather than calling ServiceURL directly.
func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}
