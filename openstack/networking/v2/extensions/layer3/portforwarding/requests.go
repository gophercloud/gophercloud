package portforwarding

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type ListOptsBuilder interface {
	ToPortForwardingListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the port forwarding attributes you want to see returned. SortKey allows you to
// sort by a particular network attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	ID                string `q:"id"`
	Description       string `q:"description"`
	InternalPortID    string `q:"internal_port_id"`
	ExternalPort      string `q:"external_port"`
	InternalIPAddress string `q:"internal_ip_address"`
	Protocol          string `q:"protocol"`
	InternalPort      string `q:"internal_port"`
	SortKey           string `q:"sort_key"`
	SortDir           string `q:"sort_dir"`
	Fields            string `q:"fields"`
	Limit             int    `q:"limit"`
	Marker            string `q:"marker"`
}

// ToPortForwardingListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToPortForwardingListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// Port Forwarding resources. It accepts a ListOpts struct, which allows you to
// filter and sort the returned collection for greater efficiency.
func List(c gophercloud.Client, opts ListOptsBuilder, id string) pagination.Pager {
	url := portForwardingUrl(c, id)
	if opts != nil {
		query, err := opts.ToPortForwardingListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return PortForwardingPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a particular port forwarding resource based on its unique ID.
func Get(ctx context.Context, c gophercloud.Client, floatingIpId string, pfId string) (r GetResult) {
	resp, err := c.Get(ctx, singlePortForwardingUrl(c, floatingIpId, pfId), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOpts contains all the values needed to create a new port forwarding
// resource. All attributes are required.
type CreateOpts struct {
	Description       string `json:"description,omitempty"`
	InternalPortID    string `json:"internal_port_id"`
	InternalIPAddress string `json:"internal_ip_address"`
	InternalPort      int    `json:"internal_port"`
	ExternalPort      int    `json:"external_port"`
	Protocol          string `json:"protocol"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToPortForwardingCreateMap() (map[string]any, error)
}

// ToPortForwardingCreateMap allows CreateOpts to satisfy the CreateOptsBuilder
// interface
func (opts CreateOpts) ToPortForwardingCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "port_forwarding")
}

// Create accepts a CreateOpts struct and uses the values provided to create a
// new port forwarding for an existing floating IP.
func Create(ctx context.Context, c gophercloud.Client, floatingIpId string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPortForwardingCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, portForwardingUrl(c, floatingIpId), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOpts contains the values used when updating a port forwarding resource.
type UpdateOpts struct {
	Description       *string `json:"description,omitempty"`
	InternalPortID    string  `json:"internal_port_id,omitempty"`
	InternalIPAddress string  `json:"internal_ip_address,omitempty"`
	InternalPort      int     `json:"internal_port,omitempty"`
	ExternalPort      int     `json:"external_port,omitempty"`
	Protocol          string  `json:"protocol,omitempty"`
}

// ToPortForwardingUpdateMap allows UpdateOpts to satisfy the UpdateOptsBuilder
// interface
func (opts UpdateOpts) ToPortForwardingUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "port_forwarding")
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToPortForwardingUpdateMap() (map[string]any, error)
}

// Update allows port forwarding resources to be updated.
func Update(ctx context.Context, c gophercloud.Client, fipID string, pfID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPortForwardingUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, singlePortForwardingUrl(c, fipID, pfID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will permanently delete a particular port forwarding for a given floating ID.
func Delete(ctx context.Context, c gophercloud.Client, floatingIpId string, pfId string) (r DeleteResult) {
	resp, err := c.Delete(ctx, singlePortForwardingUrl(c, floatingIpId, pfId), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
