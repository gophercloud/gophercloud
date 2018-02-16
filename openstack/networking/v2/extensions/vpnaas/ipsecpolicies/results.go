package ipsecpolicies

import (
	"github.com/gophercloud/gophercloud"
)

// Policy is an IPSec Policy
type Policy struct {
	// The ID of the project
	TenantID string `json:"tenant_id"`
	// The human readable description of the policy
	Description string `json:"description"`
	// The human readable name of the policy
	Name string `json:"name"`
	// The authentication hash algorithm
	AuthAlgorithm string `json:"auth_algorithm"`
	// The encapsulation mode
	EncapsulationMode string `json:"encapsulation_mode"`
	// The encryption algorithm
	EncryptionAlgorithm string `json:"encryption_algorithm"`
	// The Perfect forward secrecy (PFS) mode
	PFS string `json:"pfs"`
	// The transform protocol
	TransformProtocol string `json:"transform_protocol"`
	// The lifetime of the security association
	Lifetime *Lifetime `json:"lifetime"`
}

type Lifetime struct {
	// The unit for the lifetime
	// Default is seconds
	LifetimeUnit string `json:"units"`
	// The lifetime
	// Default is 3600
	LifetimeValue int `json:"value"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts an IPSec Policy.
func (r commonResult) Extract() (*Policy, error) {
	var s struct {
		Policy *Policy `json:"ipsecpolicy"`
	}
	err := r.ExtractInto(&s)
	return s.Policy, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Policy.
type CreateResult struct {
	commonResult
}
