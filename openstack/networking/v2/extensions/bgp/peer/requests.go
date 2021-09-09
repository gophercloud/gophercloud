package peer

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// For the sake of consistency, this function will return a pagination.Pager
func List(c *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return BGPPeerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieve the specific bgp peer by its uuid
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToPeerCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a network.
type CreateOpts struct {
	AuthType string `json:"auth_type"`
	RemoteAS int    `json:"remote_as"`
	Name     string `json:"name"`
	PeerIP   string `json:"peer_ip"`
}

// ToPeerCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToPeerCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "bgp_peer")
}

func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPeerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(createURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts a unique ID and deletes the bgp Peer associated with it.
func Delete(c *gophercloud.ServiceClient, PeerID string) (r DeleteResult) {
	resp, err := c.Delete(deleteURL(c, PeerID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
