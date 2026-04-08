package allocationcandidates

import "github.com/gophercloud/gophercloud/v2"

const apiName = "allocation_candidates"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiName)
}
