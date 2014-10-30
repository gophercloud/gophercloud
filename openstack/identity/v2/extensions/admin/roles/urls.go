package roles

import "github.com/rackspace/gophercloud"

const extPath = "OS-KSADMN/roles"

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(extPath, id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(extPath)
}

func userRoleURL(c *gophercloud.ServiceClient, tenantID, userID, roleID string) string {
	return c.ServiceURL("tenants", tenantID, "users", userID, extPath, roleID)
}
