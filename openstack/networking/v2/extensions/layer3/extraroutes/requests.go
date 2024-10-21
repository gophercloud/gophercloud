package extraroutes

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
)

// OptsBuilder allows extensions to add additional parameters to the Add or
// Remove requests.
type OptsBuilder interface {
	ToExtraRoutesUpdateMap() (map[string]any, error)
}

// Opts contains the values needed to add or remove a list og routes on a
// router.
type Opts struct {
	Routes *[]routers.Route `json:"routes,omitempty"`
}

// ToExtraRoutesUpdateMap builds a body based on Opts.
func (opts Opts) ToExtraRoutesUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "router")
}

// Add allows routers to be updated with a list of routes to be added.
func Add(ctx context.Context, c *gophercloud.ServiceClient, id string, opts OptsBuilder) (r AddResult) {
	b, err := opts.ToExtraRoutesUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, addExtraRoutesURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Remove allows routers to be updated with a list of routes to be removed.
func Remove(ctx context.Context, c *gophercloud.ServiceClient, id string, opts OptsBuilder) (r RemoveResult) {
	b, err := opts.ToExtraRoutesUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, removeExtraRoutesURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
