package endpoints

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("endpoints")
}

func endpointURL(client gophercloud.Client, endpointID string) string {
	return client.ServiceURL("endpoints", endpointID)
}
