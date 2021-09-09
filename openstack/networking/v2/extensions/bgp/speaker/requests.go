package speaker

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"log"
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
	AdvertiseFloatingIPHostRoutes bool     `json:"advertise_floating_ip_host_routes"`
	AdvertiseTenantNetworks       bool     `json:"advertise_tenant_networks"`
	LocalAS                       string   `json:"local_as"`
	Networks                      []string `json:"networks,omitempty"`
}

func BuildCreateOpts(name string, ipver int, advrts bool, advnets bool, localas string, nets []string) CreateOpts {
	var r CreateOpts
	r.Name = name
	if ipver == 4 || ipver == 6 {
		r.IPVersion = ipver
	} else {
		log.Panicf("Invalid IP version %d", ipver)
	}
	r.AdvertiseFloatingIPHostRoutes = advrts
	r.AdvertiseTenantNetworks = advnets
	r.LocalAS = localas
	r.Networks = nets
	return r
}

// ToSpeakerCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToSpeakerCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, jroot)
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

// UpdateOptsBuilder allow the extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToSpeakerUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts, once marshalled, will be sent to API with the PUT request
type UpdateOpts struct {
	Name                          string `json:"name,omitempty"`
	AdvertiseFloatingIPHostRoutes bool   `json:"advertise_floating_ip_host_routes"`
	AdvertiseTenantNetworks       bool   `json:"advertise_tenant_networks"`
}

func (opts UpdateOpts) ToSpeakerUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, jroot)
}

func BuildUpdateOpts(name string, advrts bool, advnets bool) (r UpdateOpts) {
	r.AdvertiseFloatingIPHostRoutes = advrts
	r.AdvertiseTenantNetworks = advnets
	return r
}

func Update(c *gophercloud.ServiceClient, speakerID string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToSpeakerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(updateURL(c, speakerID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts a unique ID and deletes the bgp speaker associated with it.
func Delete(c *gophercloud.ServiceClient, speakerID string) (r DeleteResult) {
	resp, err := c.Delete(deleteURL(c, speakerID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
