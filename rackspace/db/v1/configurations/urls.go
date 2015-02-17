package configurations

import "github.com/rackspace/gophercloud"

func baseURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("configurations")
}

func resourceURL(c *gophercloud.ServiceClient, configID string) string {
	return c.ServiceURL("configurations", configID)
}

func instancesURL(c *gophercloud.ServiceClient, configID string) string {
	return c.ServiceURL("configurations", configID, "instances")
}
