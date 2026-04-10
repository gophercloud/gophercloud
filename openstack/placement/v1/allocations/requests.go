package allocations

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Get retrieves the allocations for a specific consumer by its UUID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, consumerUUID string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, consumerUUID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
