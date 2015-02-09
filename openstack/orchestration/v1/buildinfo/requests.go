package buildinfo

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// Get retreives data for the given stack template.
func Get(c *gophercloud.ServiceClient) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", getURL(c), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}
