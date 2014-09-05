// +build acceptance

package v3

import (
	"testing"

	endpoints3 "github.com/rackspace/gophercloud/openstack/identity/v3/endpoints"
)

func TestListEndpoints(t *testing.T) {
	// Create a service client.
	serviceClient := createAuthenticatedClient(t)

	// Use the service to list all available endpoints.
	_, err := endpoints3.List(serviceClient, endpoints3.ListOpts{})
	if err != nil {
		t.Errorf("Unexpected error while listing endpoints: %v", err)
	}
}
