package tapmirrors

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
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

// Get retrieves a particular Tap Mirror on its ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToTapMirrorListQuery() (string, error)
}

// ListOpts allows the filtering of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the Endpoint group attributes you want to see returned.
type ListOpts struct {
	ProjectID   string     `q:"project_id"`
	Name        string     `q:"name"`
	Description string     `q:"description"`
	TenantID    string     `q:"tenant_id"`
	PortID      string     `q:"port_id"`
	MirrorType  MirrorType `q:"mirror_type"`
	RemoteIP    string     `q:"remote_ip"`
}

// ToTapMirrorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToTapMirrorListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// Tap Mirrors. It accepts a ListOpts struct, which allows you to filter
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToTapMirrorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return TapMirrorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Delete will permanently delete a Tap Mirror based on its ID.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToTapMirrorUpdateMap() (map[string]any, error)
}

// UpdateOpts contains the values used when updating a Tap Mirror.
type UpdateOpts struct {
	Description *string `json:"description,omitempty"`
	Name        *string `json:"name,omitempty"`
}

// ToTapMirrorUpdateMap casts an UpdateOpts struct to a map.
func (opts UpdateOpts) ToTapMirrorUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "tap_mirror")
}

// Update allows Tap Mirrors to be updated.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToTapMirrorUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
