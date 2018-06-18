package claims

import "github.com/gophercloud/gophercloud"

const (
	apiVersion = "v2"
	apiName    = "queues"
)

func createURL(client *gophercloud.ServiceClient, queueName string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "claims")
}

func getURL(client *gophercloud.ServiceClient, queueName string, claimId string) string {
	return client.ServiceURL(apiVersion, apiName, queueName, "claims", claimId)
}
