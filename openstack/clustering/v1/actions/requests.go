package actions

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOpts params
type ListOpts struct {
	// Limit instructs List to refrain from sending excessively large lists of clusters
	Limit int `q:"limit,omitempty"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Marker string `q:"marker,omitempty"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Sort string `q:"sort,omitempty"`

	// GlobalProject indicates whether to include resources for all projects or resources for the current project
	GlobalProject string `q:"global_project,omitempty"`

	// Name to filter the response by the specified name property of the object
	Name string `q:"name,omitempty"`

	// Status to filter the resource collection by the specified target property
	Target string `q:"target,omitempty"`

	// Status to filter the resource collection by the specified action property
	Action string `q:"action,omitempty"`

	// Status to filter the resource collection by the specified status property
	Status string `q:"status,omitempty"`
}

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToActionListQuery() (string, error)
}

// ToClusterListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToActionListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDetail instructs OpenStack to provide a list of actions.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToActionListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ActionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details of a single action. Use ExtractAction to convert its
// result into a Action.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}
