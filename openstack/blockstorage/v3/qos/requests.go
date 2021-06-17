package qos

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type CreateOptsBuilder interface {
	ToQoSCreateMap() (map[string]interface{}, error)
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToQoSListQuery() (string, error)
}

type QoSConsumer string

const (
	ConsumerFront QoSConsumer = "front-end"
	ConsumberBack QoSConsumer = "back-end"
	ConsumerBoth  QoSConsumer = "both"
)

// CreateOpts contains options for creating a QoS specification.
// This object is passed to the qos.Create function.
type CreateOpts struct {
	// The name of the QoS spec
	Name string `json:"name"`
	// The consumer of the QoS spec. Possible values are
	// both, front-end, back-end.
	Consumer QoSConsumer `json:"consumer,omitempty"`
	// Specs is a collection of miscellaneous key/values used to set
	// specifications for the QoS
	Specs map[string]string `json:"-"`
}

// ToQoSCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToQoSCreateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "qos_specs")
	if err != nil {
		return nil, err
	}

	if opts.Specs != nil {
		if v, ok := b["qos_specs"].(map[string]interface{}); ok {
			for key, value := range opts.Specs {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Create will create a new QoS based on the values in CreateOpts. To extract
// the QoS object from the response, call the Extract method on the
// CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToQoSCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteOptsBuilder allows extensions to add additional parameters to the
// Delete request.
type DeleteOptsBuilder interface {
	ToQoSDeleteQuery() (string, error)
}

// DeleteOpts contains options for deleting a QoS. This object is passed to
// the qos.Delete function.
type DeleteOpts struct {
	// Delete a QoS specification even if it is in-use
	Force bool `q:"force"`
}

// ToQoSDeleteQuery formats a DeleteOpts into a query string.
func (opts DeleteOpts) ToQoSDeleteQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Delete will delete the existing QoS with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string, opts DeleteOptsBuilder) (r DeleteResult) {
	url := deleteURL(client, id)
	if opts != nil {
		query, err := opts.ToQoSDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Delete(url, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type ListOpts struct {
	// Sort is Comma-separated list of sort keys and optional sort
	// directions in the form of < key > [: < direction > ]. A valid
	//direction is asc (ascending) or desc (descending).
	Sort string `q:"sort"`

	// Marker and Limit control paging.
	// Marker instructs List where to start listing from.
	Marker string `q:"marker"`

	// Limit instructs List to refrain from sending excessively large lists of
	// QoS.
	Limit int `q:"limit"`
}

// ToQoSListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToQoSListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List instructs OpenStack to provide a list of QoS.
// You may provide criteria by which List curtails its results for easier
// processing.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToQoSListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return QoSPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
