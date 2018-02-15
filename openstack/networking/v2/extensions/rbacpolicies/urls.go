package rbacpolicies

import "github.com/gophercloud/gophercloud"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("rbac-policies")
}

func createURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}
