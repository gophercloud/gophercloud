package traits

import "github.com/gophercloud/gophercloud/v2"

const (
	apiName = "traits"
)

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(apiName)
}

func getURL(client *gophercloud.ServiceClient, traitName string) string {
	return client.ServiceURL(apiName, traitName)
}
