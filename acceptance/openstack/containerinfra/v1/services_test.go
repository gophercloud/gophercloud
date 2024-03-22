//go:build acceptance || containerinfra
// +build acceptance containerinfra

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/services"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestServicesList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewContainerInfraV1Client()
	th.AssertNoErr(t, err)

	allPages, err := services.List(client).AllPages()
	th.AssertNoErr(t, err)

	allServices, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	for _, service := range allServices {
		tools.PrintResource(t, service)
	}
}
