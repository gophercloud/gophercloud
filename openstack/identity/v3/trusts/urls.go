package trusts

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "OS-TRUST/trusts"

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func createURL(c gophercloud.Client) string {
	return rootURL(c)
}

func deleteURL(c gophercloud.Client, id string) string {
	return resourceURL(c, id)
}

func listURL(c gophercloud.Client) string {
	return c.ServiceURL(resourcePath)
}

func listRolesURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(resourcePath, id, "roles")
}

func getRoleURL(c gophercloud.Client, id, roleID string) string {
	return c.ServiceURL(resourcePath, id, "roles", roleID)
}
