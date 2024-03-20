//go:build acceptance || networking || vpnaas

package vpnaas

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	layer3 "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2/extensions/layer3"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/services"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestServiceList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	allPages, err := services.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServices, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	for _, service := range allServices {
		tools.PrintResource(t, service)
	}
}

func TestServiceCRUD(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/wallaby")
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	router, err := layer3.CreateExternalRouter(t, client)
	th.AssertNoErr(t, err)
	defer layer3.DeleteRouter(t, client, router.ID)

	service, err := CreateService(t, client, router.ID)
	th.AssertNoErr(t, err)
	defer DeleteService(t, client, service.ID)

	newService, err := services.Get(context.TODO(), client, service.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, service)
	tools.PrintResource(t, newService)
}
