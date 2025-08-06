package tapmirrors

import (
	"github.com/gophercloud/gophercloud/v2"
)

// TapMirror represents a Tap Mirror of the networking service taas extension
type TapMirror struct {
	// The ID of the Tap Mirror.
	ID string `json:"id"`

	// The name of the Tap Mirror.
	Name string `json:"name"`

	// A human-readable description of the Tap Mirror.
	Description string `json:"description"`

	// The ID of the tenant.
	TenantID string `json:"tenant_id"`

	// The ID of the project.
	ProjectID string `json:"project_id"`

	// The Port ID of the Tap Mirror, this will be the source of the mirrored traffic,
	// and this traffic will be tunneled into the GRE or ERSPAN v1 tunnel.
	// The tunnel itself is not starting from this port.
	PortID string `json:"port_id"`

	// The type of the mirroring, it can be gre or erspanv1.
	MirrorType string `json:"mirror_type"`

	// The remote IP of the Tap Mirror, this will be the remote end of the GRE or ERSPAN v1 tunnel.
	RemoteIP string `json:"remote_ip"`

	// A dictionary of direction and tunnel_id. Directions are In and Out. In specifies
	// ingress traffic to the port will be mirrored, Out specifies egress traffic will be mirrored.
	// The values of the directions are the identifiers of the ERSPAN or GRE session between
	// the source and destination, these must be unique within the project and must be convertible to int.
	Directions Directions `json:"directions"`
}

type Directions struct {
	// Unique identifier of the tunnel with ingress traffic. Must be convertible to int.
	// Omit to not capture ingress traffic.
	In string `json:"IN,omitempty"`

	// Unique identifier of the tunnel with egress traffic. Must be convertible to int.
	// Omit to not capture egress traffic.
	Out string `json:"OUT,omitempty"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a Tap Mirror.
func (r commonResult) Extract() (*TapMirror, error) {
	var s struct {
		TapMirror *TapMirror `json:"tap_mirror"`
	}
	err := r.ExtractInto(&s)
	return s.TapMirror, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Tap Mirror.
type CreateResult struct {
	commonResult
}
