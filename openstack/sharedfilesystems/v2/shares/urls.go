package shares

import (
	"strings"

	"github.com/gophercloud/gophercloud"
)

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("shares")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id)
}

func getMicroversionsURL(c *gophercloud.ServiceClient) string {
	baseURLWithoutEndingSlashes := strings.TrimRight(c.ResourceBaseURL(), "/")
	slashIndexBeforeProjectID := strings.LastIndex(baseURLWithoutEndingSlashes, "/")
	return baseURLWithoutEndingSlashes[:slashIndexBeforeProjectID] + "/"
}
