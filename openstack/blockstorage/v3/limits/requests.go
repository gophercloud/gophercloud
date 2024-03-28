package limits

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Get returns the limits about the currently scoped tenant.
func Get(ctx context.Context, client *gophercloud.ServiceClient) (r GetResult) {
	url := getURL(client)
	resp, err := client.Get(ctx, url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
