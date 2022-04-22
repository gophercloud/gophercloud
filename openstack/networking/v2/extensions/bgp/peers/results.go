package peers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

const jroot = "bgp_peer"

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a bgp peer resource.
func (r commonResult) Extract() (*BGPPeer, error) {
	var s BGPPeer
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, jroot)
}

// BGP peer
type BGPPeer struct {
	// AuthType of the BGP Speaker
	AuthType string `json:"auth_type"`

	// UUID for the bgp peer
	ID string `json:"id"`

	// Human-readable name for the bgp peer. Might not be unique.
	Name string `json:"name"`

	// TenantID is the project owner of the bgp peer.
	TenantID string `json:"tenant_id"`

	// The IP addr of the BGP Peer
	PeerIP string `json:"peer_ip"`

	// ProjectID is the project owner of the bgp peer.
	ProjectID string `json:"project_id"`

	// Remote Autonomous System
	RemoteAS int `json:"remote_as"`
}

// BGPPeerPage is the page returned by a pager when traversing over a
// collection of bgp peers.
type BGPPeerPage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether a BGPPage struct is empty.
func (r BGPPeerPage) IsEmpty() (bool, error) {
	is, err := ExtractBGPPeers(r)
	return len(is) == 0, err
}

// ExtractBGPPeers accepts a Page struct, specifically a BGPPeerPage struct,
// and extracts the elements into a slice of BGPPeer structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractBGPPeers(r pagination.Page) ([]BGPPeer, error) {
	var s []BGPPeer
	err := ExtractBGPPeersInto(r, &s)
	return s, err
}

func ExtractBGPPeersInto(r pagination.Page, v interface{}) error {
	return r.(BGPPeerPage).Result.ExtractIntoSlicePtr(v, "bgp_peers")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a BGPPeer.
type GetResult struct {
	commonResult
}
