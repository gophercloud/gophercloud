package buildinfo

import (
	"github.com/gophercloud/gophercloud"
)

// Get retrieves details of a single buildinfo. Use ExtractBuildInfo to convert its
// result into a BuildInfo.
func Get(client *gophercloud.ServiceClient) (r GetResult) {
	_, r.Err = client.Get(getURL(client), &r.Body, nil)
	return
}
