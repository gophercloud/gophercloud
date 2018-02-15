package services

import "github.com/gophercloud/gophercloud"

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToVPNServiceCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new VPN service
type CreateOpts struct {
	// TenantID specifies a tenant to own the VPN service. The caller must have
	// an admin role in order to set this. Otherwise, this field is left unset
	// and the caller will be the owner.
	TenantID     string `json:"tenant_id,omitempty"`
	//If you specify only a subnet UUID, OpenStack Networking allocates an available IP from that subnet to the port
	//If you specify both a subnet UUID and an IP address, OpenStack Networking tries to allocate the address to the port
	SubnetID     string `json:"subnet_id,omitempty"`
	//The ID of the router
	RouterID     string `json:"router_id" required:"true"`
	//Human readable description of the service
	Description  string `json:"description,omitempty"`
	//The administrative state of the resource, which is up (true) or down (false).
	AdminStateUp *bool  `json:"admin_state_up"`
	//The ID of the project
	ProjectID    string `json:"project_id"`
	//The human readable name of the service
	Name         string `json:"name,omitempty"`
	//The ID of the flavor
	FlavorID     string `json:"flavor_id,omitempty"`
}

// ToVPNServiceCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToVPNServiceCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "vpnservice")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// VPN service
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToVPNServiceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}
