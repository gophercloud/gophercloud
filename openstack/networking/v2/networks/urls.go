package networks

import "github.com/gophercloud/gophercloud/v2"

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("networks", id)
}

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL("networks")
}

func getURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

func listURL(c gophercloud.Client) string {
	return rootURL(c)
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
