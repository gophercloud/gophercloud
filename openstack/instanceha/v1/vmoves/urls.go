package vmoves

import "github.com/gophercloud/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient, notificationID string) string {
	return client.ServiceURL("notifications", notificationID, "vmoves")
}
func getURL(client *gophercloud.ServiceClient, notificationID string, id string) string {
	return client.ServiceURL("notifications", notificationID, "vmoves", id)
}
