package buildinfo

import "github.com/gophercloud/gophercloud"

var apiVersion = "v1"
var apiName = "build-info"

func commonURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiVersion, apiName)
}

func getURL(client *gophercloud.ServiceClient) string {
	return commonURL(client)
}
