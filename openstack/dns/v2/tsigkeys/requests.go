package tsigkeys

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add parameters to the List request.
type ListOptsBuilder interface {
	ToTSIGKeyListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
// https://docs.openstack.org/api-ref/dns/
type ListOpts struct {
	// Integer value for the limit of values to return.
	Limit int `q:"limit"`

	// UUID of the TSIG key at which you want to set a marker.
	Marker string `q:"marker"`

	// Name of the TSIG key.
	Name string `q:"name"`

	// Algorithm used by the TSIG key.
	Algorithm string `q:"algorithm"`

	// Scope of the TSIG key (ZONE or POOL).
	Scope string `q:"scope"`
}

// ToTSIGKeyListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToTSIGKeyListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List implements a TSIG key List request.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := baseURL(client)
	if opts != nil {
		query, err := opts.ToTSIGKeyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return TSIGKeyPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get returns information about a TSIG key, given its ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, tsigkeyID string) (r GetResult) {
	resp, err := client.Get(ctx, tsigkeyURL(client, tsigkeyID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional attributes to the
// Create request.
type CreateOptsBuilder interface {
	ToTSIGKeyCreateMap() (map[string]any, error)
}

// CreateOpts specifies the attributes used to create a TSIG key.
type CreateOpts struct {
	// Name of the TSIG key.
	Name string `json:"name" required:"true"`

	// Algorithm is the TSIG algorithm (e.g., hmac-sha256, hmac-sha512).
	Algorithm string `json:"algorithm" required:"true"`

	// Secret is the base64-encoded secret key.
	Secret string `json:"secret" required:"true"`

	// Scope defines the scope of the TSIG key (ZONE or POOL).
	Scope string `json:"scope" required:"true"`

	// ResourceID is the ID of the resource (zone or pool) this key is associated with.
	ResourceID string `json:"resource_id,omitempty"`
}

// ToTSIGKeyCreateMap formats a CreateOpts structure into a request body.
func (opts CreateOpts) ToTSIGKeyCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create implements a TSIG key create request.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTSIGKeyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, baseURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional attributes to the
// Update request.
type UpdateOptsBuilder interface {
	ToTSIGKeyUpdateMap() (map[string]any, error)
}

// UpdateOpts specifies the attributes to update a TSIG key.
type UpdateOpts struct {
	// Name of the TSIG key.
	Name string `json:"name,omitempty"`

	// Algorithm is the TSIG algorithm.
	Algorithm string `json:"algorithm,omitempty"`

	// Secret is the base64-encoded secret key.
	Secret string `json:"secret,omitempty"`

	// Scope defines the scope of the TSIG key.
	Scope string `json:"scope,omitempty"`

	// ResourceID is the ID of the resource this key is associated with.
	ResourceID string `json:"resource_id,omitempty"`
}

// ToTSIGKeyUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToTSIGKeyUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Update implements a TSIG key update request.
func Update(ctx context.Context, client *gophercloud.ServiceClient, tsigkeyID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToTSIGKeyUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, tsigkeyURL(client, tsigkeyID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete implements a TSIG key delete request.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, tsigkeyID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, tsigkeyURL(client, tsigkeyID), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
