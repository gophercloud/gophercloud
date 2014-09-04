package endpoints

import (
	"errors"

	"github.com/rackspace/gophercloud"
)

// Interface describes the availability of a specific service endpoint.
type Interface string

const (
	// InterfaceAdmin makes an endpoint only available to administrators.
	InterfaceAdmin Interface = "admin"

	// InterfacePublic makes an endpoint available to everyone.
	InterfacePublic Interface = "public"

	// InterfaceInternal makes an endpoint only available within the cluster.
	InterfaceInternal Interface = "internal"
)

// EndpointOpts contains the subset of Endpoint attributes that should be used to create or update an Endpoint.
type EndpointOpts struct {
	Interface Interface
	Name      string
	Region    string
	URL       string
	ServiceID string
}

// Create inserts a new Endpoint into the service catalog.
func Create(client *gophercloud.ServiceClient, opts EndpointOpts) (*Endpoint, error) {
	return nil, errors.New("Not implemented")
}

// ListOpts allows finer control over the the endpoints returned by a List call.
// All fields are optional.
type ListOpts struct {
	Interface Interface
	ServiceID string
	Page      int
	PerPage   int
}

// List enumerates endpoints in a paginated collection, optionally filtered by ListOpts criteria.
func List(client *gophercloud.ServiceClient, opts ListOpts) (*EndpointList, error) {
	return nil, errors.New("Not implemented")
}

// Update changes an existing endpoint with new data.
func Update(client *gophercloud.ServiceClient, endpointID string, opts EndpointOpts) (*Endpoint, error) {
	return nil, errors.New("Not implemented")
}

// Delete removes an endpoint from the service catalog.
func Delete(client *gophercloud.ServiceClient, endpointID string) error {
	return errors.New("Not implemented")
}
