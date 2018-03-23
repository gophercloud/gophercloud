package nodes

import "github.com/gophercloud/gophercloud"

var apiVersion = "v1"
var apiName = "nodes"

func commonURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiVersion, apiName)
}

func createURL(client *gophercloud.ServiceClient) string {
	return commonURL(client)
}
