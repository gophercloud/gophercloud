package aggregates

import (
	"context"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
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

func (opts CreateOpts) ToAggregatesCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "aggregate")
}

// Create makes a request against the API to create an aggregate.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOpts) (r CreateResult) {
	b, err := opts.ToAggregatesCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, aggregatesCreateURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete makes a request against the API to delete an aggregate.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, aggregateID int) (r DeleteResult) {
	v := strconv.Itoa(aggregateID)
	resp, err := client.Delete(ctx, aggregatesDeleteURL(client, v), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get makes a request against the API to get details for a specific aggregate.
func Get(ctx context.Context, client *gophercloud.ServiceClient, aggregateID int) (r GetResult) {
	v := strconv.Itoa(aggregateID)
	resp, err := client.Get(ctx, aggregatesGetURL(client, v), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateOpts struct {
	// The name of the host aggregate.
	Name string `json:"name,omitempty"`

	// The availability zone of the host aggregate.
	// You should use a custom availability zone rather than
	// the default returned by the os-availability-zone API.
	// The availability zone must not include ‘:’ in its name.
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

func (opts UpdateOpts) ToAggregatesUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "aggregate")
}

// Update makes a request against the API to update a specific aggregate.
func Update(ctx context.Context, client *gophercloud.ServiceClient, aggregateID int, opts UpdateOpts) (r UpdateResult) {
	v := strconv.Itoa(aggregateID)

	b, err := opts.ToAggregatesUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, aggregatesUpdateURL(client, v), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type AddHostOpts struct {
	// The name of the host.
	Host string `json:"host" required:"true"`
}

func (opts AddHostOpts) ToAggregatesAddHostMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "add_host")
}

// AddHost makes a request against the API to add host to a specific aggregate.
func AddHost(ctx context.Context, client *gophercloud.ServiceClient, aggregateID int, opts AddHostOpts) (r ActionResult) {
	v := strconv.Itoa(aggregateID)

	b, err := opts.ToAggregatesAddHostMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, aggregatesAddHostURL(client, v), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type RemoveHostOpts struct {
	// The name of the host.
	Host string `json:"host" required:"true"`
}

func (opts RemoveHostOpts) ToAggregatesRemoveHostMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "remove_host")
}

// RemoveHost makes a request against the API to remove host from a specific aggregate.
func RemoveHost(ctx context.Context, client *gophercloud.ServiceClient, aggregateID int, opts RemoveHostOpts) (r ActionResult) {
	v := strconv.Itoa(aggregateID)

	b, err := opts.ToAggregatesRemoveHostMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, aggregatesRemoveHostURL(client, v), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type SetMetadataOpts struct {
	Metadata map[string]any `json:"metadata" required:"true"`
}

func (opts SetMetadataOpts) ToSetMetadataMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "set_metadata")
}

// SetMetadata makes a request against the API to set metadata to a specific aggregate.
func SetMetadata(ctx context.Context, client *gophercloud.ServiceClient, aggregateID int, opts SetMetadataOpts) (r ActionResult) {
	v := strconv.Itoa(aggregateID)

	b, err := opts.ToSetMetadataMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, aggregatesSetMetadataURL(client, v), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
