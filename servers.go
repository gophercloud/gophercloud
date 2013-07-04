package gophercloud

// genericCloudProvider structures provide the implementation for generic OpenStack-compatible
// ComputeProvider interfaces.
type genericCloudProvider struct {
	// endpoint refers to the provider's API endpoint base URL.  This will be used to construct
	// and issue queries.
	endpoint string

	// Test context (if any) in which to issue requests.
	context *Context
}

// See the ComputeProvider interface for details.
func (gcp *genericCloudProvider) ListServers() ([]Server, error) {
	return nil, nil
}

// Server structures provide data about a server running in your provider's cloud.
type Server struct {
	Id string
}
