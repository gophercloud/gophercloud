package recordsets

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToRecordSetListQuery() (string, error)
}

// ListOptsHeadersBuilder allows extensions to add additional headers to the List request.
type ListOptsHeadersBuilder interface {
	ToRecordSetListHeaders() (map[string]string, error)
}

// RequestOptsHeadersBuilder allows extensions to add additional headers to
// Create and Update requests.
type RequestOptsHeadersBuilder interface {
	ToRecordSetRequestHeaders() (map[string]string, error)
}

// DeleteOptsBuilder allows extensions to add additional headers to the
// Delete request.
type DeleteOptsBuilder interface {
	ToRecordSetDeleteHeaders() (map[string]string, error)
}

// DeleteOpts specifies headers that may be set on a Delete request.
type DeleteOpts struct {
	// AllProjects header.
	AllProjects bool `h:"X-Auth-All-Projects"`

	// SudoTenantID impersonates the given project.
	SudoTenantID string `h:"X-Auth-Sudo-Tenant-ID"`
}

// ToRecordSetDeleteHeaders formats a DeleteOpts into header parameters.
func (opts DeleteOpts) ToRecordSetDeleteHeaders() (map[string]string, error) {
	return gophercloud.BuildHeaders(opts)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
// https://developer.openstack.org/api-ref/dns/
type ListOpts struct {
	// Integer value for the limit of values to return.
	Limit int `q:"limit"`

	// UUID of the recordset at which you want to set a marker.
	Marker string `q:"marker"`

	Data         string `q:"data"`
	Description  string `q:"description"`
	Name         string `q:"name"`
	SortDir      string `q:"sort_dir"`
	SortKey      string `q:"sort_key"`
	Status       string `q:"status"`
	TTL          int    `q:"ttl"`
	Type         string `q:"type"`
	ZoneID       string `q:"zone_id"`
	AllProjects  bool   `h:"X-Auth-All-Projects"`
	SudoTenantID string `h:"X-Auth-Sudo-Tenant-ID"`
}

// ToRecordSetListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRecordSetListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ToRecordSetListHeaders formats a ListOpts into header parameters.
func (opts ListOpts) ToRecordSetListHeaders() (map[string]string, error) {
	return gophercloud.BuildHeaders(opts)
}

// ListByZone implements the recordset list request for a specific zone.
func ListByZone(client *gophercloud.ServiceClient, zoneID string, opts ListOptsBuilder) pagination.Pager {
	url := baseURL(client, zoneID)
	var h map[string]string

	if opts != nil {
		query, err := opts.ToRecordSetListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query

		// Check if opts implements the optional headers interface
		if optsWithHeaders, ok := opts.(ListOptsHeadersBuilder); ok {
			h, err = optsWithHeaders.ToRecordSetListHeaders()
			if err != nil {
				return pagination.Pager{Err: err}
			}
		}
	}

	pager := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RecordSetPage{pagination.LinkedPageBase{PageResult: r}}
	})
	pager.Headers = h
	return pager
}

// ListAll implements the recordset list request across all zones.
func ListAll(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listAllRecordSetsURL(client)
	var h map[string]string

	if opts != nil {
		query, err := opts.ToRecordSetListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query

		// Check if opts implements the optional headers interface
		if optsWithHeaders, ok := opts.(ListOptsHeadersBuilder); ok {
			h, err = optsWithHeaders.ToRecordSetListHeaders()
			if err != nil {
				return pagination.Pager{Err: err}
			}
		}
	}

	pager := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RecordSetPage{pagination.LinkedPageBase{PageResult: r}}
	})
	pager.Headers = h
	return pager
}

// Get implements the recordset Get request.
func Get(ctx context.Context, client *gophercloud.ServiceClient, zoneID string, rrsetID string) (r GetResult) {
	resp, err := client.Get(ctx, rrsetURL(client, zoneID, rrsetID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional attributes to the
// Create request.
type CreateOptsBuilder interface {
	ToRecordSetCreateMap() (map[string]any, error)
}

// CreateOpts specifies the base attributes that may be used to create a
// RecordSet.
type CreateOpts struct {
	// Name is the name of the RecordSet.
	Name string `json:"name" required:"true"`

	// Description is a description of the RecordSet.
	Description string `json:"description,omitempty"`

	// Records are the DNS records of the RecordSet.
	Records []string `json:"records,omitempty"`

	// TTL is the time to live of the RecordSet.
	TTL int `json:"ttl,omitempty"`

	// Type is the RRTYPE of the RecordSet.
	Type string `json:"type,omitempty"`

	// AllProjects header.
	AllProjects bool `h:"X-Auth-All-Projects" json:"-"`

	// SudoTenantID impersonates the given project.
	SudoTenantID string `h:"X-Auth-Sudo-Tenant-ID" json:"-"`
}

// ToRecordSetCreateMap formats an CreateOpts structure into a request body.
func (opts CreateOpts) ToRecordSetCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// ToRecordSetRequestHeaders formats a CreateOpts into header parameters.
func (opts CreateOpts) ToRecordSetRequestHeaders() (map[string]string, error) {
	return gophercloud.BuildHeaders(opts)
}

// Create creates a recordset in a given zone.
func Create(ctx context.Context, client *gophercloud.ServiceClient, zoneID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRecordSetCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := &gophercloud.RequestOpts{OkCodes: []int{201, 202}}
	// Check if opts implements the optional headers interface
	if optsWithHeaders, ok := opts.(RequestOptsHeadersBuilder); ok {
		reqOpts.MoreHeaders, err = optsWithHeaders.ToRecordSetRequestHeaders()
		if err != nil {
			r.Err = err
			return
		}
	}
	resp, err := client.Post(ctx, baseURL(client, zoneID), &b, &r.Body, reqOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional attributes to the
// Update request.
type UpdateOptsBuilder interface {
	ToRecordSetUpdateMap() (map[string]any, error)
}

// UpdateOpts specifies the base attributes that may be updated on an existing
// RecordSet.
type UpdateOpts struct {
	// Description is a description of the RecordSet.
	Description *string `json:"description,omitempty"`

	// TTL is the time to live of the RecordSet.
	TTL *int `json:"ttl,omitempty"`

	// Records are the DNS records of the RecordSet.
	Records []string `json:"records,omitempty"`

	// AllProjects header.
	AllProjects bool `h:"X-Auth-All-Projects" json:"-"`

	// SudoTenantID impersonates the given project.
	SudoTenantID string `h:"X-Auth-Sudo-Tenant-ID" json:"-"`
}

// ToRecordSetRequestHeaders formats an UpdateOpts into header parameters.
func (opts UpdateOpts) ToRecordSetRequestHeaders() (map[string]string, error) {
	return gophercloud.BuildHeaders(opts)
}

// ToRecordSetUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToRecordSetUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	// If opts.TTL was actually set, use 0 as a special value to send "null",
	// even though the result from the API is 0.
	//
	// Otherwise, don't send the TTL field.
	if opts.TTL != nil {
		ttl := *(opts.TTL)
		if ttl > 0 {
			b["ttl"] = ttl
		} else {
			b["ttl"] = nil
		}
	}

	return b, nil
}

// Update updates a recordset in a given zone.
func Update(ctx context.Context, client *gophercloud.ServiceClient, zoneID string, rrsetID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRecordSetUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := &gophercloud.RequestOpts{OkCodes: []int{200, 202}}
	// Check if opts implements the optional headers interface
	if optsWithHeaders, ok := opts.(RequestOptsHeadersBuilder); ok {
		reqOpts.MoreHeaders, err = optsWithHeaders.ToRecordSetRequestHeaders()
		if err != nil {
			r.Err = err
			return
		}
	}
	resp, err := client.Put(ctx, rrsetURL(client, zoneID, rrsetID), &b, &r.Body, reqOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete removes an existing RecordSet.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, zoneID string, rrsetID string, opts DeleteOptsBuilder) (r DeleteResult) {
	reqOpts := &gophercloud.RequestOpts{OkCodes: []int{202}}
	if opts != nil {
		h, err := opts.ToRecordSetDeleteHeaders()
		if err != nil {
			r.Err = err
			return
		}
		reqOpts.MoreHeaders = h
	}
	resp, err := client.Delete(ctx, rrsetURL(client, zoneID, rrsetID), reqOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
