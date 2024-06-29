package executions

import "github.com/gophercloud/gophercloud/v2"

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("executions")
}

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("executions", id)
}

func deleteURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("executions", id)
}

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("executions")
}
