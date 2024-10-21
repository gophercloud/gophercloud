//go:build acceptance || networking || layer3 || router

package layer3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	v2 "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/agents"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestLayer3RouterScheduling(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	v2.RequireNeutronExtension(t, client, "l3_agent_scheduler")

	network, err := v2.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer v2.DeleteNetwork(t, client, network.ID)

	subnet, err := v2.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer v2.DeleteSubnet(t, client, subnet.ID)

	router, err := CreateRouter(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteRouter(t, client, router.ID)
	tools.PrintResource(t, router)

	routerInterface, err := CreateRouterInterfaceOnSubnet(t, client, subnet.ID, router.ID)
	tools.PrintResource(t, routerInterface)
	th.AssertNoErr(t, err)
	defer DeleteRouterInterface(t, client, routerInterface.PortID, router.ID)

	// List hosting agent
	allPages, err := routers.ListL3Agents(client, router.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	hostingAgents, err := routers.ExtractL3Agents(allPages)
	th.AssertNoErr(t, err)
	th.AssertIntGreaterOrEqual(t, len(hostingAgents), 1)
	hostingAgent := hostingAgents[0]
	t.Logf("Router %s is scheduled on %s", router.ID, hostingAgent.ID)

	// remove from hosting agent
	err = agents.RemoveL3Router(context.TODO(), client, hostingAgent.ID, router.ID).ExtractErr()
	th.AssertNoErr(t, err)

	containsRouterFunc := func(rs []routers.Router, routerID string) bool {
		for _, r := range rs {
			if r.ID == routerID {
				return true
			}
		}
		return false
	}

	// List routers on hosting agent
	routersOnHostingAgent, err := agents.ListL3Routers(context.TODO(), client, hostingAgent.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, containsRouterFunc(routersOnHostingAgent, router.ID), false)
	t.Logf("Router %s is not scheduled on %s", router.ID, hostingAgent.ID)

	// schedule back
	err = agents.ScheduleL3Router(context.TODO(), client, hostingAgents[0].ID, agents.ScheduleL3RouterOpts{RouterID: router.ID}).ExtractErr()
	th.AssertNoErr(t, err)

	// List hosting agent after readding
	routersOnHostingAgent, err = agents.ListL3Routers(context.TODO(), client, hostingAgent.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, containsRouterFunc(routersOnHostingAgent, router.ID), true)
	t.Logf("Router %s is scheduled on %s", router.ID, hostingAgent.ID)
}
