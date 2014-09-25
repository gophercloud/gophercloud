package images

import "github.com/rackspace/gophercloud"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("images", "detail")
}

func imageURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("images", id)
}
