// +build acceptance compute services

package v2

import (
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/services"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestServicesList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	allPages, err := services.List(client, nil).AllPages()
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

	allPages, err := services.List(client, opts).AllPages()
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

func TestServicesUpdate(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	client.Microversion = "2.53"
	allPages, err := services.List(client).AllPages()
	th.AssertNoErr(t, err)

	allServices, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	// disable all services
	for _, service := range allServices {
		// only update nova-compute service
		if !strings.EqualFold(service.Binary, "nova-compute") {
			continue
		}

		opts := services.UpdateOpts{
			Status: services.ServiceDisabled,
		}
		updated, err := services.Update(client, service.ID, opts).Extract()
		th.AssertNoErr(t, err)

		th.AssertEquals(t, updated.ID, service.ID)
	}

	// verify all services are disabled
	allPages, err = services.List(client).AllPages()
	th.AssertNoErr(t, err)

	allServices, err = services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	for _, service := range allServices {
		// only check nova-compute service
		if !strings.EqualFold(service.Binary, "nova-compute") {
			continue
		}

		th.AssertEquals(t, service.Status, "disabled")
	}

	// reenable all services
	allPages, err = services.List(client).AllPages()
	th.AssertNoErr(t, err)

	allServices, err = services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	for _, service := range allServices {
		// only update nova-compute service
		if !strings.EqualFold(service.Binary, "nova-compute") {
			continue
		}

		opts := services.UpdateOpts{
			Status: services.ServiceEnabled,
		}
		updated, err := services.Update(client, service.ID, opts).Extract()
		th.AssertNoErr(t, err)

		th.AssertEquals(t, updated.ID, service.ID)
	}
}
