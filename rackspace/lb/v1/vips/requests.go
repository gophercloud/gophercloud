package vips

import (
	"github.com/racker/perigee"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List is the operation responsible for returning a paginated collection of
// load balancer virtual IP addresses.
func List(client *gophercloud.ServiceClient, loadBalancerID int) pagination.Pager {
	url := rootURL(client, loadBalancerID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return VIPPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToVIPCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	ID string

	Type string

	Version string
}

// ToVIPCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToVIPCreateMap() (map[string]interface{}, error) {
	lb := make(map[string]interface{})

	if opts.ID != "" {
		lb["id"] = opts.ID
	}
	if opts.Type != "" {
		lb["type"] = opts.Type
	}
	if opts.Version != "" {
		lb["ipVersion"] = opts.Version
	}

	return lb, nil
}

func Create(c *gophercloud.ServiceClient, lbID int, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToVIPCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("POST", rootURL(c, lbID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{202},
	})

	return res
}
