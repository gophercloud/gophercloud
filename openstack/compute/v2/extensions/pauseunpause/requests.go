package pauseunpause

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions"
)

// Pause is the operation responsible for pausing a Compute server.
func Pause(ctx context.Context, client *gophercloud.ServiceClient, id string) (r PauseResult) {
	resp, err := client.PostWithContext(ctx, extensions.ActionURL(client, id), map[string]interface{}{"pause": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Unpause is the operation responsible for unpausing a Compute server.
func Unpause(ctx context.Context, client *gophercloud.ServiceClient, id string) (r UnpauseResult) {
	resp, err := client.PostWithContext(ctx, extensions.ActionURL(client, id), map[string]interface{}{"unpause": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
