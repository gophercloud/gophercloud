package sharetransfers

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToTransferCreateMap() (map[string]any, error)
}

// CreateOpts contains options for a Share transfer.
type CreateOpts struct {
	// The ID of the share to transfer.
	ShareID string `json:"share_id" required:"true"`

	// The name of the share transfer.
	Name string `json:"name,omitempty"`
}

// ToCreateMap assembles a request body based on the contents of a
// TransferOpts.
func (opts CreateOpts) ToTransferCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "transfer")
}

// Create will create a share tranfer request based on the values in CreateOpts.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTransferCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, transferURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// AcceptOpts contains options for a Share transfer accept reqeust.
type AcceptOpts struct {
	// The auth key of the share transfer to accept.
	AuthKey string `json:"auth_key" required:"true"`

	// Whether to clear access rules when accept the share.
	ClearAccessRules bool `json:"clear_access_rules,omitempty"`
}

// ToAcceptMap assembles a request body based on the contents of a
// AcceptOpts.
func (opts AcceptOpts) ToAcceptMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "accept")
}

// Accept will accept a share tranfer request based on the values in AcceptOpts.
func Accept(ctx context.Context, client *gophercloud.ServiceClient, id string, opts AcceptOpts) (r AcceptResult) {
	b, err := opts.ToAcceptMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, acceptURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a share transfer.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), &gophercloud.RequestOpts{
		// DELETE requests response with a 200 code, adding it here
		OkCodes: []int{200, 202, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToTransferListQuery() (string, error)
}

// ListOpts holds options for listing Transfers. It is passed to the sharetransfers.List
// or sharetransfers.ListDetail functions.
type ListOpts struct {
	// AllTenants will retrieve transfers of all tenants/projects. Admin
	// only.
	AllTenants bool `q:"all_tenants"`

	// The user defined name of the share transfer to filter resources by.
	Name string `q:"name"`

	// The name pattern that can be used to filter share transfers.
	NamePattern string `q:"name~"`

	// The key to sort a list of transfers. A valid value is id, name,
	// resource_type, resource_id, source_project_id, destination_project_id,
	// created_at, expires_at.
	SortKey string `q:"sort_key"`

	// The direction to sort a list of resources. A valid value is asc, or
	// desc.
	SortDir string `q:"sort_dir"`

	// Requests a page size of items.
	Limit int `q:"limit"`

	// Used in conjunction with limit to return a slice of items.
	Offset int `q:"offset"`

	// The ID of the last-seen item.
	Marker string `q:"marker"`
}

// ToTransferListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToTransferListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns Transfers optionally limited by the conditions provided in ListOpts.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToTransferListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := TransferPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// List returns Transfers with details optionally limited by the conditions
// provided in ListOpts.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listDetailURL(client)
	if opts != nil {
		query, err := opts.ToTransferListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := TransferPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// Get retrieves the Transfer with the provided ID. To extract the Transfer object
// from the response, call the Extract method on the GetResult.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
