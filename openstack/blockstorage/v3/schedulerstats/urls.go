package schedulerstats

import "github.com/gophercloud/gophercloud/v2"

func storagePoolsListURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("scheduler-stats", "get_pools")
}
