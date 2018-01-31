package aggregates

import (
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List makes a request against the API to list aggregates.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, aggregatesListURL(client), func(r pagination.PageResult) pagination.Page {
		return AggregatesPage{pagination.SinglePageBase(r)}
	})
}

type CreateOpts struct {
	// The name of the host aggregate.
	Name string `json:"name" required:"true"`

	// The availability zone of the host aggregate.
	// You should use a custom availability zone rather than
	// the default returned by the os-availability-zone API.
	// The availability zone must not include ‘:’ in its name.
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

func (opts CreateOpts) ToAggregatesCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "aggregate")
}

// Create makes a request against the API to create an aggregate.
func Create(client *gophercloud.ServiceClient, opts CreateOpts) (r CreateResult) {
	b, err := opts.ToAggregatesCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(aggregatesCreateURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete makes a request against the API to delete an aggregate.
func Delete(client *gophercloud.ServiceClient, aggregateID int) (r DeleteResult) {
	v := strconv.Itoa(aggregateID)
	_, r.Err = client.Delete(aggregatesDeleteURL(client, v), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get makes a request against the API to get details for an specific aggregate.
func Get(client *gophercloud.ServiceClient, aggregateID int) (r GetResult) {
	v := strconv.Itoa(aggregateID)
	_, r.Err = client.Get(aggregatesGetURL(client, v), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
