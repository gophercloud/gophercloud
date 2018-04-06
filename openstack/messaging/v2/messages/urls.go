package messages

import "github.com/gophercloud/gophercloud"

const ApiVersion = "v2"
const ApiName = "queues"

func createURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(ApiVersion, ApiName, queueName, "messages")
}
