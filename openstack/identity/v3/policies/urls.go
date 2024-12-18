package policies

import "github.com/gophercloud/gophercloud/v2"

const policyPath = "policies"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL(policyPath)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL(policyPath)
}

func getURL(client gophercloud.Client, policyID string) string {
	return client.ServiceURL(policyPath, policyID)
}

func updateURL(client gophercloud.Client, policyID string) string {
	return client.ServiceURL(policyPath, policyID)
}

func deleteURL(client gophercloud.Client, policyID string) string {
	return client.ServiceURL(policyPath, policyID)
}
