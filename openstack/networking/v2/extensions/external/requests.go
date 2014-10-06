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

func (o CreateOpts) ToNetworkCreateMap() map[string]map[string]interface{} {
	outer := o.Parent.ToNetworkCreateMap()

	outer["network"]["router:external"] = o.External

	return outer
}

type UpdateOpts struct {
	Parent   networks.UpdateOpts
	External bool
}

func (o UpdateOpts) ToNetworkUpdateMap() map[string]map[string]interface{} {
	outer := o.Parent.ToNetworkUpdateMap()

	outer["network"]["router:external"] = o.External

	return outer
}
