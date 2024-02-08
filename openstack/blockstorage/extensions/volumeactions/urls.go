package volumeactions

import "github.com/gophercloud/gophercloud/v2"

func actionURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("volumes", id, "action")
}
