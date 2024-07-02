package quotas

import "github.com/gophercloud/gophercloud"

func URL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL("quotas", projectID)
}
