package diagnostics

import "github.com/gophercloud/gophercloud/v2"

// serverDiagnosticsURL returns the diagnostics url for a nova instance/server
func serverDiagnosticsURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("servers", id, "diagnostics")
}
