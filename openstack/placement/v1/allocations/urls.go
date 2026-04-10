package allocations

import "github.com/gophercloud/gophercloud/v2"

func getURL(client *gophercloud.ServiceClient, consumerUUID string) string {
	return client.ServiceURL("allocations", consumerUUID)
}

func updateURL(client *gophercloud.ServiceClient, consumerUUID string) string {
	return client.ServiceURL("allocations", consumerUUID)
}
