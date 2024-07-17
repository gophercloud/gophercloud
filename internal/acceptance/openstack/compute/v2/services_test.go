//go:build acceptance || compute || services

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/services"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestServicesList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	allPages, err := services.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServices, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, service := range allServices {
		tools.PrintResource(t, service)

		if service.Binary == "nova-scheduler" {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestServicesListWithOpts(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	opts := services.ListOpts{
		Binary: "nova-scheduler",
	}

	allPages, err := services.List(client, opts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServices, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, service := range allServices {
		tools.PrintResource(t, service)
		th.AssertEquals(t, service.Binary, "nova-scheduler")

		if service.Binary == "nova-scheduler" {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}
