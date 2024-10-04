package services

import "github.com/gophercloud/gophercloud/v2"

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("os-services")
}
