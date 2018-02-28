package endpointgroups

import (
	"github.com/gophercloud/gophercloud"
)

// EndpointGroup is an endpoint group.
type EndpointGroup struct {
	// TenantID specifies a tenant to own the endpoint group.
	TenantID string `json:"tenant_id"`

	// TenantID specifies a tenant to own the endpoint group.
	ProjectID string `json:"project_id"`

	// Description is the human readable description of the endpoint group.
	Description string `json:"description"`

	// Name is the human readable name of the endpoint group.
	Name string `json:"name"`

	// Type is the type of the endpoints in the group.
	Type string `json:"type"`

	// Endpoints is a list of endpoints.
	Endpoints []string `json:"endpoints"`

	// ID is the id of the endpoint group
	ID string `json:"id"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts an endpoint group.
func (r commonResult) Extract() (*EndpointGroup, error) {
	var s struct {
		Service *EndpointGroup `json:"endpoint_group"`
	}
	err := r.ExtractInto(&s)
	return s.Service, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as an endpoint group.
type CreateResult struct {
	commonResult
}
