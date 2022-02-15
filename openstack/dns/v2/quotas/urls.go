package quotas

import "github.com/gophercloud/gophercloud/v2"

func URL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL("quotas", projectID)
}
