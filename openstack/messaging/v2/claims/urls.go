package claims

import "github.com/gophercloud/gophercloud/v2"

const (
	apiVersion = "v2"
	apiName    = "queues"
)

func createURL(client gophercloud.Client, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "claims")
}

func getURL(client gophercloud.Client, queueName string, claimID string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "claims", claimID)
}

func updateURL(client gophercloud.Client, queueName string, claimID string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "claims", claimID)
}

func deleteURL(client gophercloud.Client, queueName string, claimID string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "claims", claimID)
}
