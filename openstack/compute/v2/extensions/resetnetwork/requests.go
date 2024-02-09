package resetnetwork

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions"
)

// ResetNetwork will reset the network of a server
func ResetNetwork(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ResetResult) {
	b := map[string]interface{}{
		"resetNetwork": nil,
	}
	resp, err := client.PostWithContext(ctx, extensions.ActionURL(client, id), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
