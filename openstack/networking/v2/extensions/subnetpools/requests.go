package subnetpools

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToSubnetPoolListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the subnetpool attributes you want to see returned.
// SortKey allows you to sort by a particular subnetpool attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type ListOpts struct {
	ID               string `q:"id"`
	Name             string `q:"name"`
	DefaultQuota     int    `q:"default_quota"`
	TenantID         string `q:"tenant_id"`
	ProjectID        string `q:"project_id"`
	DefaultPrefixLen int    `q:"default_prefixlen"`
	MinPrefixLen     int    `q:"min_prefixlen"`
	MaxPrefixLen     int    `q:"max_prefixlen"`
	AddressScopeID   string `q:"address_scope_id"`
	IPVersion        int    `q:"ip_version"`
	Shared           *bool  `q:"shared"`
	Description      string `q:"description"`
	IsDefault        *bool  `q:"is_default"`
	Limit            int    `q:"limit"`
	Marker           string `q:"marker"`
	SortKey          string `q:"sort_key"`
	SortDir          string `q:"sort_dir"`
	Tags             string `q:"tags"`
	TagsAny          string `q:"tags-any"`
	NotTags          string `q:"not-tags"`
	NotTagsAny       string `q:"not-tags-any"`
	// type int does not allow to filter with revision_number=0
	RevisionNumber int `q:"revision_number"`
}

// ToSubnetPoolListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToSubnetPoolListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// subnetpools. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only the subnetpools owned by the project
// of the user submitting the request, unless the user has the administrative role.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToSubnetPoolListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return SubnetPoolPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific subnetpool based on its ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSubnetPoolCreateMap() (map[string]any, error)
}

// CreateOpts specifies parameters of a new subnetpool.
type CreateOpts struct {
	// Name is the human-readable name of the subnetpool.
	Name string `json:"name"`

	// DefaultQuota is the per-project quota on the prefix space
	// that can be allocated from the subnetpool for project subnets.
	DefaultQuota int `json:"default_quota,omitempty"`

	// TenantID is the id of the Identity project.
	TenantID string `json:"tenant_id,omitempty"`

	// ProjectID is the id of the Identity project.
	ProjectID string `json:"project_id,omitempty"`

	// Prefixes is the list of subnet prefixes to assign to the subnetpool.
	// Neutron API merges adjacent prefixes and treats them as a single prefix.
	// Each subnet prefix must be unique among all subnet prefixes in all subnetpools
	// that are associated with the address scope.
	Prefixes []string `json:"prefixes"`

	// DefaultPrefixLen is the size of the prefix to allocate when the cidr
	// or prefixlen attributes are omitted when you create the subnet.
	// Defaults to the MinPrefixLen.
	DefaultPrefixLen int `json:"default_prefixlen,omitempty"`

	// MinPrefixLen is the smallest prefix that can be allocated from a subnetpool.
	// For IPv4 subnetpools, default is 8.
	// For IPv6 subnetpools, default is 64.
	MinPrefixLen int `json:"min_prefixlen,omitempty"`

	// MaxPrefixLen is the maximum prefix size that can be allocated from the subnetpool.
	// For IPv4 subnetpools, default is 32.
	// For IPv6 subnetpools, default is 128.
	MaxPrefixLen int `json:"max_prefixlen,omitempty"`

	// AddressScopeID is the Neutron address scope to assign to the subnetpool.
	AddressScopeID string `json:"address_scope_id,omitempty"`

	// Shared indicates whether this network is shared across all projects.
	Shared bool `json:"shared,omitempty"`

	// Description is the human-readable description for the resource.
	Description string `json:"description,omitempty"`

	// IsDefault indicates if the subnetpool is default pool or not.
	IsDefault bool `json:"is_default,omitempty"`
}

// ToSubnetPoolCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToSubnetPoolCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "subnetpool")
}

// Create requests the creation of a new subnetpool on the server.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSubnetPoolCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToSubnetPoolUpdateMap() (map[string]any, error)
}

// UpdateOpts represents options used to update a network.
type UpdateOpts struct {
	// Name is the human-readable name of the subnetpool.
	Name string `json:"name,omitempty"`

	// DefaultQuota is the per-project quota on the prefix space
	// that can be allocated from the subnetpool for project subnets.
	DefaultQuota *int `json:"default_quota,omitempty"`

	// TenantID is the id of the Identity project.
	TenantID string `json:"tenant_id,omitempty"`

	// ProjectID is the id of the Identity project.
	ProjectID string `json:"project_id,omitempty"`

	// Prefixes is the list of subnet prefixes to assign to the subnetpool.
	// Neutron API merges adjacent prefixes and treats them as a single prefix.
	// Each subnet prefix must be unique among all subnet prefixes in all subnetpools
	// that are associated with the address scope.
	Prefixes []string `json:"prefixes,omitempty"`

	// DefaultPrefixLen is yhe size of the prefix to allocate when the cidr
	// or prefixlen attributes are omitted when you create the subnet.
	// Defaults to the MinPrefixLen.
	DefaultPrefixLen int `json:"default_prefixlen,omitempty"`

	// MinPrefixLen is the smallest prefix that can be allocated from a subnetpool.
	// For IPv4 subnetpools, default is 8.
	// For IPv6 subnetpools, default is 64.
	MinPrefixLen int `json:"min_prefixlen,omitempty"`

	// MaxPrefixLen is the maximum prefix size that can be allocated from the subnetpool.
	// For IPv4 subnetpools, default is 32.
	// For IPv6 subnetpools, default is 128.
	MaxPrefixLen int `json:"max_prefixlen,omitempty"`

	// AddressScopeID is the Neutron address scope to assign to the subnetpool.
	AddressScopeID *string `json:"address_scope_id,omitempty"`

	// Description is thehuman-readable description for the resource.
	Description *string `json:"description,omitempty"`

	// IsDefault indicates if the subnetpool is default pool or not.
	IsDefault *bool `json:"is_default,omitempty"`

	// RevisionNumber implements extension:standard-attr-revisions. If != "" it
	// will set revision_number=%s. If the revision number does not match, the
	// update will fail.
	RevisionNumber *int `json:"-" h:"If-Match"`
}

// ToSubnetPoolUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToSubnetPoolUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "subnetpool")
}

// Update accepts a UpdateOpts struct and updates an existing subnetpool using the
// values provided.
func Update(ctx context.Context, c *gophercloud.ServiceClient, subnetPoolID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSubnetPoolUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	h, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		r.Err = err
		return
	}
	for k := range h {
		if k == "If-Match" {
			h[k] = fmt.Sprintf("revision_number=%s", h[k])
		}
	}
	resp, err := c.Put(ctx, updateURL(c, subnetPoolID), b, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts a unique ID and deletes the subnetpool associated with it.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
