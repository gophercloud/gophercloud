package shareaccessrules

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Get retrieves details about a share access rule.
func Get(ctx context.Context, client *gophercloud.ServiceClient, accessID string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, accessID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List gets all access rules of a share.
func List(ctx context.Context, client *gophercloud.ServiceClient, shareID string) (r ListResult) {
	resp, err := client.Get(ctx, listURL(client, shareID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
