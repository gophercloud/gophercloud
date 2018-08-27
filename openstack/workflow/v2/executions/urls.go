package executions

import "github.com/gophercloud/gophercloud"

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("executions")
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("executions", id)
}
