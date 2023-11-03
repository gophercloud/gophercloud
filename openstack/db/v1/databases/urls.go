package databases

import "github.com/gophercloud/gophercloud"

func baseURL(c *gophercloud.ServiceClient, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "databases")
}

func dbURL(c *gophercloud.ServiceClient, instanceID, dbName string) string {
	return c.ServiceURL("instances", instanceID, "databases", dbName)
}

func dbGrantAccessURL(c *gophercloud.ServiceClient, instanceID, userName string) string {
	return c.ServiceURL("instances", instanceID, "users", userName, "databases")
}

func dbRevokeAccessURL(c *gophercloud.ServiceClient, instanceID, userName, dbName string) string {
	return c.ServiceURL("instances", instanceID, "users", userName, "databases", dbName)
}
