package gophercloud

import "errors"

var (
	// ErrEndpointNotFound is returned when no available endpoints match the provided EndpointOpts.
	ErrEndpointNotFound = errors.New("No suitable endpoint could be found in the service catalog.")
)

// EndpointOpts contains options for finding an endpoint for an Openstack client.
type EndpointOpts struct {

	// Type is the service type for the client (e.g., "compute", "object-store").
	// Type is a required field.
	Type string

	// Name is the service name for the client (e.g., "nova").
	// Name is not a required field, but it is used if present.
	// Services can have the same Type but a different Name, which is one example of when both Type and Name are needed.
	Name string

	// Region is the region in which the service resides.
	Region string

	// URLType is they type of endpoint to be returned (e.g., "public", "private").
	// URLType is not required, and defaults to "public".
	URLType string
}

// EndpointLocator is a function that describes how to locate a single endpoint from a service catalog for a specific ProviderClient.
// It should be set during ProviderClient initialization and used to discover related ServiceClients.
type EndpointLocator func(EndpointOpts) (string, error)
