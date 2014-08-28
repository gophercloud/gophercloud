package gophercloud

import (
	"strings"
)

// ServiceClient stores details about a specific service that are necessary for further interactions with that service API, as well
// as utility methods for service implementation.
//
// Generally, you will acquire a ServiceClient by calling the NewClient() function in the appropriate service package.
type ServiceClient struct {

	// Authority caches results of the most recent authentication.
	Authority AuthResults

	// Options remembers the original authentication parameters, if reauthentication is enabled.
	Options AuthOptions

	// Endpoint is the base URL of the relevant API.
	Endpoint string

	// TokenID is the most recently valid token issued.
	TokenID string
}

// ServiceURL constructs a URL for a resource belonging to this client.
func (client *ServiceClient) ServiceURL(parts ...string) string {
	return client.Endpoint + strings.Join(parts, "/")
}

// AuthenticatedHeaders returns a map of HTTP headers that are common for all authenticated service
// requests.
func (client *ServiceClient) AuthenticatedHeaders() map[string]string {
	return map[string]string{"X-Auth-Token": client.TokenID}
}
