package apiversions

import (
	"strings"

	"github.com/rackspace/gophercloud"
)

func apiVersionsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("")
}

func apiInfoURL(c *gophercloud.ServiceClient, version string) string {
	return c.ServiceURL(strings.TrimRight(version, "/") + "/")
}
