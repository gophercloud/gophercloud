package ipsecpolicies

import (
	"github.com/gophercloud/gophercloud"
)

// Policy is an IPSec Policy
type Policy struct {
	// TenantID specifies a tenant to own the IPSec policy. The caller must have
	// an admin role in order to set this. Otherwise, this field is left unset
	// and the caller will be the owner.
	TenantID string `json:"tenant_id"`
	// The human readable description of the policy
	Description string `json:"description"`
	// The human readable name of the policy
	// Does not have to be unique
	Name string `json:"name"`
	// The authentication hash algorithm
	// Valid values are sha1, sha256, sha384, sha512
	// The default is sha1.
	AuthAlgorithm string `json:"auth_algorithm"`
	// The encapsulation mode
	// A valid value is tunnel or transport
	// Default is tunnel.
	EncapsulationMode string `json:"encapsulation_mode"`
	// The encryption algorithm
	// A valid value is 3des, aes-128, aes-192, aes-256, and so on
	// Default is aes-128
	EncryptionAlgorithm string `json:"encryption_algorithm"`
	// Perfect forward secrecy (PFS)
	// A valid value is Group2, Group5, Group14, and so on
	// Default is Group5
	PFS string `json:"pfs"`
	// The transform protocol
	// A valid value is ESP, AH, or AH- ESP
	// Default is ESP.
	TransformProtocol string `json:"transform_protocol"`
	// The units for the lifetime of the security association
	// The lifetime consists of a unit and integer value
	// You can omit either the unit or value portion of the lifetime
	// Default unit is seconds and default value is 3600.
	LifetimeUnit string `json:"lifetime.units"`
	LifetimeName string `json:"lifetime.name"`
}

type commonResult struct {
	gophercloud.Result
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Policy.
type CreateResult struct {
	commonResult
}
