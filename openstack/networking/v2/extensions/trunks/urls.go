package trunks

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "trunks"

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func createURL(c gophercloud.Client) string {
	return rootURL(c)
}

func deleteURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

func listURL(c gophercloud.Client) string {
	return rootURL(c)
}

func getURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

func updateURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

func getSubportsURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id, "get_subports")
}

func addSubportsURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id, "add_subports")
}

func removeSubportsURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id, "remove_subports")
}
