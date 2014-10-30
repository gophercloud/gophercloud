package roles

import "github.com/rackspace/gophercloud"

const (
	ExtPath  = "OS-KSADMN/roles"
	UserPath = "users"
)

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(ExtPath, id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(ExtPath)
}

func userRoleURL(c *gophercloud.ServiceClient, tenantID, userID, roleID string) string {
	return c.ServiceURL("tenants", tenantID, UserPath, userID, ExtPath, roleID)
}
