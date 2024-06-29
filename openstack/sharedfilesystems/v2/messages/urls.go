package messages

import "github.com/gophercloud/gophercloud/v2"

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("messages")
}

func getURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("messages", id)
}

func deleteURL(c gophercloud.Client, id string) string {
	return getURL(c, id)
}
