package gophercloud

import (
	"strings"
)

// ServiceClient stores details about a specific service that are necessary for further interactions
// with that service API, as well as utility methods for service implementation.
//
// Generally, you will acquire a ServiceClient by calling the NewClient() function in the
// appropriate service package.
type ServiceClient struct {
	authority AuthResults
	options   AuthOptions
	endpoint  string
	tokenID   string
}

// ServiceURL constructs a URL for a resource belonging to this client.
func (client *ServiceClient) ServiceURL(parts ...string) string {
	return client.endpoint + strings.Join(parts, "/")
}

// AuthenticatedHeaders returns a map of HTTP headers that are common for all authenticated service
// requests.
func (client *ServiceClient) AuthenticatedHeaders() map[string]string {
	return map[string]string{"X-Auth-Token": client.tokenID}
}
