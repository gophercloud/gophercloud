package nodes

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToNodeCreateMap() (map[string]interface{}, error)
}

// CreateOpts params
type CreateOpts struct {
	Role      string                 `json:"role,omitempty"`
	ProfileID string                 `json:"profile_id" required:"true"`
	ClusterID string                 `json:"cluster_id,omitempty"`
	Name      string                 `json:"name" required:"true"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ToNodeCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToNodeCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "node")
}

// Create requests the creation of a new node.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToNodeCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	var result *http.Response
	result, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	if r.Err == nil {
		r.Header = result.Header
	}

	return
}

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToNodeListQuery() (string, error)
}

// ListOpts params
type ListOpts struct {
	Limit         int    `q:"limit"`
	Marker        string `q:"marker"`
	Sort          string `q:"sort"`
	GlobalProject *bool  `q:"global_project"`
	ClusterID     string `q:"cluster_id"`
	Name          string `q:"name"`
	Status        string `q:"status"`
}

// ToNodeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToNodeListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List instructs OpenStack to provide a list of nodes.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToNodeListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return NodePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Delete deletes the specified node ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(deleteURL(client, id), nil)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
