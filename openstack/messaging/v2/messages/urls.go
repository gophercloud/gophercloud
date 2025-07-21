package messages

import (
	"net/url"

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

// builds next page full url based on service endpoint
func nextPageURL(endpointURL, next string) (string, error) {
	base, err := url.Parse(endpointURL)
	if err != nil {
		return "", err
	}
	rel, err := url.Parse(next)
	if err != nil {
		return "", err
	}
	combined := base.JoinPath(rel.Path)
	combined.RawQuery = rel.RawQuery
	return combined.String(), nil
}
