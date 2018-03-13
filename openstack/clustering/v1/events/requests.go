package events

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToEventsListQuery() (string, error)
}

// ListOpts params
type ListOpts struct {
	// Limit instructs List to refrain from sending excessively large lists of profiles
	Limit         int    `q:"limit"`
	Level         string `q:"level"`
	MarkerUUID    string `q:"marker"`
	Sort          string `q:"sort"`
	GlobalProject bool   `q:"global_project"`
	OID           string `q:"oid"`
	OType         string `q:"otype"`
	OName         string `q:"otype"`
	ClusterID     string `q:"cluster_id"`
	Action        string `q:"action"`
}

// ToEventsListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToEventsListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get retrieves details of a single event. Use ExtractEvent to convert its
// result into a Event.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// ListDetail instructs OpenStack to provide a list of events
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToEventsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return EventPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
