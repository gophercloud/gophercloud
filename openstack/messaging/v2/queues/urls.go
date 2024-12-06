package queues

import (
	"net/url"

	"github.com/gophercloud/gophercloud/v2"
)

const (
	apiVersion = "v2"
	apiName    = "queues"
)

func commonURL(client gophercloud.Client) string {
	return client.ServiceURL(apiVersion, apiName)
}

func createURL(client gophercloud.Client, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName)
}

func listURL(client gophercloud.Client) string {
	return commonURL(client)
}

func updateURL(client gophercloud.Client, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName)
}

func getURL(client gophercloud.Client, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName)
}

func deleteURL(client gophercloud.Client, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName)
}

func statURL(client gophercloud.Client, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "stats")
}

func shareURL(client gophercloud.Client, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "share")
}

func purgeURL(client gophercloud.Client, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "purge")
}

// builds next page full url based on current url
func nextPageURL(currentURL string, next string) (string, error) {
	base, err := url.Parse(currentURL)
	if err != nil {
		return "", err
	}
	rel, err := url.Parse(next)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(rel).String(), nil
}
