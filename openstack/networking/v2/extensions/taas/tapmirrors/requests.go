package tapmirrors

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

type MirrorType string

const (
	MirrorTypeErspanv1 MirrorType = "erspanv1"
	MirrorTypeGre      MirrorType = "gre"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToTapMirrorCreateMap() (map[string]any, error)
}

// CreateOpts contains all the values needed to create a new tap mirror
type CreateOpts struct {
	// The name of the Tap Mirror.
	Name string `json:"name"`

	// A human-readable description of the Tap Mirror.
	Description string `json:"description,omitempty"`

	// The ID of the project. The caller must have an admin role in
	// order to set this. Otherwise, this field is left unset
	// and the caller will be the owner.
	TenantID string `json:"tenant_id,omitempty"`

	// The Port ID of the Tap Mirror, this will be the source of the mirrored traffic,
	// and this traffic will be tunneled into the GRE or ERSPAN v1 tunnel.
	// The tunnel itself is not starting from this port.
	PortID string `json:"port_id"`

	// The type of the mirroring, it can be gre or erspanv1.
	MirrorType MirrorType `json:"mirror_type"`

	// The remote IP of the Tap Mirror, this will be the remote end of the GRE or ERSPAN v1 tunnel.
	RemoteIP string `json:"remote_ip"`

	// A dictionary of direction and tunnel_id. Directions are IN and OUT.
	// The values of the directions must be unique within the project and
	// must be convertible to int.
	Directions Directions `json:"directions"`
}

// ToTapMirrorCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToTapMirrorCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "tap_mirror")
}

// Create accepts a CreateOpts struct and uses the values to create a new Tap Mirror.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTapMirrorCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
