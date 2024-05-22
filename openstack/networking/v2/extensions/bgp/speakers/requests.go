package speakers

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List the bgp speakers
func List(c *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return BGPSpeakerPage{pagination.SinglePageBase(r)}
	})
}

// Get retrieve the specific bgp speaker by its uuid
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOpts represents options used to create a BGP Speaker.
type CreateOpts struct {
	Name                          string   `json:"name"`
	IPVersion                     int      `json:"ip_version"`
	AdvertiseFloatingIPHostRoutes bool     `json:"advertise_floating_ip_host_routes"`
	AdvertiseTenantNetworks       bool     `json:"advertise_tenant_networks"`
	LocalAS                       string   `json:"local_as"`
	Networks                      []string `json:"networks,omitempty"`
}

// CreateOptsBuilder declare a function that build CreateOpts into a Create request body.
type CreateOptsBuilder interface {
	ToSpeakerCreateMap() (map[string]any, error)
}

// ToSpeakerCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToSpeakerCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, jroot)
}

// Create accepts a CreateOpts and create a BGP Speaker.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOpts) (r CreateResult) {
	b, err := opts.ToSpeakerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, createURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts a unique ID and deletes the bgp speaker associated with it.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, speakerID string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, speakerID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOpts represents options used to update a BGP Speaker.
type UpdateOpts struct {
	Name                          string `json:"name,omitempty"`
	AdvertiseFloatingIPHostRoutes bool   `json:"advertise_floating_ip_host_routes"`
	AdvertiseTenantNetworks       bool   `json:"advertise_tenant_networks"`
}

// ToSpeakerUpdateMap build a request body from UpdateOpts
func (opts UpdateOpts) ToSpeakerUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, jroot)
}

// UpdateOptsBuilder allow the extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToSpeakerUpdateMap() (map[string]any, error)
}

// Update accepts a UpdateOpts and update the BGP Speaker.
func Update(ctx context.Context, c *gophercloud.ServiceClient, speakerID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSpeakerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, updateURL(c, speakerID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// AddBGPPeerOpts represents options used to add a BGP Peer to a BGP Speaker
type AddBGPPeerOpts struct {
	BGPPeerID string `json:"bgp_peer_id"`
}

// AddBGPPeerOptsBuilder declare a funtion that encode AddBGPPeerOpts into a request body
type AddBGPPeerOptsBuilder interface {
	ToBGPSpeakerAddBGPPeerMap() (map[string]any, error)
}

// ToBGPSpeakerAddBGPPeerMap build a request body from AddBGPPeerOpts
func (opts AddBGPPeerOpts) ToBGPSpeakerAddBGPPeerMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// AddBGPPeer add the BGP peer to the speaker a.k.a. PUT /v2.0/bgp-speakers/{bgp-speaker-id}/add_bgp_peer
func AddBGPPeer(ctx context.Context, c *gophercloud.ServiceClient, bgpSpeakerID string, opts AddBGPPeerOptsBuilder) (r AddBGPPeerResult) {
	b, err := opts.ToBGPSpeakerAddBGPPeerMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, addBGPPeerURL(c, bgpSpeakerID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveBGPPeerOpts represents options used to remove a BGP Peer to a BGP Speaker
type RemoveBGPPeerOpts AddBGPPeerOpts

// RemoveBGPPeerOptsBuilder declare a funtion that encode RemoveBGPPeerOpts into a request body
type RemoveBGPPeerOptsBuilder interface {
	ToBGPSpeakerRemoveBGPPeerMap() (map[string]any, error)
}

// ToBGPSpeakerRemoveBGPPeerMap build a request body from RemoveBGPPeerOpts
func (opts RemoveBGPPeerOpts) ToBGPSpeakerRemoveBGPPeerMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// RemoveBGPPeer remove the BGP peer from the speaker, a.k.a. PUT /v2.0/bgp-speakers/{bgp-speaker-id}/add_bgp_peer
func RemoveBGPPeer(ctx context.Context, c *gophercloud.ServiceClient, bgpSpeakerID string, opts RemoveBGPPeerOptsBuilder) (r RemoveBGPPeerResult) {
	b, err := opts.ToBGPSpeakerRemoveBGPPeerMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, removeBGPPeerURL(c, bgpSpeakerID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetAdvertisedRoutes a.k.a. GET /v2.0/bgp-speakers/{bgp-speaker-id}/get_advertised_routes
func GetAdvertisedRoutes(c *gophercloud.ServiceClient, bgpSpeakerID string) pagination.Pager {
	url := getAdvertisedRoutesURL(c, bgpSpeakerID)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return AdvertisedRoutePage{pagination.SinglePageBase(r)}
	})
}

// AddGatewayNetworkOptsBuilder declare a function that build AddGatewayNetworkOpts into a request body.
type AddGatewayNetworkOptsBuilder interface {
	ToBGPSpeakerAddGatewayNetworkMap() (map[string]any, error)
}

// AddGatewayNetworkOpts represents the data that would be PUT to the endpoint
type AddGatewayNetworkOpts struct {
	// The uuid of the network
	NetworkID string `json:"network_id"`
}

// ToBGPSpeakerAddGatewayNetworkMap implements the function
func (opts AddGatewayNetworkOpts) ToBGPSpeakerAddGatewayNetworkMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// AddGatewayNetwork a.k.a. PUT /v2.0/bgp-speakers/{bgp-speaker-id}/add_gateway_network
func AddGatewayNetwork(ctx context.Context, c *gophercloud.ServiceClient, bgpSpeakerID string, opts AddGatewayNetworkOptsBuilder) (r AddGatewayNetworkResult) {
	b, err := opts.ToBGPSpeakerAddGatewayNetworkMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, addGatewayNetworkURL(c, bgpSpeakerID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveGatewayNetworkOptsBuilder declare a function that build RemoveGatewayNetworkOpts into a request body.
type RemoveGatewayNetworkOptsBuilder interface {
	ToBGPSpeakerRemoveGatewayNetworkMap() (map[string]any, error)
}

// RemoveGatewayNetworkOpts represent the data that would be PUT to the endpoint
type RemoveGatewayNetworkOpts AddGatewayNetworkOpts

// ToBGPSpeakerRemoveGatewayNetworkMap implement the function
func (opts RemoveGatewayNetworkOpts) ToBGPSpeakerRemoveGatewayNetworkMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// RemoveGatewayNetwork a.k.a. PUT /v2.0/bgp-speakers/{bgp-speaker-id}/remove_gateway_network
func RemoveGatewayNetwork(ctx context.Context, c *gophercloud.ServiceClient, bgpSpeakerID string, opts RemoveGatewayNetworkOptsBuilder) (r RemoveGatewayNetworkResult) {
	b, err := opts.ToBGPSpeakerRemoveGatewayNetworkMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, removeGatewayNetworkURL(c, bgpSpeakerID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
