package peers

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List the bgp peers
func List(c *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return BGPPeerPage{pagination.SinglePageBase(r)}
	})
}

// Get retrieve the specific bgp peer by its uuid
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToPeerCreateMap() (map[string]any, error)
}

// CreateOpts represents options used to create a BGP Peer.
type CreateOpts struct {
	AuthType string `json:"auth_type"`
	RemoteAS int    `json:"remote_as"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	PeerIP   string `json:"peer_ip"`
}

// ToPeerCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToPeerCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, jroot)
}

// Create a BGP Peer
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOpts) (r CreateResult) {
	b, err := opts.ToPeerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, createURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts a unique ID and deletes the bgp Peer associated with it.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, bgpPeerID string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, bgpPeerID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToPeerUpdateMap() (map[string]any, error)
}

// UpdateOpts represents options used to update a BGP Peer.
type UpdateOpts struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

// ToPeerUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToPeerUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, jroot)
}

// Update accept a BGP Peer ID and an UpdateOpts and update the BGP Peer
func Update(ctx context.Context, c *gophercloud.ServiceClient, bgpPeerID string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToPeerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, updateURL(c, bgpPeerID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
