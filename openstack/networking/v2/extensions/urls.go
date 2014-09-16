package networks

import "github.com/rackspace/gophercloud"

const Version = "v2.0"

func ExtensionURL(c *gophercloud.ServiceClient, name string) string {
	return c.ServiceURL(Version, "extensions", name)
}
