package speaker

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// For the sake of consistency, this function will return a pagination.Pager
func List(c *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return BGPSpeakerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieve the specific bgp speaker by its uuid
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSpeakerCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a network.
type CreateOpts struct {
	Name                          string   `json:"name"`
	IPVersion                     int      `json:"ip_version"`
	AdvertiseFloatingIPHostRoutes bool     `json:"advertise_floating_ip_host_routes,omitempty"`
	AdvertiseTenantNetworks       bool     `json:"advertise_tenant_networks,omitempty"`
	LocalAS                       string   `json:"local_as"`
	Networks                      []string `json:"networks,omitempty"`
}

// ToSpeakerCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToSpeakerCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "bgp_speaker")
}

func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSpeakerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(createURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts a unique ID and deletes the bgp speaker associated with it.
func Delete(c *gophercloud.ServiceClient, SpeakerID string) (r DeleteResult) {
	resp, err := c.Delete(deleteURL(c, SpeakerID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
