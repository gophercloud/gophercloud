package backups

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("backups")
}

func resourceURL(c *gophercloud.ServiceClient, backupID string) string {
	return c.ServiceURL("backups", backupID)
}

func instanceURL(c *gophercloud.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "backups")
}
