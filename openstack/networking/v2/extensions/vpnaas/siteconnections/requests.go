package siteconnections

import "github.com/gophercloud/gophercloud"

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToConnectionCreateMap() (map[string]interface{}, error)
}
type Action string
type Initiator string

const (
	ActionHold             Action    = "hold"
	ActionClear            Action    = "clear"
	ActionRestart          Action    = "restart"
	ActionDisabled         Action    = "disabled"
	ActionRestartByPeer    Action    = "restart-by-peer"
	InitiatorBiDirectional Initiator = "bi-directional"
	InitiatorResponseOnly  Initiator = "response-only"
)

type DPDCreateOpts struct {
	Action   Action `json:"action,omitempty"`
	Timeout  int    `json:"timeout,omitempty"`
	Interval int    `json:"interval,omitempty"`
}

// CreateOpts contains all the values needed to create a new IPSec site connection
type CreateOpts struct {
	IKEPolicyID        string         `json:"ikepolicy_id,omitempty"`
	VPNServiceID       string         `json:"vpnservice_id,omitempty"`
	LocalEPGroupID     string         `json:"local_ep_group_id,omitempty"`
	IPSecPolicyID      string         `json:"ipsecpolicy_id,omitempty"`
	PeerID             string         `json:"peer_id"`
	TenantID           string         `json:"tenant_id,omitempty"`
	PeerEPGroupID      string         `json:"peer_ep_group_id,omitempty"`
	LocalID            string         `json:"local_id,omitempty"`
	Name               string         `json:"name,omitempty"`
	Description        string         `json:"description,omitempty"`
	PeerAddress        string         `json:"peer_address"`
	RouteMode          string         `json:"route_mode,omitempty"`
	PSK                string         `json:"psk"`
	Initiator          Initiator      `json:"initiator,omitempty"`
	PeerCIDRs          string         `json:"peer_cidrs,omitempty"`
	AdminStateUp       *bool          `json:"admin_state_up,omitempty"`
	DPD                *DPDCreateOpts `json:"dpd,omitempty"`
	AuthenticationMode string         `json:"auth_mode,omitempty"`
}

// ToServiceCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToConnectionCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "ipsec_site_connection")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// IPSec site connection.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToConnectionCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}
