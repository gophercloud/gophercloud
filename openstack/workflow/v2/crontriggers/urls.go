package crontriggers

import "github.com/gophercloud/gophercloud/v2"

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("cron_triggers")
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("cron_triggers", id)
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("cron_triggers", id)
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("cron_triggers")
}
