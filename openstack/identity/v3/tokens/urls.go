package tokens

import identity "github.com/rackspace/gophercloud/openstack/identity/v3"

func getTokenURL(c *identity.Client) string {
	return c.ServiceURL("auth", "tokens")
}
