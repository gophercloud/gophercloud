package crontriggers

import "github.com/gophercloud/gophercloud/v2"

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("cron_triggers")
}

func deleteURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("cron_triggers", id)
}

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("cron_triggers", id)
}

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("cron_triggers")
}
