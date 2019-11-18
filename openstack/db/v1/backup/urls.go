package backup

import "github.com/gophercloud/gophercloud"

func baseURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("backups")
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("backups", id)
}

func instanceBackupURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "backups")
}
