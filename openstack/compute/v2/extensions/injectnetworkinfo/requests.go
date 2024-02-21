package injectnetworkinfo

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions"
)

// InjectNetworkInfo will inject the network info into a server
func InjectNetworkInfo(ctx context.Context, client *gophercloud.ServiceClient, id string) (r InjectNetworkResult) {
	b := map[string]interface{}{
		"injectNetworkInfo": nil,
	}
	resp, err := client.Post(ctx, extensions.ActionURL(client, id), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
