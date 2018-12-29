package addressscopes

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToAddressScopeListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the subnetpool attributes you want to see returned.
// SortKey allows you to sort by a particular subnetpool attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type ListOpts struct {
	ID          string `q:"id"`
	Name        string `q:"name"`
	TenantID    string `q:"tenant_id"`
	ProjectID   string `q:"project_id"`
	IPVersion   int    `q:"ip_version"`
	Shared      *bool  `q:"shared"`
	Description string `q:"description"`
	Limit       int    `q:"limit"`
	Marker      string `q:"marker"`
	SortKey     string `q:"sort_key"`
	SortDir     string `q:"sort_dir"`
}

// ToAddressScopeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToAddressScopeListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// address-scopes. It accepts a ListOpts struct, which allows you to filter and
// sort the returned collection for greater efficiency.
//
// Default policy settings return only the address-scopes owned by the project
// of the user submitting the request, unless the user has the administrative
// role.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToAddressScopeListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return AddressScopePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific address-scope based on its ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}
