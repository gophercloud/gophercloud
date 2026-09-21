package datastores

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("datastores")
}

func createVersionURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("mgmt", "datastore-versions")
}

func mgmtVersionURL(c *gophercloud.ServiceClient, versionID string) string {
	return c.ServiceURL("mgmt", "datastore-versions", versionID)
}

func resourceURL(c *gophercloud.ServiceClient, dsID string) string {
	return c.ServiceURL("datastores", dsID)
}

func versionsURL(c *gophercloud.ServiceClient, dsID string) string {
	return c.ServiceURL("datastores", dsID, "versions")
}

func versionURL(c *gophercloud.ServiceClient, dsID, versionID string) string {
	return c.ServiceURL("datastores", dsID, "versions", versionID)
}
