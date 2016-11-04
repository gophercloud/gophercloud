package recordsets

import "github.com/gophercloud/gophercloud"

func listURL(c *gophercloud.ServiceClient, zoneID string) string {
	return c.ServiceURL("v2/zones", zoneID, "recordsets")
}

func rrsetURL(c *gophercloud.ServiceClient, zoneID string, rrsetID string) string {
	return c.ServiceURL("v2/zones", zoneID, "recordsets", rrsetID)
}
