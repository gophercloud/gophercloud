//go:build acceptance || networking || networkipavailabilities

package networkipavailabilities

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/networkipavailabilities"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestNetworkIPAvailabilityList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	allPages, err := networkipavailabilities.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allAvailabilities, err := networkipavailabilities.ExtractNetworkIPAvailabilities(allPages)
	th.AssertNoErr(t, err)

	for _, availability := range allAvailabilities {
		for _, subnet := range availability.SubnetIPAvailabilities {
			tools.PrintResource(t, subnet)
			tools.PrintResource(t, subnet.TotalIPs)
			tools.PrintResource(t, subnet.UsedIPs)
		}
	}
}
