package volumetransfers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOpts contains options for a Volume transfer.
type CreateOpts struct {
	// The ID of the volume to transfer.
	VolumeID string `json:"volume_id" required:"true"`

	// The name of the volume transfer
	Name string `json:"name,omitempty"`
}

// ToCreateMap assembles a request body based on the contents of a
// TransferOpts.
func (opts CreateOpts) ToCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "transfer")
}

// Create will create a volume tranfer request based on the values in CreateOpts.
func Create(client *gophercloud.ServiceClient, opts CreateOpts) (r CreateResult) {
	b, err := opts.ToCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(transferURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

// AcceptOpts contains options for a Volume transfer accept reqeust.
type AcceptOpts struct {
	// The auth key of the volume transfer to accept.
	AuthKey string `json:"auth_key" required:"true"`
}

// ToAcceptMap assembles a request body based on the contents of a
// AcceptOpts.
func (opts AcceptOpts) ToAcceptMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "accept")
}

// Accept will accept a volume tranfer request based on the values in AcceptOpts.
func Accept(client *gophercloud.ServiceClient, id string, opts AcceptOpts) (r CreateResult) {
	b, err := opts.ToAcceptMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(acceptURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

// Delete deletes a volume transfer.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToTransferListQuery() (string, error)
}

// ListOpts holds options for listing Transfers. It is passed to the transfers.List
// function.
type ListOpts struct {
	// AllTenants will retrieve transfers of all tenants/projects.
	AllTenants bool `q:"all_tenants"`

	// Comma-separated list of sort keys and optional sort directions in the
	// form of <key>[:<direction>].
	Sort string `q:"sort"`

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
		return TransferPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves the Transfer with the provided ID. To extract the Transfer object
// from the response, call the Extract method on the GetResult.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}
