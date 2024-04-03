//go:build acceptance || sharedfilesystems || services

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/services"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestServicesList(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)

	client.Microversion = "2.7"
	allPages, err := services.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServices, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	th.AssertIntGreaterOrEqual(t, len(allServices), 1)

	for _, s := range allServices {
		tools.PrintResource(t, &s)
		th.AssertEquals(t, s.Status, "enabled")
	}
}
