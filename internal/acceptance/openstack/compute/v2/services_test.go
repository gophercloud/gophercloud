//go:build acceptance || compute || services

package v2

import (
	"context"
	"testing"
	"time"

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

	th.AssertEquals(t, true, found)
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
		th.AssertEquals(t, "nova-scheduler", service.Binary)

		if service.Binary == "nova-scheduler" {
			found = true
		}
	}

	th.AssertEquals(t, true, found)
}

func TestServicesUpdate(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	listOpts := services.ListOpts{
		Binary: "nova-compute",
	}

	client.Microversion = "2.53"
	allPages, err := services.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServices, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	// disable all services
	for _, service := range allServices {
		opts := services.UpdateOpts{
			Status: services.ServiceDisabled,
		}
		updated, err := services.Update(context.TODO(), client, service.ID, opts).Extract()
		th.AssertNoErr(t, err)

		th.AssertEquals(t, updated.ID, service.ID)
	}

	// verify all services are disabled
	allPages, err = services.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServices, err = services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	for _, service := range allServices {
		th.AssertEquals(t, "disabled", service.Status)
	}

	// reenable all services
	allPages, err = services.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServices, err = services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	for _, service := range allServices {
		opts := services.UpdateOpts{
			Status: services.ServiceEnabled,
		}
		updated, err := services.Update(context.TODO(), client, service.ID, opts).Extract()
		th.AssertNoErr(t, err)

		th.AssertEquals(t, updated.ID, service.ID)
	}

	// verify all services are enabled
	allPages, err = services.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServices, err = services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	for _, service := range allServices {
		th.AssertEquals(t, "enabled", service.Status)
	}

	// Just checking if compute eventually becomes available
	time.Sleep(5 * time.Second)
}
