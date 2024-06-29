package nodegroups

import (
	"github.com/gophercloud/gophercloud/v2"
)

func getURL(c gophercloud.Client, clusterID, nodeGroupID string) string {
	return c.ServiceURL("clusters", clusterID, "nodegroups", nodeGroupID)
}

func listURL(c gophercloud.Client, clusterID string) string {
	return c.ServiceURL("clusters", clusterID, "nodegroups")
}

func createURL(c gophercloud.Client, clusterID string) string {
	return c.ServiceURL("clusters", clusterID, "nodegroups")
}

func updateURL(c gophercloud.Client, clusterID, nodeGroupID string) string {
	return c.ServiceURL("clusters", clusterID, "nodegroups", nodeGroupID)
}

func deleteURL(c gophercloud.Client, clusterID, nodeGroupID string) string {
	return c.ServiceURL("clusters", clusterID, "nodegroups", nodeGroupID)
}
