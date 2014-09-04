package endpoints

import "github.com/rackspace/gophercloud"

// Endpoint describes the entry point for another service's API.
type Endpoint struct {
	ID        string    `json:"id"`
	Interface Interface `json:"interface"`
	Name      string    `json:"name"`
	Region    string    `json:"region"`
	ServiceID string    `json:"service_id"`
	URL       string    `json:"url"`
}

// EndpointList contains a page of Endpoint results.
type EndpointList struct {
	gophercloud.Pagination
	Endpoints []Endpoint
}
