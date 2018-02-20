package services

import "github.com/gophercloud/gophercloud"

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToServiceCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new VPN service
type CreateOpts struct {
	// TenantID specifies a tenant to own the VPN service. The caller must have
	// an admin role in order to set this. Otherwise, this field is left unset
	// and the caller will be the owner.
	TenantID string `json:"tenant_id,omitempty"`

	// SubnetID is the ID of the subnet.
	SubnetID string `json:"subnet_id,omitempty"`

	// RouterID is the ID of the router.
	RouterID string `json:"router_id" required:"true"`

	// Description is the human readable description of the service.
	Description string `json:"description,omitempty"`

	// AdminStateUp is the administrative state of the resource, which is up (true) or down (false).
	AdminStateUp *bool `json:"admin_state_up"`

	// Name is the human readable name of the service.
	Name string `json:"name,omitempty"`

	// The ID of the flavor.
	FlavorID string `json:"flavor_id,omitempty"`
}

// ToServiceCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToServiceCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "vpnservice")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// VPN service.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToServiceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// Delete will permanently delete a particular VPN service based on its
// unique ID.
func Delete(c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
