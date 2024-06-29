package recordsets

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c gophercloud.Client, zoneID string) string {
	return c.ServiceURL("zones", zoneID, "recordsets")
}

func rrsetURL(c gophercloud.Client, zoneID string, rrsetID string) string {
	return c.ServiceURL("zones", zoneID, "recordsets", rrsetID)
}
