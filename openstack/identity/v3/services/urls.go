package services

import "github.com/rackspace/gophercloud"

func getListURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("services")
}

func getServiceURL(client *gophercloud.ServiceClient, serviceID string) string {
	return ""
}
