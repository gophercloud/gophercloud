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
	// The units for the lifetime of the security association
	// The lifetime consists of a unit and integer value
	LifetimeUnit  string `json:"lifetime.units"`
	LifetimeValue int    `json:"lifetime.value"`
}

type commonResult struct {
	gophercloud.Result
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Policy.
type CreateResult struct {
	commonResult
}
