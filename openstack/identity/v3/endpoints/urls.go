package endpoints

import "github.com/rackspace/gophercloud"

func getListURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("endpoints")
}

func getEndpointURL(client *gophercloud.ServiceClient, endpointID string) string {
	return client.ServiceURL("endpoints", endpointID)
}
