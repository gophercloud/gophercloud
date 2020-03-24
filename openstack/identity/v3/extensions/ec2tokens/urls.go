package ec2tokens

import "github.com/gophercloud/gophercloud"

func ec2tokensURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("ec2tokens")
}
