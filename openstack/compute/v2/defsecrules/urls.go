package defsecrules

import "github.com/gophercloud/gophercloud/v2"

const rulepath = "os-security-group-default-rules"

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rulepath, id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rulepath)
}
