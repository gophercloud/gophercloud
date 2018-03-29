package policytypes

import "github.com/gophercloud/gophercloud"

var apiVersion = "v1"
var apiName = "policy-types"

func policyTypeListURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiVersion, apiName)
}
