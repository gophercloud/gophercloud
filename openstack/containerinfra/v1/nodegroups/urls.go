package nodegroups

import (
	"github.com/gophercloud/gophercloud"
)

func getURL(c *gophercloud.ServiceClient, clusterID, nodeGroupID string) string {
	return c.ServiceURL("clusters", clusterID, "nodegroups", nodeGroupID)
}
