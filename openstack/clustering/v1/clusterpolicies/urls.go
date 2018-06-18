package clusterpolicies

import "github.com/gophercloud/gophercloud"

var apiVersion = "v1"
var apiName = "clusters"

func getURL(client *gophercloud.ServiceClient, clusterID string, policyID string) string {
	return client.ServiceURL(apiVersion, apiName, clusterID, "policies", policyID)
}
