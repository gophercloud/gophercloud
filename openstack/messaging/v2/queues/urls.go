package queues

import "github.com/gophercloud/gophercloud"

const ApiVersion = "v2"
const ApiName = "queues"

func commonURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(ApiVersion, ApiName)
}

func commonURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiVersion, apiName)
}

func createURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(ApiVersion, ApiName, queueName)
}

func listURL(client *gophercloud.ServiceClient) string {
	return commonURL(client)
}

func listURL(client *gophercloud.ServiceClient) string {
	return commonURL(client)
}
