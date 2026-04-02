package resourceclasses

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List retrieves a list of resource classes.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(client)

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ResourceClassesPage{pagination.SinglePageBase(r)}
	})
}

// Get retrieves the resource class with the provided name.
func Get(ctx context.Context, client *gophercloud.ServiceClient, name string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, name), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToResourceClassCreateMap() (map[string]any, error)
}

// CreateOpts represents the attributes of a new resource class.
type CreateOpts struct {
	Name string `json:"name" required:"true"`
}

// ToResourceClassCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToResourceClassCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create creates a new resource class.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToResourceClassCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Update ensures the existence of a custom resource class with the
// provided name (can be safely called multiple times).
func Update(ctx context.Context, client *gophercloud.ServiceClient, name string) (r UpdateResult) {
	resp, err := client.Put(ctx, updateURL(client, name), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{201, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes the resource class with the provided name.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, name string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, name), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
