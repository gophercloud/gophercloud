package services

import (
	"github.com/rackspace/gophercloud"

	"github.com/mitchellh/mapstructure"
)

// Service is the result of a list or information query.
type Service struct {
	Description *string `json:"description,omitempty"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
}

// ExtractServices extracts a slice of Services from a Collection acquired from List.
func ExtractServices(page gophercloud.Page) ([]Service, error) {
	var response struct {
		Services []Service `mapstructure:"services"`
	}

	err := mapstructure.Decode(page.(gophercloud.LinkedPage).Body, &response)
	return response.Services, err
}
