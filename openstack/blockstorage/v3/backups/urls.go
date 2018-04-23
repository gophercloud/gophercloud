package backups

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("backups")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("backups", id)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return deleteURL(c, id)
}

func listURL(c *gophercloud.ServiceClient) string {
	return createURL(c)
}

func metadataURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("backups", id, "metadata")
}

func updateMetadataURL(c *gophercloud.ServiceClient, id string) string {
	return metadataURL(c, id)
}
