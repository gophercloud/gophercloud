package endpoints

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

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
