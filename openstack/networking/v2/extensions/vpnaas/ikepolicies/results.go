package ikepolicies

// Policy is an IKE Policy
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

	// TransformProtocol is the transform protocol
	TransformProtocol string `json:"transform_protocol"`

	// Lifetime is the lifetime of the security association
	Lifetime Lifetime `json:"lifetime"`

	// ID is the ID of the policy
	ID string `json:"id"`

	// Phase1NegotiationMode is the IKE mode
	Phase1NegotiationMode string `json:"phase1_negotiation_mode"`

	// IkeVersion is the IKE version.
	IkeVersion string `json:"ike_version"`
}

type Lifetime struct {
	// Units is the unit for the lifetime
	// Default is seconds
	Units string `json:"units"`

	// Value is the lifetime
	// Default is 3600
	Value int `json:"value"`
}
