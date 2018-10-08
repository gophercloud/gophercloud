package attributestags

import "github.com/gophercloud/gophercloud"

const (
	tagsPath = "tags"
)

func replaceURL(c *gophercloud.ServiceClient, r_type string, id string) string {
	return c.ServiceURL(r_type, id, tagsPath)
}

func listURL(c *gophercloud.ServiceClient, r_type string, id string) string {
	return c.ServiceURL(r_type, id, tagsPath)
}
