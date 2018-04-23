package nodes

import "github.com/gophercloud/gophercloud"

var apiVersion = "v1"
var apiName = "nodes"

func idURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL(apiVersion, apiName, id)
}

func updateURL(client *gophercloud.ServiceClient, id string) string {
	return idURL(client, id)
}
