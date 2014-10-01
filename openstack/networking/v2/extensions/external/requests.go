package external

import "github.com/rackspace/gophercloud/openstack/networking/v2/networks"

type AdminState *bool

// Convenience vars for AdminStateUp values.
var (
	iTrue  = true
	iFalse = false

	Nothing AdminState = nil
	Up      AdminState = &iTrue
	Down    AdminState = &iFalse
)

type CreateOpts struct {
	Parent   networks.CreateOpts
	External bool
}

func (o CreateOpts) ToMap() map[string]map[string]interface{} {
	outer := o.Parent.ToMap()

	outer["network"]["router:external"] = o.External

	return outer
}

func (o CreateOpts) IsCreateOpts() bool { return true }

type UpdateOpts struct {
	Parent   networks.UpdateOpts
	External bool
}

func (o UpdateOpts) ToMap() map[string]map[string]interface{} {
	outer := o.Parent.ToMap()

	outer["network"]["router:external"] = o.External

	return outer
}

func (o UpdateOpts) IsUpdateOpts() bool { return true }
