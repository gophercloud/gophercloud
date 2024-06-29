package databases

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c gophercloud.Client, instanceID string) string {
	return c.ServiceURL("instances", instanceID, "databases")
}

func dbURL(c gophercloud.Client, instanceID, dbName string) string {
	return c.ServiceURL("instances", instanceID, "databases", dbName)
}
