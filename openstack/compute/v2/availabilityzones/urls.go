package availabilityzones

import "github.com/gophercloud/gophercloud/v2"

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("os-availability-zone")
}

func listDetailURL(c gophercloud.Client) string {
	return c.ServiceURL("os-availability-zone", "detail")
}
