package profiles

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"net/http"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToProfileCreateMap() (map[string]interface{}, error)
}

// CreateOpts params
type CreateOpts struct {
	// Name is the name of the cluster.
	Name string `json:"name" required:"true"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`

	Spec map[string]interface{} `json:"spec" required:"true"`
}

// ToProfileCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToProfileCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "profile")
}

// Create requests the creation of a new profile.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToProfileCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	r.Header = result.Header
	return
}

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToProfileListQuery() (string, error)
}

// ListOpts params
type ListOpts struct {
	// GlobalProject indicates whether to include resources for all projects or resources for the current project
	GlobalProject bool `q:"global_project"`

	// Limit instructs List to refrain from sending excessively large lists of profiles
	Limit int `q:"limit"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Marker string `q:"marker"`

	// Name to filter the response by the specified name property of the object
	Name string `q:"name"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Sort string `q:"sort"`

	// Filter the response by the specified type property of the object
	Type string `q:"type"`
}

// ToProfileListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToProfileListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get retrieves details of a single profile. Use ExtractProfile to convert its
// result into a Node.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	r.Header = result.Header
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToProfileUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts implements UpdateOpts
type UpdateOpts struct {
	// Profile A structured description of a profile object.
	Profile map[string]interface{} `json:"-"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Name is the name of the new profile.
	Name string `json:"name,omitempty"`
}

// ToProfileUpdateMap assembles a request body based on the contents of
// UpdateOpts.
func (opts UpdateOpts) ToProfileUpdateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "profile")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Update implements profile updated request.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToProfileUpdateMap()
	if err != nil {
		r.Err = err
		return r
	}
	var result *http.Response
	result, r.Err = client.Patch(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	r.Header = result.Header
	return
}

// ListDetail instructs OpenStack to provide a list of profiles.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToProfileListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ProfilePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Delete deletes the specified node ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(deleteURL(client, id), nil)
	r.Header = result.Header
	return
}
