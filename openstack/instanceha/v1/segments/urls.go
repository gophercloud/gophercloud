package segments

import "github.com/gophercloud/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient) string { return client.ServiceURL("segments") }
func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("segments", id)
}
