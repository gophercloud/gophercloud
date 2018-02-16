package ipsecpolicies

import "github.com/gophercloud/gophercloud"

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToPolicyCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new IPSec policy
type CreateOpts struct {
	// TenantID specifies a tenant to own the IPSec policy. The caller must have
	// an admin role in order to set this. Otherwise, this field is left unset
	// and the caller will be the owner.
	TenantID string `json:"tenant_id,omitempty"`
	// The human readable description of the policy
	Description string `json:"description,omitempty"`
	// The human readable name of the policy
	// Does not have to be unique
	Name string `json:"name,omitempty"`
	// The authentication hash algorithm
	// Valid values are sha1, sha256, sha384, sha512
	// The default is sha1.
	AuthAlgorithm string `json:"auth_algorithm,omitempty"`
	// The encapsulation mode
	// A valid value is tunnel or transport
	// Default is tunnel.
	EncapsulationMode string `json:"encapsulation_mode,omitempty"`
	// The encryption algorithm
	// A valid value is 3des, aes-128, aes-192, aes-256, and so on
	// Default is aes-128
	EncryptionAlgorithm string `json:"encryption_algorithm,omitempty"`
	// Perfect forward secrecy (PFS)
	// A valid value is Group2, Group5, Group14, and so on
	// Default is Group5
	PFS string `json:"pfs,omitempty"`
	// The transform protocol
	// A valid value is ESP, AH, or AH- ESP
	// Default is ESP.
	TransformProtocol string `json:"transform_protocol,omitempty"`
	// The units for the lifetime of the security association
	// The lifetime consists of a unit and integer value
	// You can omit either the unit or value portion of the lifetime
	// Default unit is seconds and default value is 3600.
	LifetimeUnit string `json:"lifetime.units,omitempty"`
	LifetimeName string `json:"lifetime.name,omitempty"`
}

// ToPolicyCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToPolicyCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "ipsecpolicy")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// IPSec policy
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}
