package policies

import "github.com/gophercloud/gophercloud"

const policyPath = "policies"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(policyPath)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(policyPath)
}

func deleteURL(client *gophercloud.ServiceClient, policyID string) string {
	return client.ServiceURL(policyPath, policyID)
}
