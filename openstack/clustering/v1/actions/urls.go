package actions

import "github.com/gophercloud/gophercloud"

var apiVersion = "v1"
var apiName = "actions"

func commonURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiVersion, apiName)
}

func listURL(client *gophercloud.ServiceClient) string {
	return commonURL(client)
}
