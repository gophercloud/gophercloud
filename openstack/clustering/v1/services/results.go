package services

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

// ExtractServices provides access to the list of services in a page acquired from the ListDetail operation.
func ExtractServices(r pagination.Page) ([]Service, error) {
	var s struct {
		Services []Service `json:"services"`
	}
	err := (r.(ServicePage)).ExtractInto(&s)
	return s.Services, err
}

// ServicePage contains a single page of all services from a ListDetails call.
type ServicePage struct {
	pagination.LinkedPageBase
}

// Service represents a detailed service
type Service struct {
	Binary        string `json:"binary"`
	DisableReason string `json:"disabled_reason"`
	Host          string `json:"host"`
	ID            string `json:"id"`
	State         string `json:"state"`
	Status        string `json:"status"`
	Topic         string `json:"topic"`

	// TODO: Need to figure out properly parse time
	//UpdatedAt     time.Time `json:"updated_at"`
	UpdatedAt string `json:"updated_at"`
}

// IsEmpty determines if a NodePage contains any results.
func (page ServicePage) IsEmpty() (bool, error) {
	services, err := ExtractServices(page)
	return len(services) == 0, err
}
