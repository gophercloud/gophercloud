package projects

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToProjectListQuery() (string, error)
}

// ListOpts allows you to query the List method.
type ListOpts struct {
	// DomainID filters the response by a domain ID.
	DomainID string `q:"domain_id"`

	// Enabled filters the response by enabled projects.
	Enabled *bool `q:"enabled"`

	// IsDomain filters the response by projects that are domains.
	// Setting this to true is effectively listing domains.
	IsDomain *bool `q:"is_domain"`

	// Name filters the response by project name.
	Name string `q:"name"`

	// ParentID filters the response by projects of a given parent project.
	ParentID string `q:"parent_id"`
}

// ToProjectListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToProjectListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerats the Projects to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToProjectListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ProjectPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetOptsBuilder allows extensions to add additional parameters to
// the Get request.
type GetOptsBuilder interface {
	ToProjectGetQuery() (string, error)
}

// GetOpts allows you to modify the details included in the Get request.
type GetOpts struct{}

// ToProjectGetQuery formats a GetOpts into a query string.
func (opts GetOpts) ToProjectGetQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get retrieves details on a single project, by ID.
func Get(client *gophercloud.ServiceClient, id string, opts GetOptsBuilder) (r GetResult) {
	url := getURL(client, id)
	if opts != nil {
		query, err := opts.ToProjectGetQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	_, r.Err = client.Get(url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
