package securityservices

import "github.com/gophercloud/gophercloud/v2"

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("security-services")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("security-services", id)
}

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("security-services", "detail")
}

func getURL(c gophercloud.Client, id string) string {
	return deleteURL(c, id)
}

func updateURL(c gophercloud.Client, id string) string {
	return deleteURL(c, id)
}
