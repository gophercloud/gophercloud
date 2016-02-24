package diskconfig

import "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

type ServerWithDiskConfig struct {
	servers.Server
	DiskConfig DiskConfig `json:"OS-DCF:diskConfig"`
}
