package notifications

import "github.com/gophercloud/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient) string { return client.ServiceURL("notifications") }
func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("notifications", id)
}
