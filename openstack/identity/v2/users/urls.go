package users

import "github.com/gophercloud/gophercloud/v2"

const (
	tenantPath = "tenants"
	userPath   = "users"
	rolePath   = "roles"
)

func ResourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(userPath, id)
}

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(userPath)
}

func listRolesURL(c gophercloud.Client, tenantID, userID string) string {
	return c.ServiceURL(tenantPath, tenantID, userPath, userID, rolePath)
}
