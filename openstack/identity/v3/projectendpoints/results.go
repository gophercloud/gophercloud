package projectendpoints

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CreateResult is the response from a Create operation. Call its Extract
// method to interpret it as an Endpoint.
type CreateResult struct {
	gophercloud.ErrResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Endpoint describes the entry point for another service's API.
type Endpoint struct {
	// ID is the unique ID of the endpoint.
	ID string `json:"id"`

	// Availability is the interface type of the Endpoint (admin, internal,
	// or public), referenced by the gophercloud.Availability type.
	Availability gophercloud.Availability `json:"interface"`

	// Region is the region the Endpoint is located in.
	Region string `json:"region"`

	// ServiceID is the ID of the service the Endpoint refers to.
	ServiceID string `json:"service_id"`

	// URL is the url of the Endpoint.
	URL string `json:"url"`
}

// EndpointPage is a single page of Endpoint results.
type EndpointPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if no Endpoints were returned.
func (r EndpointPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	es, err := ExtractEndpoints(r)
	return len(es) == 0, err
}

// ExtractEndpoints extracts an Endpoint slice from a Page.
func ExtractEndpoints(r pagination.Page) ([]Endpoint, error) {
	var s struct {
		Endpoints []Endpoint `json:"endpoints"`
	}
	err := (r.(EndpointPage)).ExtractInto(&s)
	return s.Endpoints, err
}
