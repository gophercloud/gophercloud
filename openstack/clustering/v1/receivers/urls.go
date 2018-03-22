package receivers

import "github.com/gophercloud/gophercloud"

var apiVersion = "v1"
var apiName = "receivers"

func idURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL(apiVersion, apiName, id)
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return idURL(client, id)
}
