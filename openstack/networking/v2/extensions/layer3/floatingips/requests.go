package floatingips

import (
	"fmt"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

type CreateOpts struct {
	FloatingNetworkID string
	PortID            string
	FixedIP           string
	TenantID          string
}

var (
	errFloatingNetworkIDRequired = fmt.Errorf("A NetworkID is required")
	errPortIDRequired            = fmt.Errorf("A PortID is required")
)

func Create(c *gophercloud.ServiceClient, opts CreateOpts) CreateResult {
	var res CreateResult

	// Validate
	if opts.FloatingNetworkID == "" {
		res.Err = errFloatingNetworkIDRequired
		return res
	}
	if opts.PortID == "" {
		res.Err = errPortIDRequired
		return res
	}

	// Define structures
	type floatingIP struct {
		FloatingNetworkID string `json:"floating_network_id"`
		PortID            string `json:"port_id"`
		FixedIP           string `json:"fixed_ip_address,omitempty"`
		TenantID          string `json:"tenant_id,omitempty"`
	}
	type request struct {
		FloatingIP floatingIP `json:"floatingip"`
	}

	// Populate request body
	reqBody := request{FloatingIP: floatingIP{
		FloatingNetworkID: opts.FloatingNetworkID,
		PortID:            opts.PortID,
		FixedIP:           opts.FixedIP,
		TenantID:          opts.TenantID,
	}}

	// Send request to API
	_, err := perigee.Request("POST", rootURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{201},
	})

	res.Err = err

	return res
}
