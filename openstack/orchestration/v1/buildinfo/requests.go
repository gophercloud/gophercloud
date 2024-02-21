package buildinfo

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Get retreives data for the given stack template.
func Get(ctx context.Context, c *gophercloud.ServiceClient) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
