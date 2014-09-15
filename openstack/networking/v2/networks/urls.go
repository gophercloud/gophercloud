package networks

import (
	"strings"

	"github.com/rackspace/gophercloud"
)

const Version = "v2.0"

func APIVersionsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("")
}

func APIInfoURL(c *gophercloud.ServiceClient, version string) string {
	return c.ServiceURL(strings.TrimRight(version, "/") + "/")
}

func ExtensionURL(c *gophercloud.ServiceClient, name string) string {
	return c.ServiceURL(Version, "extensions", name)
}

func NetworkURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(Version, "networks", id)
}

func CreateURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(Version, "networks")
}
