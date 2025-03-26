package messages

import (
	"github.com/gophercloud/gophercloud/v2"
)

const (
	apiVersion = "v2"
	apiName    = "queues"
)

func createURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "messages")
}

func listURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "messages")
}

func getURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "messages")
}

func deleteURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "messages")
}

func DeleteMessageURL(client *gophercloud.ServiceClient, queueName string, messageID string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "messages", messageID)
}

func messageURL(client *gophercloud.ServiceClient, queueName string, messageID string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "messages", messageID)
}
