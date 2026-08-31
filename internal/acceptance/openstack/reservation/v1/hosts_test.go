//go:build acceptance || reservation || hosts

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/hosts"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestHostsList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewReservationV1Client()
	th.AssertNoErr(t, err)

	allPages, err := hosts.List(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allHosts, err := hosts.ExtractHosts(allPages)
	th.AssertNoErr(t, err)

	for _, host := range allHosts {
		tools.PrintResource(t, host)
		tools.PrintResource(t, host.CreatedAt)
		tools.PrintResource(t, host.UpdatedAt)
		tools.PrintResource(t, host.ExtraCapabilities)
	}
}
