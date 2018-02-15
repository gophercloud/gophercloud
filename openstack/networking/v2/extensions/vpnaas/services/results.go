package services

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Service is a VPN Service
type Service struct {
	TenantID     string `json:"tenant_id"`
	SubnetID     string `json:"subnet_id,omitempty"`
	RouterID     string `json:"router_id"`
	Description  string `json:"description,omitempty"`
	AdminStateUp *bool  `json:"admin_state_up"`
	ProjectID    string `json:"project_id"`
	Name         string `json:"name,omitempty"`
	Status       string `json:"status"`
	ID           string `json:"id"`
	ExternalV6IP string `json:"external_v6_ip"`
	ExternalV4IP string `json:"external_v4_ip"`
	FlavorID     string `json:"flavor_id"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a VPN service.
func (r commonResult) Extract() (*Service, error) {
	var s struct {
		Service *Service `json:"vpnservice"`
	}
	err := r.ExtractInto(&s)
	return s.Service, err
}

// ServicePage is the page returned by a pager when traversing over a
// collection of VPN services.
type ServicePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of VPN services has
// reached the end of a page and the pager seeks to traverse over a new one.
// In order to do this, it needs to construct the next page's URL.
func (r ServicePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"vpnservices_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ServicePage struct is empty.
func (r ServicePage) IsEmpty() (bool, error) {
	is, err := ExtractServices(r)
	return len(is) == 0, err
}

// ExtractServices accepts a Page struct, specifically a Service struct,
// and extracts the elements into a slice of Service structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractServices(r pagination.Page) ([]Service, error) {
	var s struct {
		Services []Service `json:"vpnservices"`
	}
	err := (r.(ServicePage)).ExtractInto(&s)
	return s.Services, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Service.
type GetResult struct {
	commonResult
}
