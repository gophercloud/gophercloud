package networkipavailabilities

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "network-ip-availabilities"

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c gophercloud.Client, networkIPAvailabilityID string) string {
	return c.ServiceURL(resourcePath, networkIPAvailabilityID)
}

func listURL(c gophercloud.Client) string {
	return rootURL(c)
}

func getURL(c gophercloud.Client, networkIPAvailabilityID string) string {
	return resourceURL(c, networkIPAvailabilityID)
}
