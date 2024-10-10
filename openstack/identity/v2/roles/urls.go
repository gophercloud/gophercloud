package roles

import "github.com/gophercloud/gophercloud/v2"

const (
	ExtPath  = "OS-KSADM"
	RolePath = "roles"
	UserPath = "users"
)

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(ExtPath, RolePath)
}

func userTenantRoleURL(c gophercloud.Client, tenantID, userID, roleID string) string {
	return c.ServiceURL("tenants", tenantID, UserPath, userID, RolePath, ExtPath, roleID)
}
