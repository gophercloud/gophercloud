package request

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add parameters to the List request.
type ListOptsBuilder interface {
	ToTransferRequestListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned.
// https://developer.openstack.org/api-ref/dns/
type ListOpts struct {
	Status string `q:"status"`
}

// ToTransferRequestListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToTransferRequestListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List implements a transfer request List request.
func List(client gophercloud.Client, opts ListOptsBuilder) pagination.Pager {
	url := baseURL(client)
	if opts != nil {
		query, err := opts.ToTransferRequestListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return TransferRequestPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get returns information about a transfer request, given its ID.
func Get(ctx context.Context, client gophercloud.Client, transferRequestID string) (r GetResult) {
	resp, err := client.Get(ctx, resourceURL(client, transferRequestID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional attributes to the
// Create request.
type CreateOptsBuilder interface {
	ToTransferRequestCreateMap() (map[string]any, error)
}

// CreateOpts specifies the attributes used to create a transfer request.
type CreateOpts struct {
	// TargetProjectID is ID that the request will be limited to. No other project
	// will be allowed to accept this request.
	TargetProjectID string `json:"target_project_id,omitempty"`

	// Description of the transfer request.
	Description string `json:"description,omitempty"`
}

// ToTransferRequestCreateMap formats an CreateOpts structure into a request body.
func (opts CreateOpts) ToTransferRequestCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Create implements a transfer request create request.
func Create(ctx context.Context, client gophercloud.Client, zoneID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTransferRequestCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client, zoneID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusCreated, http.StatusAccepted},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional attributes to the
// Update request.
type UpdateOptsBuilder interface {
	ToTransferRequestUpdateMap() (map[string]any, error)
}

// UpdateOpts specifies the attributes to update a transfer request.
type UpdateOpts struct {
	// TargetProjectID is ID that the request will be limited to. No other project
	// will be allowed to accept this request.
	TargetProjectID string `json:"target_project_id,omitempty"`

	// Description of the transfer request.
	Description string `json:"description,omitempty"`
}

// ToTransferRequestUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToTransferRequestUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Update implements a transfer request update request.
func Update(ctx context.Context, client gophercloud.Client, transferID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToTransferRequestUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, resourceURL(client, transferID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK, http.StatusAccepted},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete implements a transfer request delete request.
func Delete(ctx context.Context, client gophercloud.Client, transferID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, resourceURL(client, transferID), &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
