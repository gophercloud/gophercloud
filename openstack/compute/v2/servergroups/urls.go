package servergroups

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "os-server-groups"

func resourceURL(c gophercloud.Client) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c gophercloud.Client) string {
	return resourceURL(c)
}

func createURL(c gophercloud.Client) string {
	return resourceURL(c)
}

func getURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func deleteURL(c gophercloud.Client, id string) string {
	return getURL(c, id)
}
