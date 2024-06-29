package buildinfo

import "github.com/gophercloud/gophercloud/v2"

func getURL(c gophercloud.Client) string {
	return c.ServiceURL("build_info")
}
