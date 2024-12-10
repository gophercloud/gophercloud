package zones

import "github.com/gophercloud/gophercloud/v2"

func baseURL(c gophercloud.Client) string {
	return c.ServiceURL("zones")
}

func zoneURL(c gophercloud.Client, zoneID string) string {
	return c.ServiceURL("zones", zoneID)
}
