package rules

import (
	"github.com/gophercloud/gophercloud"
)

type (
	// Protocol represents a valid rule protocol
	Protocol string
)

const (
	// ProtocolAny is to allow any protocol
	ProtocolAny Protocol = "any"

	// ProtocolICMP is to allow the ICMP protocol
	ProtocolICMP Protocol = "icmp"

	// ProtocolTCP is to allow the TCP protocol
	ProtocolTCP Protocol = "tcp"

	// ProtocolUDP is to allow the UDP protocol
	ProtocolUDP Protocol = "udp"
)

type (
	// Action represents a valid rule protocol
	Action string
)

const (
	// ActionAllow is to allow traffic
	ActionAllow Action = "allow"

	// ActionDeny is to deny traffic
	ActionDeny Action = "deny"

	// ActionTCP is to reject traffic
	ActionReject Action = "reject"
)

type CreateOptsBuilder interface {
	ToRuleCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new firewall rule.
type CreateOpts struct {
	Protocol             Protocol              `json:"protocol" required:"true"`
	Action               Action                `json:"action" required:"true"`
	TenantID             string                `json:"tenant_id,omitempty"`
	Name                 string                `json:"name,omitempty"`
	Description          string                `json:"description,omitempty"`
	IPVersion            gophercloud.IPVersion `json:"ip_version,omitempty"`
	SourceIPAddress      string                `json:"source_ip_address,omitempty"`
	DestinationIPAddress string                `json:"destination_ip_address,omitempty"`
	SourcePort           string                `json:"source_port,omitempty"`
	DestinationPort      string                `json:"destination_port,omitempty"`
	Shared               *bool                 `json:"shared,omitempty"`
	Enabled              *bool                 `json:"enabled,omitempty"`
}

// ToRuleCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToRuleCreateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "firewall_rule")
	if err != nil {
		return nil, err
	}

	if m := b["firewall_rule"].(map[string]interface{}); m["protocol"] == "any" {
		m["protocol"] = nil
	}

	return b, nil
}

// Create accepts a CreateOpts struct and uses the values to create a new firewall rule
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// Delete will permanently delete a particular firewall rule based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
