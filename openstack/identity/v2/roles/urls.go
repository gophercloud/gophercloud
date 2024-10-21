package roles

import "github.com/gophercloud/gophercloud/v2"

const (
	ExtPath  = "OS-KSADM"
	RolePath = "roles"
	UserPath = "users"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(ExtPath, RolePath)
}

func userTenantRoleURL(c *gophercloud.ServiceClient, tenantID, userID, roleID string) string {
	return c.ServiceURL("tenants", tenantID, UserPath, userID, RolePath, ExtPath, roleID)
}
