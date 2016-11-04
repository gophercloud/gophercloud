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

// Get returns additional information about a service, given its ID.
func Get(client *gophercloud.ServiceClient, zoneID string) (r GetResult) {
	_, r.Err = client.Get(zoneURL(client, zoneID), &r.Body, nil)
	return
}
