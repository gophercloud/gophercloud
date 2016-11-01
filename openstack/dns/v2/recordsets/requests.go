package recordsets

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToRecordSetListQuery() (string, error)
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

	Data        string `q:"data"`
	Description string `q:"description"`
	Name        string `q:"name"`
	SortDir     string `q:"sort_dir"`
	SortKey     string `q:"sort_key"`
	Status      string `q:"status"`
	TTL         int    `q:"ttl"`
	Type        string `q:"type"`
	ZoneID      string `q:"zone_id"`
}

// ToRecordSetListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRecordSetListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListByZone implements the recordset list request.
func ListByZone(client *gophercloud.ServiceClient, zoneID string, opts ListOptsBuilder) pagination.Pager {
	url := baseURL(client, zoneID)
	if opts != nil {
		query, err := opts.ToRecordSetListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RecordSetPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get implements the recordset get request.
func Get(client *gophercloud.ServiceClient, zoneID string, rrsetID string) (r GetResult) {
	_, r.Err = client.Get(rrsetURL(client, zoneID, rrsetID), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional attributes to the Update request.
type CreateOptsBuilder interface {
	ToZoneCreateMap() (map[string]interface{}, error)
}

// UpdateOpts specifies the base attributes that may be updated on an existing server.
type CreateOpts struct {
	ZoneID  string   `json:"zone_id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Type    string   `json:"type,omitempty"`
	Records []string `json:"records,omitempty"`
	TTL     int      `json:"ttl,omitempty"`
}

// ToServerUpdateMap formats an UpdateOpts structure into a request body.
func (opts CreateOpts) ToZoneCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Update changes the service type of an existing service.
func Create(client *gophercloud.ServiceClient, zoneID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToZoneCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(listURL(client, zoneID), &b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional attributes to the Update request.
type UpdateOptsBuilder interface {
	ToRRSetUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts specifies the base attributes that may be updated on an existing server.
type UpdateOpts struct {
	TTL     int      `json:"ttl,omitempty"`
	Records []string `json:"records,omitempty"`
}

// ToServerUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToRRSetUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Update changes the service type of an existing service.
func Update(client *gophercloud.ServiceClient, zoneID string, rrsetID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRRSetUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(rrsetURL(client, zoneID, rrsetID), &b, &r.Body, nil)
	return
}

// Delete removes an existing service.
// It either deletes all associated endpoints, or fails until all endpoints are deleted.
func Delete(client *gophercloud.ServiceClient, zoneID string, rrsetID string) (r DeleteResult) {
	_, r.Err = client.Delete(rrsetURL(client, zoneID, rrsetID), nil)
	return
}
