package quotas

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL("mgmt", "quotas", projectID)
}
