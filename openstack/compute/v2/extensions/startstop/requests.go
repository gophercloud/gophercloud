package startstop

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions"
)

// Start is the operation responsible for starting a Compute server.
func Start(ctx context.Context, client *gophercloud.ServiceClient, id string) (r StartResult) {
	resp, err := client.PostWithContext(ctx, extensions.ActionURL(client, id), map[string]interface{}{"os-start": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Stop is the operation responsible for stopping a Compute server.
func Stop(ctx context.Context, client *gophercloud.ServiceClient, id string) (r StopResult) {
	resp, err := client.PostWithContext(ctx, extensions.ActionURL(client, id), map[string]interface{}{"os-stop": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
