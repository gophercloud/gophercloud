package images

import "github.com/rackspace/gophercloud"

func getListURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("images", "detail")
}

func getImageURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("images", id)
}
