package qos

import (
	"github.com/gophercloud/gophercloud"
)

type CreateOptsBuilder interface {
	ToQoSCreateMap() (map[string]interface{}, error)
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
