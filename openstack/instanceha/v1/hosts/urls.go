package hosts

import "github.com/gophercloud/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient, segmentID string) string {
	return client.ServiceURL("segments", segmentID, "hosts")
}
func getURL(client *gophercloud.ServiceClient, segmentID string, id string) string {
	return client.ServiceURL("segments", segmentID, "hosts", id)
}
