package portgroups

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToPortGroupCreateMap() (map[string]any, error)
}

// CreateOpts specifies port group creation parameters
type CreateOpts struct {
	// NodeUUID is the UUID of the Node this resource belongs to
	NodeUUID string `json:"node_uuid" required:"true"`

	// Address is the physical hardware address of this Portgroup,
	// typically the hardware MAC address
	Address string `json:"address,omitempty"`

	// Name is a human-readable identifier for the Portgroup resource
	Name string `json:"name,omitempty"`

	// Mode is the mode of the port group. For possible values, refer to
	// https://www.kernel.org/doc/Documentation/networking/bonding.txt
	// If not specified, it will be set to the value of the
	// [DEFAULT]default_portgroup_mode configuration option.
	// When set, cannot be removed from the port group.
	Mode string `json:"mode,omitempty"`

	// StandalonePortsSupported indicates whether ports that are members
	// of this portgroup can be used as stand-alone ports
	StandalonePortsSupported bool `json:"standalone_ports_supported,omitempty"`

	// Properties contains key/value properties related to the port
	// group's configuration
	Properties map[string]interface{} `json:"properties,omitempty"`

	// Extra is a set of one or more arbitrary metadata key and value pairs
	Extra map[string]string `json:"extra,omitempty"`

	// UUID is the UUID for the resource
	UUID string `json:"uuid,omitempty"`
}

// ToPortGroupCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToPortGroupCreateMap() (map[string]any, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Create requests a node to be created
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToPortGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, createURL(client), reqBody, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToPortGroupListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing specific query parameters to the API.
type ListOpts struct {
	// Node filters the list to return only Portgroups associated with this
	// specific node (name or UUID)
	Node string `q:"node,omitempty"`

	// Address filters the list to return only Portgroups with the specified
	// physical hardware address (typically MAC)
	Address string `q:"address,omitempty"`

	// Fields specifies which fields to return in the response
	// For example: "uuid,name" will return only those fields
	Fields []string `q:"fields,omitempty"`

	// Limit requests a page size of items. Returns a number of items up to a limit value.
	// Use with marker to implement pagination. Cannot exceed max_limit set in configuration.
	Limit int `q:"limit,omitempty"`

	// Marker is the ID of the last-seen item. Use with limit to implement pagination.
	// Use the ID from the response as marker in subsequent limited requests.
	Marker string `q:"marker,omitempty"`

	// SortDir sorts the response by the requested direction.
	// Valid values are "asc" or "desc". Default is "asc".
	SortDir string `q:"sort_dir,omitempty"`

	// SortKey sorts the response by this attribute value.
	// Default is "id". Multiple sort key/direction pairs can be specified.
	SortKey string `q:"sort_key,omitempty"`

	// Detail indicates whether to show detailed information about the resource.
	// Cannot be true if Fields parameter is specified.
	Detail bool `q:"detail,omitempty"`
}

// ToPortGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToPortGroupListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list portgroups accessible to you.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToPortGroupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return PortGroupsPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get requests the details of an portgroup by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete requests the deletion of an portgroup
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
