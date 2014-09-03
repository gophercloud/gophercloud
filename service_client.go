package gophercloud

import "strings"

// ServiceClient stores details required to interact with a specific service API implemented by a provider.
// Generally, you'll acquire these by calling the appropriate `New` method on a ProviderClient.
type ServiceClient struct {
	// Provider is a reference to the provider that implements this service.
	Provider *ProviderClient

	// Endpoint is the base URL of the service's API, acquired from a service catalog.
	// It should NOT end with a /.
	Endpoint string
}

// ServiceURL constructs a URL for a resource belonging to this provider.
func (client *ServiceClient) ServiceURL(parts ...string) string {
	return client.Endpoint + "/" + strings.Join(parts, "/")
}
