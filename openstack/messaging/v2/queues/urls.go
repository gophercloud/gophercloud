package queues

import "github.com/gophercloud/gophercloud"

const ApiVersion = "v2"
const ApiName = "queues"

func commonURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(ApiVersion, ApiName)
}

func createURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(ApiVersion, ApiName, queueName)
}

func listURL(client *gophercloud.ServiceClient) string {
	return commonURL(client)
}

func updateURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(ApiVersion, ApiName, queueName)
}

func getURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(ApiVersion, ApiName, queueName)
}

func deleteURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(ApiVersion, ApiName, queueName)
}
