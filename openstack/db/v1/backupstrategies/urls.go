package backupstrategies

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("backup_strategies")
}
