package apiversions

import (
	"fmt"
	"net/url"

	"github.com/gophercloud/gophercloud"
)

func getURL(c *gophercloud.ServiceClient, version string) string {
	u, _ := url.Parse(c.ServiceURL(""))
	u.Path = fmt.Sprintf("/%s/", version)
	return u.String()
}

func listURL(c *gophercloud.ServiceClient) string {
	u, _ := url.Parse(c.ServiceURL(""))
	u.Path = "/"
	return u.String()
}
