package services

import (
	"github.com/gophercloud/gophercloud"
)

// Service is a VPN Service
type Service struct {
	// TenantID is the ID of the project.
	TenantID string `json:"tenant_id"`

	// SubnetID is the ID of the subnet.
	SubnetID string `json:"subnet_id"`

	// RouterID is the ID of the router.
	RouterID string `json:"router_id"`

	// Description is a human-readable description for the resource.
	// Default is an empty string
	Description string `json:"description"`

	// AdminStateUp is the administrative state of the resource, which is up (true) or down (false).
	AdminStateUp bool `json:"admin_state_up"`

	// Name is the human readable name of the service.
	Name string `json:"name"`

	// Status indicates whether IPsec VPN service is currently operational.
	// Values are ACTIVE, DOWN, BUILD, ERROR, PENDING_CREATE, PENDING_UPDATE, or PENDING_DELETE.
	Status string `json:"status"`

	// ID is the unique ID of the VPN service.
	ID string `json:"id"`

	// ExternalV6IP is the read-only external (public) IPv6 address that is used for the VPN service.
	ExternalV6IP string `json:"external_v6_ip"`

	// ExternalV4IP is the read-only external (public) IPv4 address that is used for the VPN service.
	ExternalV4IP string `json:"external_v4_ip"`

	// FlavorID is the ID of the flavor.
	FlavorID string `json:"flavor_id"`
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

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the operation succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
