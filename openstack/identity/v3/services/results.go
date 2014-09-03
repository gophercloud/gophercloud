package services

import "github.com/rackspace/gophercloud"

// ServiceResult is the result of a list or information query.
type ServiceResult struct {
	Description *string `json:"description,omitempty"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
}

// ServiceListResult is a paged collection of Services.
type ServiceListResult struct {
	gophercloud.Pagination

	Services []ServiceResult `json:"services"`
}
