package flavors

import "github.com/gophercloud/gophercloud/v2"

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("flavors", id)
}

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("flavors")
}
