package queues

import "github.com/gophercloud/gophercloud"

var apiVersion = "v2"
var apiName = "queues"

func createURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName)
}
