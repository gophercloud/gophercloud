package services

import (
	"github.com/gophercloud/gophercloud"
)

// Service is a VPN Service
type Service struct {
	//The ID of the project
	TenantID     string `json:"tenant_id"`
	//The ID of the subnet
	SubnetID     string `json:"subnet_id"`
	//The ID of the router
	RouterID     string `json:"router_id"`
	//A human-readable description for the resource
	//Default is an empty string
	Description  string `json:"description"`
	//The administrative state of the resource, which is up (true) or down (false).
	AdminStateUp *bool  `json:"admin_state_up"`
	//The ID of the project
	ProjectID    string `json:"project_id"`
	//The human readable name of the service
	Name         string `json:"name"`
	//Indicates whether IPsec VPN service is currently operational
	//Values are ACTIVE, DOWN, BUILD, ERROR, PENDING_CREATE, PENDING_UPDATE, or PENDING_DELETE.
	Status       string `json:"status"`
	//The unique ID of the VPN service
	ID           string `json:"id"`
	//Read-only external (public) IPv6 address that is used for the VPN service
	ExternalV6IP string `json:"external_v6_ip"`
	//Read-only external (public) IPv4 address that is used for the VPN service
	ExternalV4IP string `json:"external_v4_ip"`
	//The ID of the flavor
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
