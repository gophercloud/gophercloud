package ipsecpolicies

import (
	"github.com/gophercloud/gophercloud"
)

// Policy is an IPSec Policy
type Policy struct {
	// TenantID is the ID of the project
	TenantID string `json:"tenant_id"`

	// Description is the human readable description of the policy
	Description string `json:"description"`

	// Name is the human readable name of the policy
	Name string `json:"name"`

	// AuthAlgorithm is the authentication hash algorithm
	AuthAlgorithm string `json:"auth_algorithm"`

	// EncapsulationMode is the encapsulation mode
	EncapsulationMode string `json:"encapsulation_mode"`

	// EncryptionAlgorithm is the encryption algorithm
	EncryptionAlgorithm string `json:"encryption_algorithm"`

	// PFS is the Perfect forward secrecy (PFS) mode
	PFS string `json:"pfs"`

	// TranformProtocol is the transform protocol
	TransformProtocol string `json:"transform_protocol"`

	// Lifetime is the lifetime of the security association
	Lifetime *Lifetime `json:"lifetime"`
}

type Lifetime struct {
	// LifetimeUnits is the unit for the lifetime
	// Default is seconds
	LifetimeUnits string `json:"units"`

	// LifetimeValue is the lifetime
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
