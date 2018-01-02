package subnetpools

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
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
	ID               string   `q:"id"`
	Name             string   `q:"name"`
	DefaultQuota     int      `q:"default_quota"`
	TenantID         string   `q:"tenant_id"`
	ProjectID        string   `q:"project_id"`
	CreatedAt        string   `q:"created_at"`
	UpdatedAt        string   `q:"updated_at"`
	Prefixes         []string `q:"prefixes"`
	DefaultPrefixLen int      `q:"default_prefixlen"`
	MinPrefixLen     int      `q:"min_prefixlen"`
	MaxPrefixLen     int      `q:"max_prefixlen"`
	AddressScopeID   string   `q:"address_scope_id"`
	IPversion        int      `q:"ip_version"`
	Shared           bool     `q:"shared"`
	Description      string   `q:"description"`
	IsDefault        bool     `q:"is_default"`
	RevisionNumber   int      `q:"revision_number"`
	Limit            int      `q:"limit"`
	Marker           string   `q:"marker"`
	SortKey          string   `q:"sort_key"`
	SortDir          string   `q:"sort_dir"`
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
