package profiles

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToProfileCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used for creating a profile.
type CreateOpts struct {
	Name     string                 `json:"name" required:"true"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Spec     Spec                   `json:"spec" required:"true"`
}

// ToProfileCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToProfileCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "profile")
}

// Create requests the creation of a new profile on the server.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToProfileCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.PostWithContext(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves detail of a single profile.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.GetWithContext(ctx, getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToProfileListQuery() (string, error)
}

// ListOpts represents options used to list profiles.
type ListOpts struct {
	GlobalProject *bool  `q:"global_project"`
	Limit         int    `q:"limit"`
	Marker        string `q:"marker"`
	Name          string `q:"name"`
	Sort          string `q:"sort"`
	Type          string `q:"type"`
}

// ToProfileListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToProfileListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List instructs OpenStack to provide a list of profiles.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
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

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToProfileUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a profile.
type UpdateOpts struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Name     string                 `json:"name,omitempty"`
}

// ToProfileUpdateMap constructs a request body from UpdateOpts.
func (opts UpdateOpts) ToProfileUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "profile")
}

// Update updates a profile.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToProfileUpdateMap()
	if err != nil {
		r.Err = err
		return r
	}
	resp, err := client.PatchWithContext(ctx, updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes the specified profile via profile id.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.DeleteWithContext(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ValidateOptsBuilder allows extensions to add additional parameters to the
// Validate request.
type ValidateOptsBuilder interface {
	ToProfileValidateMap() (map[string]interface{}, error)
}

// ValidateOpts params
type ValidateOpts struct {
	Spec Spec `json:"spec" required:"true"`
}

// ToProfileValidateMap formats a CreateOpts into a body map.
func (opts ValidateOpts) ToProfileValidateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "profile")
}

// Validate profile.
func Validate(ctx context.Context, client *gophercloud.ServiceClient, opts ValidateOpts) (r ValidateResult) {
	b, err := opts.ToProfileValidateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.PostWithContext(ctx, validateURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
