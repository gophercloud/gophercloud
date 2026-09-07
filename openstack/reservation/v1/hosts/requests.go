package hosts

import (
	"context"
	"fmt"
	"maps"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add parameters to the List request.
type ListOptsBuilder interface {
	ToHostListQuery() (string, error)
}

// ListOpts allows the filtering of paginated collections through the API.
//
// Blazar accepts query parameters on this endpoint but ignores them. It returns every host.
type ListOpts struct{}

// ToHostListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToHostListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List retrieves a list of compute hosts.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToHostListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return HostPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder allows extensions to add parameters to the Create request.
type CreateOptsBuilder interface {
	ToHostCreateMap() (map[string]any, error)
}

// CreateOpts specifies the parameters for enrolling a host into the freepool.
type CreateOpts struct {
	// Name is the name by which Nova knows the hypervisor to enroll.
	Name string `json:"name" required:"true"`

	// ExtraCapabilities holds the operator-defined capabilities to set on the host.
	ExtraCapabilities map[string]any `json:"-"`
}

// ToHostCreateMap formats a CreateOpts into a request body.
func (opts CreateOpts) ToHostCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	for k, v := range opts.ExtraCapabilities {
		if _, ok := b[k]; ok {
			return nil, fmt.Errorf("extra capability %q collides with a host field", k)
		}
		b[k] = v
	}

	return b, nil
}

// Create enrols a compute host into the Blazar freepool.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToHostCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a specific compute host based on its unique ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add parameters to the Update request.
type UpdateOptsBuilder interface {
	ToHostUpdateMap() (map[string]any, error)
}

// UpdateOpts specifies the parameters for updating a host. Only the extra
// capabilities of a host can be changed.
type UpdateOpts struct {
	// ExtraCapabilities holds the operator-defined capabilities to set on the host.
	ExtraCapabilities map[string]any
}

// ToHostUpdateMap formats an UpdateOpts into a request body.
func (opts UpdateOpts) ToHostUpdateMap() (map[string]any, error) {
	if len(opts.ExtraCapabilities) == 0 {
		return nil, gophercloud.ErrMissingInput{Argument: "ExtraCapabilities"}
	}

	b := make(map[string]any, len(opts.ExtraCapabilities))
	maps.Copy(b, opts.ExtraCapabilities)

	return b, nil
}

// Update changes the extra capabilities of a compute host.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToHostUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete removes a compute host from the Blazar freepool.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
