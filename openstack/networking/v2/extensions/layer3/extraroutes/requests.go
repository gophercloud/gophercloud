package extraroutes

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

// OptsBuilder allows extensions to add additional parameters to the Add or
// Remove requests.
type OptsBuilder interface {
	ToExtraRoutesUpdateMap() (map[string]interface{}, error)
}

// Opts contains the values needed to add or remove a list og routes on a
// router.
type Opts struct {
	Routes *[]routers.Route `json:"routes,omitempty"`
}

// ToExtraRoutesUpdateMap builds a body based on Opts.
func (opts Opts) ToExtraRoutesUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "router")
}

// Add allows routers to be updated with a list of routes to be added.
func Add(c *gophercloud.ServiceClient, id string, opts OptsBuilder) (r AddResult) {
	b, err := opts.ToExtraRoutesUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(addExtraRoutesURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Remove allows routers to be updated with a list of routes to be removed.
func Remove(c *gophercloud.ServiceClient, id string, opts OptsBuilder) (r RemoveResult) {
	b, err := opts.ToExtraRoutesUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(removeExtraRoutesURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
