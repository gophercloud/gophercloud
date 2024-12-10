package policies

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "qos/policies"

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func listURL(c gophercloud.Client) string {
	return rootURL(c)
}

func getURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

func createURL(c gophercloud.Client) string {
	return rootURL(c)
}

func updateURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}
