package keypairs

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "os-keypairs"

func resourceURL(c gophercloud.Client) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c gophercloud.Client) string {
	return resourceURL(c)
}

func createURL(c gophercloud.Client) string {
	return resourceURL(c)
}

func getURL(c gophercloud.Client, name string) string {
	return c.ServiceURL(resourcePath, name)
}

func deleteURL(c gophercloud.Client, name string) string {
	return getURL(c, name)
}
