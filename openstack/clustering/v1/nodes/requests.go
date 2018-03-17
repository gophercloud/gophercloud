package nodes

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"net/http"
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
	r.Header = result.Header
	return
}

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToNodeListQuery() (string, error)
}

// ListOpts params
type ListOpts struct {
	// Limit instructs List to refrain from sending excessively large lists of nodes
	Limit int `q:"limit"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Marker string `q:"marker"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Sort string `q:"sort"`

	// GlobalProject indicates whether to include resources for all projects or resources for the current project
	GlobalProject string `q:"global_project"`

	// ClusterID the name of the cluster object
	ClusterID string `q:"cluster_id"`

	// Name to filter the response by the specified name property of the object
	Name string `q:"name"`

	// Status to filter the resource collection by the specified status property
	Status string `q:"status"`
}

// ToNodeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToNodeListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get retrieves details of a single node. Use ExtractNode to convert its
// result into a Node.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	r.Header = result.Header
	return
}

// ListDetail instructs OpenStack to provide a list of nodes.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
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

// UpdateOpts params
type UpdateOpts struct {
	Node map[string]interface{} `json:"-"`

	// Name is the name of the node.
	Name string `json:"name,omitempty"`

	ProfileID string `json:"profile_id,omitempty"`

	Role string `json:"role,omitempty"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateOptsBuilder params
type UpdateOptsBuilder interface {
	ToNodeUpdateMap() (map[string]interface{}, error)
}

// ToClusterUpdateMap constructs a request body from CreateOpts.
func (opts UpdateOpts) ToNodeUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "node")
}

// Update requests the update of a node.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToNodeUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	var result *http.Response
	result, r.Err = client.Patch(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// Delete deletes the specified node ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(deleteURL(client, id), nil)
	r.Header = result.Header
	return
}
