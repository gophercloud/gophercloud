package ipsecpolicies

import (
	"github.com/gophercloud/gophercloud"
)

// Policy is an IPSec Policy
type Policy struct {
	// TenantID is the ID of the project
	TenantID string `json:"tenant_id"`

	// ProjectID is the ID of the project
	ProjectID string `json:"project_id"`

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

	// TransformProtocol is the transform protocol
	TransformProtocol string `json:"transform_protocol"`

	// Lifetime is the lifetime of the security association
	Lifetime Lifetime `json:"lifetime"`

	// ID is the ID of the policy
	ID string `json:"id"`
}

type Lifetime struct {
	// Units is the unit for the lifetime
	// Default is seconds
	Units string `json:"units"`

	// Value is the lifetime
	// Default is 3600
	Value int `json:"value"`
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

// CreateResult represents the result of a delete operation. Call its ExtractErr method
// to determine if the operation succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Policy.
type GetResult struct {
	commonResult
}
