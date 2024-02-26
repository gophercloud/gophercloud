package apiversions

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// List lists all the API versions available to end users.
func List(ctx context.Context, client *gophercloud.ServiceClient) (r ListResult) {
	resp, err := client.Get(ctx, listURL(client), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get will get a specific API version, specified by major ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, v string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, v), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
