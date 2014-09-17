package images

import "github.com/rackspace/gophercloud"

func getListURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("images", "detail")
}
