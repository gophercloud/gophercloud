package conductors

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToConductorListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the conductor attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListOpts struct {
	// One or more fields to be returned in the response.
	Fields []string `q:"fields" format:"comma-separated"`

	// Requests a page size of items.
	Limit int `q:"limit"`

	// The ID of the last-seen item.
	Marker string `q:"marker"`

	// Sorts the response by the requested sort direction.
	SortDir string `q:"sort_dir"`

	// Sorts the response by the this attribute value.
	SortKey string `q:"sort_key"`

	// Provide additional information for the BIOS Settings
	Detail bool `q:"detail"`
}

// ToConductorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToConductorListQuery() (string, error) {
	if opts.Detail && len(opts.Fields) > 0 {
		return "", fmt.Errorf("cannot have both fields and detail options for conductors")
	}

	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list conductors accessible to you.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToConductorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ConductorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get requests details on a single conductor by hostname
func Get(ctx context.Context, client *gophercloud.ServiceClient, name string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, name), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
