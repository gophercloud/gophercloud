package endpointgroups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract interprets a GetResult, CreateResult or UpdateResult as a concrete
// EndpointGroup. An error is returned if the original call or the extraction
// failed.
func (r commonResult) Extract() (*EndpointGroup, error) {
	var s struct {
		EndpointGroup *EndpointGroup `json:"endpoint_group"`
	}
	err := r.ExtractInto(&s)
	return s.EndpointGroup, err
}

// GetResult is the response from a Get operation. Call its Extract method
// to interpret it as an EndpointGroup.
type GetResult struct {
	commonResult
}

// EndpointFilter represents a set of one or several criteria to match endpoints.
type EndpointFilter struct {
	// Availability is the interface type to filter (admin, internal,
	// or public), referenced by the gophercloud.Availability type.
	Availability gophercloud.Availability `json:"interface,omitempty"`

	// ServiceID is the ID of the service to filter
	ServiceID string `json:"service_id,omitempty"`

	// RegionID is the ID of the region to filter
	RegionID string `json:"region_id,omitempty"`
}

// EndpointGroup represents a group of endpoints matching one or several filter
// criteria.
type EndpointGroup struct {
	// ID is the unique ID of the endpoint group.
	ID string `json:"id"`

	// Name is the name of the new endpoint group.
	Name string `json:"name"`

	// Filters is an EndpointFilter type describing the filter criteria
	Filters EndpointFilter `json:"filters"`

	// Description is the description of the endpoint group
	Description string `json:"description"`
}

// EndpointGroupPage is a single page of EndpointGroup results.
type EndpointGroupPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if no EndpointGroups were returned.
func (r EndpointGroupPage) IsEmpty() (bool, error) {
	es, err := ExtractEndpointGroups(r)
	return len(es) == 0, err
}

// ExtractEndpointGroups extracts an EndpointGroup slice from a Page.
func ExtractEndpointGroups(r pagination.Page) ([]EndpointGroup, error) {
	var s struct {
		EndpointGroups []EndpointGroup `json:"endpoint_groups"`
	}
	err := (r.(EndpointGroupPage)).ExtractInto(&s)
	return s.EndpointGroups, err
}
