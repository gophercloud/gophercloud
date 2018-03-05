package siteconnections

import (
	"github.com/gophercloud/gophercloud"
)

type DPD struct {
	Action   string `json:"action"`
	Timeout  int    `json:"timeout"`
	Interval int    `json:"interval"`
}

// Connection is an IPSec site connection
type Connection struct {
	IKEPolicyID        string   `json:"ikepolicy_id"`
	VPNServiceID       string   `json:"vpnservice_id"`
	LocalEPGroupID     string   `json:"local_ep_group_id"`
	IPSecPolicyID      string   `json:"ipsecpolicy_id"`
	PeerID             string   `json:"peer_id"`
	TenantID           string   `json:"tenant_id"`
	PeerEPGroupID      string   `json:"peer_ep_group_id"`
	LocalID            string   `json:"local_id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	PeerAddress        string   `json:"peer_address"`
	RouteMode          string   `json:"route_mode"`
	PSK                string   `json:"psk"`
	Initiator          string   `json:"initiator"`
	PeerCIDRs          []string `json:"peer_cidrs"`
	AdminStateUp       bool     `json:"admin_state_up"`
	DPD                DPD      `json:"dpd"`
	AuthenticationMode string   `json:"auth_mode"`
	MTU                int      `json:"mtu"`
	Status             string   `json:"status"`
	ProjectID          string   `json:"project_id"`
	ID                 string   `json:"id"`
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
