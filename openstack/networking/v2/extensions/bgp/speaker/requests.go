package speaker

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"log"
	"strconv"
	"strings"
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

func BuildCreateOpts(name string, localas string, nets []string, opts map[string]string) (r CreateOpts) {
	r.Name = name
	r.LocalAS = localas
	r.Networks = nets

	ipver, ok := opts["IPVersion"]
	if !ok {
		r.IPVersion = 4
	} else if ipver == "4" || ipver == "6" {
		r.IPVersion, _ = strconv.Atoi(ipver)
	} else {
		log.Panic("IPVersion should be either 4 or 6")
	}

	advrts, ok := opts["AdvertiseFloatingIPHostRoutes"]
	if !ok {
		r.AdvertiseFloatingIPHostRoutes = true
	} else if strings.ToLower(advrts) == "true" {
		r.AdvertiseFloatingIPHostRoutes = true
	} else if strings.ToLower(advrts) == "false" {
		r.AdvertiseFloatingIPHostRoutes = false
	} else {
		log.Panic("AdvertiseFloatingIPHostRoutes should be either true or false")
	}

	advnets, ok := opts["AdvertiseTenantNetworks"]
	if !ok {
		r.AdvertiseTenantNetworks = true
	} else if strings.ToLower(advnets) == "true" {
		r.AdvertiseTenantNetworks = true
	} else if strings.ToLower(advnets) == "false" {
		r.AdvertiseTenantNetworks = false
	} else {
		log.Panic("AdvertiseTenantNetworks sho uld be either true or false")
	}
	return
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

func BuildUpdateOpts(c *gophercloud.ServiceClient, speakerID string, m map[string]string) (r UpdateOpts) {
	name, ok := m["Name"]
	if ok {
		r.Name = name
	}

	speaker := func() *BGPSpeaker {
		s, err := Get(c, speakerID).Extract()
		if err != nil {
			log.Panic(err)
		}
		return s
	}()

	advrts, ok := m["AdvertiseFloatingIPHostRoutes"]
	if ok {
		if strings.ToLower(advrts) == "true" {
			r.AdvertiseFloatingIPHostRoutes = true
		} else if strings.ToLower(advrts) == "false" {
			r.AdvertiseFloatingIPHostRoutes = false
		} else {
			log.Panic("AdvertiseFloatingIPHostRoutes should be either true or false")
		}
	} else {
		r.AdvertiseFloatingIPHostRoutes = (*speaker).AdvertiseFloatingIPHostRoutes
	}

	advnets, ok := m["AdvertiseTenantNetworks"]
	if ok {
		if strings.ToLower(advnets) == "true" {
			r.AdvertiseTenantNetworks = true
		} else if strings.ToLower(advnets) == "false" {
			r.AdvertiseTenantNetworks = false
		} else {
			log.Panic("AdvertiseTenantNetworks should be either true or false")
		}
	} else {
		r.AdvertiseTenantNetworks = (*speaker).AdvertiseTenantNetworks
	}
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

type AddBGPPeerOptBuilder interface {
	ToBGPSpeakerAddBGPPeerMap() (map[string]interface{}, error)
}

type AddBGPPeerOpts struct {
	BGPPeerID string `json:"bgp_peer_id"`
}

func (opts AddBGPPeerOpts) ToBGPSpeakerAddBGPPeerMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func AddBGPPeer(c *gophercloud.ServiceClient, bgpSpeakerID string, opts AddBGPPeerOpts) (r AddBGPPeerResult) {
	b, err := opts.ToBGPSpeakerAddBGPPeerMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(addBGPPeerURL(c, bgpSpeakerID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type RemoveBGPPeerOptBuilder interface {
	ToBGPSpeakerRemoveBGPPeerMap() (map[string]interface{}, error)
}

type RemoveBGPPeerOpts struct {
	BGPPeerID string `json:"bgp_peer_id"`
}

func (opts RemoveBGPPeerOpts) ToBGPSpeakerRemoveBGPPeerMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func RemoveBGPPeer(c *gophercloud.ServiceClient, bgpSpeakerID string, opts RemoveBGPPeerOpts) (r RemoveBGPPeerResult) {
	b, err := opts.ToBGPSpeakerRemoveBGPPeerMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(removeBGPPeerURL(c, bgpSpeakerID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
