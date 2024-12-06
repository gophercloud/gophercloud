package conductors

import "github.com/gophercloud/gophercloud/v2"

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("conductors")
}

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("conductors", id)
}
