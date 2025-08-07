package tapmirrors

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
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

// TapMirrorPage is the page returned by a pager when traversing over a
// collection of Policies.
type TapMirrorPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of Endpoint groups has
// reached the end of a page and the pager seeks to traverse over a new one.
// In order to do this, it needs to construct the next page's URL.
func (r TapMirrorPage) NextPageURL(endpointURL string) (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"tap_mirrors_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether an TapMirrorPage struct is empty.
func (r TapMirrorPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractTapMirrors(r)
	return len(is) == 0, err
}

// ExtractTapMirrors accepts a Page struct, specifically an TapMirrorPage struct,
// and extracts the elements into a slice of Tap Mirror structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractTapMirrors(r pagination.Page) ([]TapMirror, error) {
	var s struct {
		TapMirrors []TapMirror `json:"tap_mirrors"`
	}
	err := (r.(TapMirrorPage)).ExtractInto(&s)
	return s.TapMirrors, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Tap Mirror.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a TapMirror.
type GetResult struct {
	commonResult
}

// DeleteResult represents the results of a Delete operation. Call its ExtractErr method
// to determine whether the operation succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UpdateResult represents the result of an update operation. Call its Extract method
// to interpret it as a TapMirror.
type UpdateResult struct {
	commonResult
}
