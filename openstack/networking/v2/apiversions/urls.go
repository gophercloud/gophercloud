package apiversions

import (
	"strings"

	"github.com/rackspace/gophercloud"
)

func APIVersionsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("")
}

func APIInfoURL(c *gophercloud.ServiceClient, version string) string {
	return c.ServiceURL(strings.TrimRight(version, "/") + "/")
}
