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

// Get returns additional information about a service, given its ID.
func Get(client *gophercloud.ServiceClient, zoneID string, rrsetID string) (r GetResult) {
	_, r.Err = client.Get(rrsetURL(client, zoneID, rrsetID), &r.Body, nil)
	return
}
