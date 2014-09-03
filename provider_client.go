package gophercloud

// ProviderClient stores details that are required to interact with any services within a specific provider's API.
//
// Generally, you acquire a ProviderClient by calling the `NewClient()` method in the appropriate provider's child package,
// providing whatever authentication credentials are required.
type ProviderClient struct {

	// Options remembers the original authentication parameters, if reauthentication is enabled.
	Options AuthOptions

	// TokenID is the most recently valid token issued.
	TokenID string
}

// AuthenticatedHeaders returns a map of HTTP headers that are common for all authenticated service
// requests.
func (client *ProviderClient) AuthenticatedHeaders() map[string]string {
	return map[string]string{"X-Auth-Token": client.TokenID}
}
