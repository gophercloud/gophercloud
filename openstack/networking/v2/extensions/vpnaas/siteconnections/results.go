package siteconnections

import (
	"github.com/gophercloud/gophercloud"
)

type DPD struct {
	Action   string `q:"action"`
	Timeout  int    `q:"timeout"`
	Interval int    `q:"interval"`
}

// Connection is an IPSec site connection
type Connection struct {
	IKEPolicyID        string `q:"ikepolicy_id"`
	VPNServiceID       string `q:"vpnservice_id"`
	LocalEPGroupID     string `q:"local_ep_group_id"`
	IPSecPolicyID      string `q:"ipsecpolicy_id"`
	PeerID             string `q:"peer_id"`
	TenantID           string `q:"tenant_id"`
	PeerEPGroupID      string `q:"peer_ep_group_id"`
	LocalID            string `q:"local_id"`
	Name               string `q:"name"`
	Description        string `q:"description"`
	PeerAddress        string `q:"peer_address"`
	RouteMode          string `q:"route_mode"`
	PSK                string `q:"psk"`
	Initiator          string `q:"initiator"`
	PeerCIDRs          string `q:"peer_cidrs"`
	AdminStateUp       bool   `q:"admin_state_up"`
	DPD                DPD    `q:"dpd"`
	AuthenticationMode string `q:"auth_mode"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts an IPSec site connection.
func (r commonResult) Extract() (*Connection, error) {
	var s struct {
		Connection *Connection `json:"ipsec_site_connection"`
	}
	err := r.ExtractInto(&s)
	return s.Connection, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Connection.
type CreateResult struct {
	commonResult
}
