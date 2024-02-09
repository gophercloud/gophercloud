package lockunlock

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions"
)

// Lock is the operation responsible for locking a Compute server.
func Lock(ctx context.Context, client *gophercloud.ServiceClient, id string) (r LockResult) {
	resp, err := client.PostWithContext(ctx, extensions.ActionURL(client, id), map[string]interface{}{"lock": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Unlock is the operation responsible for unlocking a Compute server.
func Unlock(ctx context.Context, client *gophercloud.ServiceClient, id string) (r UnlockResult) {
	resp, err := client.PostWithContext(ctx, extensions.ActionURL(client, id), map[string]interface{}{"unlock": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
