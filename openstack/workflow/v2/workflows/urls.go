package workflows

import (
	"github.com/gophercloud/gophercloud/v2"
)

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("workflows")
}

func deleteURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("workflows", id)
}

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("workflows", id)
}

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("workflows")
}
