package usages

import "github.com/gophercloud/gophercloud/v2"

func getURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("usages")
}
