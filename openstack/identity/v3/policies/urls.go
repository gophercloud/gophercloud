package policies

import "github.com/gophercloud/gophercloud"

const policyPath = "policies"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(policyPath)
}
