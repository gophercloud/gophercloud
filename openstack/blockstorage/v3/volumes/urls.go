package volumes

import "github.com/gophercloud/gophercloud/v2"

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("volumes")
}

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("volumes", "detail")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("volumes", id)
}

func getURL(c gophercloud.Client, id string) string {
	return deleteURL(c, id)
}

func updateURL(c gophercloud.Client, id string) string {
	return deleteURL(c, id)
}

func actionURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("volumes", id, "action")
}
