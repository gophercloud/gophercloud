package apiversions

import (
	"github.com/gophercloud/gophercloud/v2"
)

func getURL(c gophercloud.Client, version string) string {
	return c.ServiceURL(version)
}

func listURL(c gophercloud.Client) string {
	return c.ServiceURL()
}
