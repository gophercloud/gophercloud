package diagnostics

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Diagnostics
func Get(ctx context.Context, client *gophercloud.ServiceClient, serverId string) (r serverDiagnosticsResult) {
	resp, err := client.Get(ctx, serverDiagnosticsURL(client, serverId), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
