package clusterpolicies

import "github.com/gophercloud/gophercloud"

var apiVersion = "v1"
var apiName = "clusters"

func listURL(client *gophercloud.ServiceClient, clusterID string) string {
	return client.ServiceURL(apiVersion, apiName, clusterID, "policies")
}
