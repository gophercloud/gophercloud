//go:build acceptance || networking || layer3 || router
// +build acceptance networking layer3 router

package layer3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/agents"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestLayer3RouterScheduling(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	router, err := CreateRouter(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteRouter(t, client, router.ID)
	tools.PrintResource(t, router)

	routerInterface, err := CreateRouterInterfaceOnSubnet(t, client, subnet.ID, router.ID)
	tools.PrintResource(t, routerInterface)
	th.AssertNoErr(t, err)
	defer DeleteRouterInterface(t, client, routerInterface.PortID, router.ID)

	// List hosting agent
	allPages, err := routers.ListL3Agents(client, router.ID).AllPages()
	th.AssertNoErr(t, err)
	hostingAgents, err := routers.ExtractL3Agents(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(hostingAgents) > 0, true)
	hostingAgent := hostingAgents[0]
	t.Logf("Router %s is scheduled on %s", router.ID, hostingAgent.ID)

	// remove from hosting agent
	err = agents.RemoveL3Router(client, hostingAgent.ID, router.ID).ExtractErr()
	th.AssertNoErr(t, err)

	containsRouterFunc := func(rs []routers.Router, routerID string) bool {
		for _, r := range rs {
			if r.ID == router.ID {
				return true
			}
		}
		return false
	}

	// List routers on hosting agent
	routersOnHostingAgent, err := agents.ListL3Routers(client, hostingAgent.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, containsRouterFunc(routersOnHostingAgent, router.ID), false)
	t.Logf("Router %s is not scheduled on %s", router.ID, hostingAgent.ID)

	// schedule back
	err = agents.ScheduleL3Router(client, hostingAgents[0].ID, agents.ScheduleL3RouterOpts{RouterID: router.ID}).ExtractErr()
	th.AssertNoErr(t, err)

	// List hosting agent after readding
	routersOnHostingAgent, err = agents.ListL3Routers(client, hostingAgent.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, containsRouterFunc(routersOnHostingAgent, router.ID), true)
	t.Logf("Router %s is scheduled on %s", router.ID, hostingAgent.ID)
}
