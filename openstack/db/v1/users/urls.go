package users

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c gophercloud.Client, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "users")
}

func userURL(c gophercloud.Client, instanceID, userName string) string {
	return c.ServiceURL("instances", instanceID, "users", userName)
}
