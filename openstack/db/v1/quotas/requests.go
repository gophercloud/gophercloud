package quotas

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Get retrieves the details of quotas for a specified tenant.
func Get(ctx context.Context, client *gophercloud.ServiceClient, projectID string) (r GetResult) {
	resp, err := client.Get(ctx, baseURL(client, projectID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
