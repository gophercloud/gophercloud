package v2

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vlantransparent"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

// VLANTransparentNetwork represents OpenStack V2 Networking Network with the
// "vlan-transparent" extension enabled.
type VLANTransparentNetwork struct {
	networks.Network
	vlantransparent.TransparentExt
}
