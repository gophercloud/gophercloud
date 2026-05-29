//go:build acceptance || networking || layer3 || router

package layer3

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/extraroutes"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestLayer3ExtraRoutesAddRemove(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)
	tmp := net.ParseIP(subnet.GatewayIP).To4()
	if tmp == nil {
		th.AssertNoErr(t, fmt.Errorf("invalid subnet gateway IP: %s", subnet.GatewayIP))
	}
	tmp[3] = 251
	gateway := tmp.String()

	router, err := CreateRouter(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteRouter(t, client, router.ID)

	tools.PrintResource(t, router)

	aiOpts := routers.AddInterfaceOpts{
		SubnetID: subnet.ID,
	}
	iface, err := routers.AddInterface(context.TODO(), client, router.ID, aiOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, iface)

	// 2. delete router interface
	defer func() {
		riOpts := routers.RemoveInterfaceOpts{
			SubnetID: subnet.ID,
		}
		_, err = routers.RemoveInterface(context.TODO(), client, router.ID, riOpts).Extract()
		th.AssertNoErr(t, err)
	}()

	// 1. delete routes first
	defer func() {
		routes := []routers.Route{}
		opts := routers.UpdateOpts{
			Routes: &routes,
		}
		_, err = routers.Update(context.TODO(), client, router.ID, opts).Extract()
		th.AssertNoErr(t, err)
	}()

	routes := []routers.Route{
		{
			DestinationCIDR: "192.168.11.0/30",
			NextHop:         gateway,
		},
		{
			DestinationCIDR: "192.168.12.0/30",
			NextHop:         gateway,
		},
	}
	updateOpts := routers.UpdateOpts{
		Routes: &routes,
	}
	_, err = routers.Update(context.TODO(), client, router.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	newRoutes := []routers.Route{
		{
			DestinationCIDR: "192.168.13.0/30",
			NextHop:         gateway,
		},
		{
			DestinationCIDR: "192.168.14.0/30",
			NextHop:         gateway,
		},
	}
	opts := extraroutes.Opts{
		Routes: &newRoutes,
	}
	// add new routes
	rt, err := extraroutes.Add(context.TODO(), client, router.ID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, append(routes, newRoutes...), rt.Routes)

	// remove new routes
	rt, err = extraroutes.Remove(context.TODO(), client, router.ID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, routes, rt.Routes)
}
