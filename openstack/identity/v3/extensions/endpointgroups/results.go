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

// EndpointGroup represents a group of endpoints matching one or several filter
// criteria.
type EndpointGroup struct {
	// ID is the unique ID of the endpoint group.
	ID string `json:"id"`

	// Name is the name of the new endpoint group.
	Name string `json:"name"`

	// Filters is a map type describing the filter criteria
	Filters map[string]interface{} `json:"filters"`

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
