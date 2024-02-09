package suspendresume

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions"
)

// Suspend is the operation responsible for suspending a Compute server.
func Suspend(ctx context.Context, client *gophercloud.ServiceClient, id string) (r SuspendResult) {
	resp, err := client.PostWithContext(ctx, extensions.ActionURL(client, id), map[string]interface{}{"suspend": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Resume is the operation responsible for resuming a Compute server.
func Resume(ctx context.Context, client *gophercloud.ServiceClient, id string) (r UnsuspendResult) {
	resp, err := client.PostWithContext(ctx, extensions.ActionURL(client, id), map[string]interface{}{"resume": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
