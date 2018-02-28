package endpointgroups

import "github.com/gophercloud/gophercloud"

type EndpointType string

const (
	TypeSubnet  EndpointType = "subnet"
	TypeCIDR    EndpointType = "cidr"
	TypeVLAN    EndpointType = "vlan"
	TypeNetwork EndpointType = "network"
	TypeRouter  EndpointType = "router"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToEndpointGroupCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new endpoint group
type CreateOpts struct {
	// TenantID specifies a tenant to own the endpoint group. The caller must have
	// an admin role in order to set this. Otherwise, this field is left unset
	// and the caller will be the owner.
	TenantID string `json:"tenant_id,omitempty"`

	// Description is the human readable description of the endpoint group.
	Description string `json:"description,omitempty"`

	// Name is the human readable name of the endpoint group.
	Name string `json:"name,omitempty"`

	// The type of the endpoints in the group.
	// A valid value is subnet, cidr, network, router, or vlan.
	Type EndpointType `json:"type,omitempty"`

	// List of endpoints of the same type, for the endpoint group.
	// The values will depend on the type.
	Endpoints []string `json:"endpoints"`
}

// ToEndpointGroupCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToEndpointGroupCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "endpoint_group")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// endpoint group.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToEndpointGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}
