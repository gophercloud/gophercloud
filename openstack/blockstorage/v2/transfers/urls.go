package transfers

import "github.com/gophercloud/gophercloud/v2"

func transferURL(c gophercloud.Client) string {
	return c.ServiceURL("os-volume-transfer")
}

func acceptURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("os-volume-transfer", id, "accept")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("os-volume-transfer", id)
}

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("os-volume-transfer", "detail")
}

func getURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("os-volume-transfer", id)
}
