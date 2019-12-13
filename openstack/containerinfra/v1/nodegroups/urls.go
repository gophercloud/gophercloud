package nodegroups

import (
	"github.com/gophercloud/gophercloud"
)

func getURL(c *gophercloud.ServiceClient, clusterID, nodeGroupID string) string {
	return c.ServiceURL("clusters", clusterID, "nodegroups", nodeGroupID)
}

func listURL(c *gophercloud.ServiceClient, clusterID string) string {
	return c.ServiceURL("clusters", clusterID, "nodegroups")
}

func createURL(c *gophercloud.ServiceClient, clusterID string) string {
	return c.ServiceURL("clusters", clusterID, "nodegroups")
}

func updateURL(c *gophercloud.ServiceClient, clusterID, nodeGroupID string) string {
	return c.ServiceURL("clusters", clusterID, "nodegroups", nodeGroupID)
}

func deleteURL(c *gophercloud.ServiceClient, clusterID, nodeGroupID string) string {
	return c.ServiceURL("clusters", clusterID, "nodegroups", nodeGroupID)
}
