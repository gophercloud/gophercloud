package subnetpools

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "subnetpools"

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(resourcePath)
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
