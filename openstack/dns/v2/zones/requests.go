package zones

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Zone Object

type ListOptsBuilder interface {
	ToZoneListMap() (string, error)
}

// ListOpts allows you to query the List method.
type ListOpts struct {
	Marker string `q:"marker"`
	Limit  int    `q:"limit"`
}

func (opts ListOpts) ToZoneListMap() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func List(client *gophercloud.ServiceClient) pagination.Pager {
	u := listURL(client)
	return pagination.NewPager(client, u, func(r pagination.PageResult) pagination.Page {
		return ZonePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOptsBuilder allows extensions to add additional attributes to the Update request.
type CreateOptsBuilder interface {
	ToZoneCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies the base attributes that may be updated on an existing server.
type CreateOpts struct {
	Email      string            `json:"email,omitempty"`
	TTL        int               `json:"ttl,omitempty"`
	Name       string            `json:"name"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Masters    []string          `json:"masters,omitempty"`
	Type       string            `json:"type,omitempty"`
}

// ToServerUpdateMap formats an UpdateOpts structure into a request body.
func (opts CreateOpts) ToZoneCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Update changes the service type of an existing service.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToZoneCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(listURL(client), &b, &r.Body, nil)
	return
}

// Get returns additional information about a service, given its ID.
func Get(client *gophercloud.ServiceClient, zoneID string) (r GetResult) {
	_, r.Err = client.Get(zoneURL(client, zoneID), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional attributes to the Update request.
type UpdateOptsBuilder interface {
	ToZoneUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts specifies the base attributes that may be updated on an existing server.
type UpdateOpts struct {
	Email   string   `json:"email,omitempty"`
	TTL     int      `json:"ttl,omitempty"`
	Masters []string `json:"masters,omitempty"`
}

// ToServerUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToZoneUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Update changes the service type of an existing service.
func Update(client *gophercloud.ServiceClient, zoneID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToZoneUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(zoneURL(client, zoneID), &b, &r.Body, nil)
	return
}

// Delete removes an existing service.
// It either deletes all associated endpoints, or fails until all endpoints are deleted.
func Delete(client *gophercloud.ServiceClient, zoneID string) (r DeleteResult) {
	_, r.Err = client.Delete(zoneURL(client, zoneID), nil)
	return
}
