package diagnostics

import (
	"github.com/gophercloud/gophercloud/v2"
)

// Diagnostics
func Get(client *gophercloud.ServiceClient, serverId string) (r serverDiagnosticsResult) {
	resp, err := client.Get(serverDiagnosticsURL(client, serverId), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
