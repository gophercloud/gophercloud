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
	TenantID    	string  `json:"tenant_id,omitempty"`
	SubnetID    	string 	`json:"subnet_id,omitempty"`
	RouterID    	string 	`json:"router_id", required:"true"`
	Description 	string  `json:"description,omitempty"`
	AdminStateUp 	*bool 	`json:"admin_state_up"`
	ProjectID   	string  `json:"project_id"`
	Name 			string  `json:"name,omitempty"`
	FlavorID        string  `json:"flavor_id"`
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

