package clusterpolicies

import "github.com/gophercloud/gophercloud"

var apiVersion = "v1"
var apiName = "clusters"

func commonURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiVersion, apiName)
}

func getURL(client *gophercloud.ServiceClient, clusterID string) string {
	return client.ServiceURL(apiVersion, apiName, clusterID, "policies")
}

func policyURL(client *gophercloud.ServiceClient, clusterID string, policyID string) string {
	return client.ServiceURL(apiVersion, apiName, clusterID, "policies", policyID)
}

func listURL(client *gophercloud.ServiceClient, clusterID string) string {
	return getURL(client, clusterID)
}

func getDetailURL(client *gophercloud.ServiceClient, clusterID string, policyID string) string {
	return policyURL(client, clusterID, policyID)
}
