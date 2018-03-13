package policies

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"net/http"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToPolicyCreateMap() (map[string]interface{}, error)
}

// CreateOpts params
type CreateOpts struct {
	// Name is the name of the cluster.
	Name string                 `json:"name" required:"true"`
	Spec map[string]interface{} `json:"spec" required:"true"`
}

// ToPolicyCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToPolicyCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "policy")
}

// Create requests the creation of a new policy.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToPolicyListQuery() (string, error)
}

// ListOpts params
type ListOpts struct {
	// Limit instructs List to refrain from sending excessively large lists of profiles
	Limit int `q:"limit"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Marker string `q:"marker"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Sort string `q:"sort"`

	// GlobalProject indicates whether to include resources for all projects or resources for the current project
	GlobalProject bool `q:"global_project"`

	// Name to filter the response by the specified name property of the object
	Name string `q:"name"`

	// Filter the response by the specified type property of the object
	Type string `q:"type"`
}

// ToPolicyListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToPolicyListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get retrieves details of a single policy. Use ExtractPolicy to convert its
// result into a Node.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	r.Header = result.Header
	return
}

// ListDetail instructs OpenStack to provide a list of policies.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToPolicyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return PolicyPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToPolicyUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts implements UpdateOpts
type UpdateOpts struct {
	// Policy A structured description of a policy object.
	Policy map[string]interface{} `json:"-"`

	// Name is the name of the new policy.
	Name string `json:"name,omitempty"`
}

// ToPolicyUpdateMap assembles a request body based on the contents of
// UpdateOpts.
func (opts UpdateOpts) ToPolicyUpdateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "policy")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Update implements profile updated request.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPolicyUpdateMap()
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

// Delete deletes the specified node ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(deleteURL(client, id), nil)
	r.Header = result.Header
	return
}

// ValidatePolicyOpts params
type ValidatePolicyOpts struct {
	Spec map[string]interface{} `json:"spec" required:"true"`
}

// ToPolicyCreateMap constructs a request body from CreateOpts.
func (opts ValidatePolicyOpts) ToValidatePolicyMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "policy")
}

// Validate policy.
func Validate(client *gophercloud.ServiceClient, opts ValidatePolicyOpts) (r CreateResult) {
	b, err := opts.ToValidatePolicyMap()
	if err != nil {
		r.Err = err
		return
	}

	var result *http.Response
	result, r.Err = client.Post(validateURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	r.Header = result.Header
	return
}
