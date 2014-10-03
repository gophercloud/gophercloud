package endpoints

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.CommonResult
}

// Extract interprets a GetResult, CreateResult or UpdateResult as a concrete Endpoint.
// An error is returned if the original call or the extraction failed.
func (r commonResult) Extract() (*Endpoint, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Endpoint `json:"endpoint"`
	}

	err := mapstructure.Decode(r.Resp, &res)
	if err != nil {
		return nil, fmt.Errorf("Error decoding Endpoint: %v", err)
	}

	return &res.Endpoint, nil
}

// CreateResult is the deferred result of a Create call.
type CreateResult struct {
	commonResult
}

// createErr quickly wraps an error in a CreateResult.
func createErr(err error) CreateResult {
	return CreateResult{commonResult{gophercloud.CommonResult{Err: err}}}
}

// UpdateResult is the deferred result of an Update call.
type UpdateResult struct {
	commonResult
}

// Endpoint describes the entry point for another service's API.
type Endpoint struct {
	ID           string                   `mapstructure:"id" json:"id"`
	Availability gophercloud.Availability `mapstructure:"interface" json:"interface"`
	Name         string                   `mapstructure:"name" json:"name"`
	Region       string                   `mapstructure:"region" json:"region"`
	ServiceID    string                   `mapstructure:"service_id" json:"service_id"`
	URL          string                   `mapstructure:"url" json:"url"`
}

// EndpointPage is a single page of Endpoint results.
type EndpointPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if no Endpoints were returned.
func (p EndpointPage) IsEmpty() (bool, error) {
	es, err := ExtractEndpoints(p)
	if err != nil {
		return true, err
	}
	return len(es) == 0, nil
}

// ExtractEndpoints extracts an Endpoint slice from a Page.
func ExtractEndpoints(page pagination.Page) ([]Endpoint, error) {
	var response struct {
		Endpoints []Endpoint `mapstructure:"endpoints"`
	}

	err := mapstructure.Decode(page.(EndpointPage).Body, &response)
	if err != nil {
		return nil, err
	}
	return response.Endpoints, nil
}
