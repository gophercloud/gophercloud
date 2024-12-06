package services

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("services")
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("services")
}

func serviceURL(client gophercloud.Client, serviceID string) string {
	return client.ServiceURL("services", serviceID)
}

func updateURL(client gophercloud.Client, serviceID string) string {
	return client.ServiceURL("services", serviceID)
}
