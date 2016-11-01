package recordsets

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Zone Object

type ListOptsBuilder interface {
	ToRRSetListMap() (string, error)
}

// ListOpts allows you to query the List method.
type ListOpts struct {
	Marker string `q:"marker"`
	Limit  int    `q:"limit"`
}

func (opts ListOpts) ToRRSetListMap() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func List(client *gophercloud.ServiceClient, zoneID string) pagination.Pager {
	u := listURL(client, zoneID)
	return pagination.NewPager(client, u, func(r pagination.PageResult) pagination.Page {
		return RRSetPage{pagination.LinkedPageBase{PageResult: r}}
	})
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

// Get returns additional information about a service, given its ID.
func Get(client *gophercloud.ServiceClient, zoneID string, rrsetID string) (r GetResult) {
	_, r.Err = client.Get(rrsetURL(client, zoneID, rrsetID), &r.Body, nil)
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
