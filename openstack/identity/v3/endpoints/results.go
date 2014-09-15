package endpoints

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
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

// ExtractEndpoints extracts an Endpoint slice from a Page.
func ExtractEndpoints(page gophercloud.Page) ([]Endpoint, error) {
	var response struct {
		Endpoints []Endpoint `mapstructure:"endpoints"`
	}

	err := mapstructure.Decode(page.(gophercloud.LinkedPage).Body, &response)
	if err != nil {
		return nil, err
	}
	return response.Endpoints, nil
}
