package services

import (
	"github.com/gophercloud/gophercloud"
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

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Service.
type CreateResult struct {
	commonResult
}
