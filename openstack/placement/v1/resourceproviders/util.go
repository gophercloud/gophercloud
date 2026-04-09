package resourceproviders

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/utils"
)

func microversionAtLeast(client *gophercloud.ServiceClient, major, minor int) bool {
	M, m, err := utils.ParseMicroversion(client.Microversion)
	if err != nil {
		return false
	}
	return M > major || (M == major && m >= minor)
}
