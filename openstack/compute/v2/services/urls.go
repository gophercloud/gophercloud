package services

import "github.com/gophercloud/gophercloud/v2"

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("os-services")
}

func updateURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("os-services", id)
}
