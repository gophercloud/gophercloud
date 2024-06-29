package containers

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("containers")
}

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("containers", id)
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("containers")
}

func deleteURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("containers", id)
}

func listConsumersURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("containers", id, "consumers")
}

func createConsumerURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("containers", id, "consumers")
}

func deleteConsumerURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("containers", id, "consumers")
}

func createSecretRefURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("containers", id, "secrets")
}

func deleteSecretRefURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("containers", id, "secrets")
}
